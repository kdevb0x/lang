package ir

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/driusan/lang/parser/ast"
)

type EnumMap map[string]int

// Compile takes an AST and writes the assembly that it compiles to to
// w.
func Generate(node ast.Node, typeInfo ast.TypeInformation, callables ast.Callables, enums EnumMap) (Func, EnumMap, error) {
	context := &variableLayout{
		make(map[ast.VarWithType]Register),
		0,
		0,
		typeInfo,
		nil,
		nil,
		enums,
		callables,
		0,
		0,
	}
	switch n := node.(type) {
	case ast.ProcDecl:
		context.funcargs = n.GetArgs()
		nargs := 0
		for i, arg := range n.Args {
			nargs++
			// Slices get passed as {n, *void}, so claim an extra argument in the
			// IR, that way code generation will make sure other variables on the
			// stack start at the right place.
			if _, ok := arg.Typ.(ast.SliceType); ok {
				nargs++
			}
			context.FuncParamRegister(arg, i)
		}
		for _, rv := range n.Return {
			context.rettypes = append(context.rettypes, context.GetTypeInfo(rv.Type()))
		}
		body, err := compileBlock(n.Body, context)
		if err != nil {
			return Func{}, nil, err
		}
		return Func{Name: n.Name, Body: body, NumArgs: uint(nargs), NumLocals: context.numLocals, LargestFuncCall: context.maxFuncCall}, enums, nil
	case ast.FuncDecl:
		nargs := 0
		for i, arg := range n.Args {
			nargs++
			if _, ok := arg.Typ.(ast.SliceType); ok {
				nargs++
			}
			context.FuncParamRegister(arg, i)
		}
		for _, rv := range n.Return {
			words := strings.Fields(string(rv.Type()))
			for _, typePiece := range words {
				context.rettypes = append(context.rettypes, context.GetTypeInfo(typePiece))
			}
		}
		body, err := compileBlock(n.Body, context)
		if err != nil {
			return Func{}, nil, err
		}
		return Func{Name: n.Name, Body: body, NumArgs: uint(nargs), NumLocals: context.numLocals, LargestFuncCall: context.maxFuncCall}, enums, nil
	case ast.SumTypeDefn:
		e := make(EnumMap)
		for i, v := range n.Options {
			e[v.Constructor] = i
		}
		return Func{}, e, nil
	case ast.TypeDefn:
		// Do nothing, the types have already been validated
		return Func{}, enums, fmt.Errorf("No IR to generate for type definitions.")
	default:
		panic(fmt.Sprintf("Unhandled Node type in compiler %v", reflect.TypeOf(n)))
	}
}

// calculate the IR to perform a function call.
func callFunc(fc ast.FuncCall, context *variableLayout, tailcall bool) ([]Opcode, error) {
	var ops []Opcode
	var argRegs []Register
	var signature ast.Callable
	if s := context.callables[fc.Name]; len(s) > 1 {
		return nil, fmt.Errorf("Multiple dispatch not yet implemented")
	} else if len(s) < 1 {
		return nil, fmt.Errorf("Can not call undefined function %v", fc.Name)
	} else {
		signature = s[0]
	}
	var funcArgs []ast.VarWithType
	if signature != nil {
		funcArgs = signature.GetArgs()
	}

	for i, arg := range fc.UserArgs {
		switch a := arg.(type) {
		case ast.EnumValue:
			argRegs = append(argRegs, getRegister(a, context))
			for _, v := range a.Parameters {
				arg, r, err := evaluateValue(v, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, arg...)
				argRegs = append(argRegs, r)
			}
		case ast.StringLiteral, ast.IntLiteral, ast.BoolLiteral:
			argRegs = append(argRegs, getRegister(a, context))
		case ast.ArrayValue:
			newops, r, err := evaluateValue(a, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, newops...)
			argRegs = append(argRegs, r)
		case ast.VarWithType:
			switch st := a.Typ.(type) {
			case ast.SliceType:
				lv := context.Get(a)
				// Slice types have 2 internal representations:
				//     struct{ n int, [n]foo}
				// where n foos directly follow the size (this form is used when they're allocated) and:
				//     struct{ n int, first *foo}
				// where a pointer to the first foo follows n (this form is used when they're passed around).
				// This should be harmonized (probably by getting rid of the first) but for now w
				// need to handle both.
				switch l := lv.(type) {
				case LocalValue:
					val1, ok := context.SafeGet(ast.VarWithType{
						Name: ast.Variable(fmt.Sprintf("%s[%d]", a.Name, 0)),
						Typ:  st.Base,
					})
					if !ok {
						val1 = context.Get(ast.VarWithType{
							Name:      ast.Variable(fmt.Sprintf("%s[%d]", a.Name, 0)),
							Typ:       st.Base,
							Reference: true,
						})
					}
					argRegs = append(argRegs, lv)
					argRegs = append(argRegs, Pointer{val1})
				case FuncArg:
					argRegs = append(argRegs, FuncArg{
						Id:   l.Id,
						Info: ast.TypeInfo{8, false},
					})
					argRegs = append(argRegs, FuncArg{
						Id:   l.Id + 1,
						Info: ast.TypeInfo{8, false},
					})
				default:
					panic(fmt.Sprintf("This should not happen %v", reflect.TypeOf(lv)))
				}
			default:
				lv := context.Get(a)
				if funcArgs != nil && funcArgs[i].Reference {
					lv = Pointer{lv}
				}
				argRegs = append(argRegs, lv)
			}
		case ast.FuncCall:
			// a function call as a parameter to a function call in
			// a return statement shouldn't be tail call optimized,
			// only the return call itself.
			fc, err := callFunc(a, context, false)
			if err != nil {
				return nil, err
			}
			ops = append(ops, fc...)

			ti := context.typeinfo[a.Returns[0].Type()]
			if ti.Size == 0 {
				ti.Size = 8
			}
			reg := context.NextLocalRegister(ast.VarWithType{"", ast.TypeLiteral(a.Returns[0].Type()), false})
			ops = append(ops,
				MOV{
					Src: FuncRetVal{0, ti},
					Dst: reg,
				},
			)
			argRegs = append(argRegs, reg)
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
			arg, r, err := evaluateValue(a, context)

			if err != nil {
				return nil, err
			}
			ops = append(ops, arg...)
			argRegs = append(argRegs, r)
		default:
			panic(fmt.Sprintf("Unhandled argument type in FuncCall %v", reflect.TypeOf(a)))
		}
	}

	// FIXME: This shouldn't need to reserve this much stack space, but
	// something gets corrupted somewhere in the standard library if we
	// don't..
	//if argSize := uint(len(argRegs)); argSize > context.maxFuncCall {
	context.maxFuncCall += uint(len(argRegs))
	//}

	// Perform the call.
	if fc.Name == "print" {
		ops = append(ops, CALL{FName: "printf", Args: argRegs, TailCall: tailcall})
	} else {
		ops = append(ops, CALL{FName: Fname(fc.Name), Args: argRegs, TailCall: tailcall})
	}
	return ops, nil
}

var loopNum uint

func getRegister(n ast.Node, context *variableLayout) Register {
	switch v := n.(type) {
	case ast.StringLiteral:
		return StringLiteral(v)
	case ast.IntLiteral:
		return IntLiteral(v)
	case ast.BoolLiteral:
		if v {
			return IntLiteral(1)
		}
		return IntLiteral(0)
	case ast.VarWithType:
		return context.Get(v)
	case ast.EnumOption:
		return IntLiteral(context.GetEnumIndex(v.Constructor))
	case ast.EnumValue:
		return IntLiteral(context.GetEnumIndex(v.Constructor.Constructor))
	default:
		panic(fmt.Sprintf("Unhandled type in getRegister: %v", reflect.TypeOf(v)))
	}
}

func compileBlock(block ast.BlockStmt, context *variableLayout) ([]Opcode, error) {
	var ops []Opcode
	for _, stmt := range block.Stmts {
		switch s := stmt.(type) {
		case ast.FuncCall:
			fc, err := callFunc(s, context, false)
			if err != nil {
				return nil, err
			}
			ops = append(ops, fc...)
		case ast.LetStmt:
			ov, oldval := context.values[s.Var]
			reg := context.NextLocalRegister(s.Var)
			if oldval {
				context.values[s.Var] = ov
			}

			switch v := s.Value.(type) {
			case ast.IntLiteral, ast.StringLiteral, ast.BoolLiteral:
				ops = append(ops, MOV{
					Src: getRegister(v, context),
					Dst: reg,
				})
			case ast.EnumValue:
				ops = append(ops, MOV{
					Src: getRegister(v, context),
					Dst: reg,
				})
				// FIXME: Need to handle parameters here.
			case ast.AdditionOperator, ast.SubtractionOperator,
				ast.DivOperator, ast.MulOperator, ast.ModOperator,
				ast.GreaterComparison, ast.GreaterOrEqualComparison,
				ast.EqualityComparison, ast.NotEqualsComparison,
				ast.LessThanComparison, ast.LessThanOrEqualComparison,
				ast.ArrayValue:
				body, r, err := evaluateValue(s.Value, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, MOV{
					Src: r,
					Dst: reg,
				})
			case ast.FuncCall:
				fc, err := callFunc(v, context, false)
				if err != nil {
					return nil, err
				}
				ops = append(ops, fc...)

				multiwordoffset := 0
				for i, v := range v.Returns {
					words := strings.Fields(string(v.Type()))
					for word := range words {
						var r Register
						ti := context.GetTypeInfo(words[word])

						if word == 0 {
							r = reg
						} else {
							r = context.NextLocalRegister(ast.VarWithType{"", ast.TypeLiteral(words[word]), false})
						}
						ops = append(ops, MOV{
							Src: FuncRetVal{uint(i + word + multiwordoffset), ti},
							Dst: r,
						})
					}
					multiwordoffset += len(words) - 1
				}
			case ast.VarWithType:
				// A let statement being assigned to a variable doesn't need any IR, it just needs to make sure that the reference points
				// to the right place.
				// The verification that nothing gets modified happens at the AST level.
				// FIXME: This should make a copy if the
				// reference to the variable.
				vr := context.Get(v)
				context.SetLocalRegister(s.Var, vr)
			case ast.ArrayLiteral:
				regs := make([]Register, len(v))
				// First generate the LocalValue registers to ensure they're consecutive if there's a variable
				// or some other expression in one of the literal pieces.
				isSlice := false
				var baseType ast.Type
				if st, ok := s.Var.Typ.(ast.SliceType); ok {
					// Move the size to the start of the slice.
					isSlice = true
					if lv, ok := reg.(LocalValue); ok {
						lv.Info = ast.TypeInfo{8, false}
						ops = append(ops, MOV{
							Src: IntLiteral(len(v)),
							Dst: lv,
						})
						context.values[s.Var] = lv
					} else {
						panic("Unexpected LocalValue")
					}
					baseType = st.Base
				}

				for i := range regs {
					entryVar := ast.VarWithType{ast.Variable(fmt.Sprintf("%s[%d]", s.Var.Name, i)), ast.TypeLiteral(v[i].Type()), false}
					if isSlice {
						entryVar.Typ = baseType
					} else if i == 0 { // && !isSlice
						// Convert the type information for the first LocalValue allocated to match
						// foo[0], not the (useless) defaults that foo generated, rather than allocating
						// a new register.
						tr, ok := reg.(LocalValue)
						if !ok {
							panic("Register for array is not a LocalValue")
						}
						tr.Info = context.GetTypeInfo(v[i].Type())
						regs[i] = tr
						context.values[entryVar] = regs[i]
						context.tempVars--
						continue
					}
					// Allocate a new LocalValue for foo[0...n]
					regs[i] = context.NextLocalRegister(entryVar)
				}

				// Then evaluate the values and put them in the appropriate registers
				for i, r := range regs {
					body, val, err := evaluateValue(v[i], context)
					if err != nil {
						return nil, err
					}
					ops = append(ops, body...)

					ops = append(ops, MOV{
						Src: val,
						Dst: r,
					})
				}
			default:
				panic(fmt.Sprintf("Unsupported let statement assignment type: %v", reflect.TypeOf(v)))
			}
			if oldval {
				context.values[s.Var] = reg
			}
		case ast.ReturnStmt:
			switch arg := s.Val.(type) {
			case ast.FuncCall:
				fc, err := callFunc(arg, context, true)
				if err != nil {
					return nil, err
				}
				ops = append(ops, fc...)
				// Calling the function already will have left
				// the value in FuncRetValRegister[0]
			case ast.EnumValue:
				// The variant of the enum goes into FR0
				ops = append(ops, MOV{
					Src: getRegister(arg, context),
					Dst: FuncRetVal{0, context.GetReturnTypeInfo(0)},
				})

				// The parameters go into FRn + i
				for i, v := range arg.Parameters {
					ti := context.GetTypeInfo(v.Type())

					ops = append(ops, MOV{
						Src: getRegister(arg.Parameters[i], context),
						Dst: FuncRetVal{1 + uint(i), ti},
					})
				}
			case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
				body, r, err := evaluateValue(arg, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, MOV{
					Src: r,
					Dst: FuncRetVal{0, context.GetReturnTypeInfo(0)},
				})
			default:
				if len(context.rettypes) != 0 {
					ops = append(ops, MOV{
						Src: getRegister(arg, context),
						Dst: FuncRetVal{0, context.GetReturnTypeInfo(0)},
					})
				}
			}
			ops = append(ops, RET{})
		case ast.MutStmt:
			reg := context.NextLocalRegister(s.Var)
			switch v := s.InitialValue.(type) {
			case ast.IntLiteral, ast.BoolLiteral, ast.StringLiteral:
				ops = append(ops, MOV{
					Src: getRegister(s.InitialValue, context),
					Dst: reg,
				})
			case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator,
				ast.ArrayValue:
				body, r, err := evaluateValue(s.InitialValue, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, MOV{
					Src: r,
					Dst: reg,
				})
			case ast.FuncCall:
				fc, err := callFunc(v, context, false)
				if err != nil {
					return nil, err
				}
				ops = append(ops, fc...)

				multiwordoffset := 0
				for i, v := range v.Returns {
					words := strings.Fields(string(v.Type()))
					for word := range words {
						var r Register
						ti := context.GetTypeInfo(words[word])

						if word == 0 {
							r = reg
						} else {
							r = context.NextLocalRegister(ast.VarWithType{"", ast.TypeLiteral(words[word]), false})
						}
						ops = append(ops, MOV{
							Src: FuncRetVal{uint(i + word + multiwordoffset), ti},
							Dst: r,
						})
					}
					multiwordoffset += len(words) - 1
				}
			case ast.VarWithType:
				val := context.Get(v)
				ops = append(ops, MOV{
					Src: val,
					Dst: reg,
				})
			case ast.ArrayLiteral:
				regs := make([]Register, len(v))
				isSlice := false
				var baseType ast.Type
				if st, ok := s.Var.Typ.(ast.SliceType); ok {
					// Move the size to the start of the slice.
					isSlice = true
					if lv, ok := reg.(LocalValue); ok {
						lv.Info = ast.TypeInfo{8, false}
						ops = append(ops, MOV{
							Src: IntLiteral(len(v)),
							Dst: lv,
						})
						context.values[s.Var] = lv
					} else {
						panic("Unexpected LocalValue")
					}
					baseType = st.Base
				}
				// First generate the LocalValue registers to ensure they're consecutive if there's a variable
				// or some other expression in one of the literal pieces.
				for i := range regs {
					entryVar := ast.VarWithType{ast.Variable(fmt.Sprintf("%s[%d]", s.Var.Name, i)), ast.TypeLiteral(v[i].Type()), false}
					if isSlice {
						entryVar.Typ = baseType
					} else if i == 0 { // && !isSlice
						// Convert the type information for the first LocalValue allocated to match
						// foo[0], not the (useless) defaults that foo generated, rather than allocating
						// a new register.
						tr, ok := reg.(LocalValue)
						if !ok {
							panic("Register for array is not a LocalValue")
						}
						tr.Info = context.GetTypeInfo(v[i].Type())
						regs[i] = tr
						context.values[entryVar] = regs[i]
						context.tempVars--
						continue
					}
					// Allocate a new LocalValue for foo[1...n]
					regs[i] = context.NextLocalRegister(entryVar)
				}

				// Then evaluate the values and put them in the appropriate registers
				for i, r := range regs {
					body, val, err := evaluateValue(v[i], context)
					if err != nil {
						return nil, err
					}
					ops = append(ops, body...)

					ops = append(ops, MOV{
						Src: val,
						Dst: r,
					})
				}
			default:
				panic(fmt.Sprintf("Unhandled type for MutStmt assignment %v", reflect.TypeOf(s.InitialValue)))
			}
		case ast.AssignmentOperator:
			dst := context.Get(s.Variable)
			switch v := s.Value.(type) {
			case ast.IntLiteral, ast.BoolLiteral, ast.StringLiteral:
				ops = append(ops, MOV{
					Src: getRegister(s.Value, context),
					Dst: dst,
				})
			case ast.AdditionOperator, ast.SubtractionOperator, ast.DivOperator, ast.MulOperator, ast.ModOperator:
				body, r, err := evaluateValue(s.Value, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, MOV{
					Src: r,
					Dst: dst,
				})
			case ast.FuncCall:
				fc, err := callFunc(v, context, false)
				if err != nil {
					return nil, err
				}
				ops = append(ops, fc...)

				multiwordoffset := 0
				for i, v := range v.Returns {
					words := strings.Fields(string(v.Type()))
					for word := range words {
						r := context.NextTempRegister()
						ti := context.GetTypeInfo(words[word])
						if ti.Size == 0 {
							ti.Size = 8
						}

						if word == 0 {
							ops = append(ops, MOV{
								Src: FuncRetVal{uint(i + word + multiwordoffset), ti},
								Dst: dst,
							})
						} else {
							ops = append(ops, MOV{
								Src: FuncRetVal{uint(i + word + multiwordoffset), ti},
								Dst: r,
							})
						}
						multiwordoffset += len(words) - 1
					}
				}
			default:
				panic(fmt.Sprintf("Statement type assignment not implemented: %v", reflect.TypeOf(s.Value)))

			}
		case ast.WhileLoop:
			lname := Label(fmt.Sprintf("loop%dend", loopNum))
			lcond := Label(fmt.Sprintf("loop%dcond", loopNum))
			loopNum++

			ops = append(ops, lcond)
			body, err := evaluateCondition(s.Condition, context, lname)
			if err != nil {
				return nil, err
			}
			ops = append(ops, body...)

			body, err = compileBlock(s.Body, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, body...)

			ops = append(ops, JMP{lcond})
			ops = append(ops, lname)
		case ast.IfStmt:
			iname := Label(fmt.Sprintf("if%delse", loopNum))
			dname := Label(fmt.Sprintf("if%delsedone", loopNum))
			loopNum++
			body, err := evaluateCondition(s.Condition, context, iname)
			if err != nil {
				return nil, err
			}
			ops = append(ops, body...)

			body, err = compileBlock(s.Body, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, body...)
			ops = append(ops, JMP{dname})
			ops = append(ops, iname)
			if len(s.Else.Stmts) != 0 {
				body, err := compileBlock(s.Else, context)
				if err != nil {
					panic(err)
				}
				ops = append(ops, body...)
			}
			ops = append(ops, dname)
		case ast.MatchStmt:
			mname := Label(fmt.Sprintf("match%d", loopNum))
			loopNum++
			body, src, err := evaluateValue(s.Condition, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, body...)

			// Generate jump table
			for i := range s.Cases {
				body, dst, err := evaluateValue(s.Cases[i].Variable, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)

				ops = append(ops, JE{
					ConditionalJump{Label: Label(fmt.Sprintf("%vv%d", mname.Inline(), i)),
						Src: src,
						Dst: dst,
					},
				})
			}
			ops = append(ops, JMP{Label(fmt.Sprintf("%vdone", mname.Inline()))})

			// Generate bodies
			for i := range s.Cases {
				ops = append(ops, Label(fmt.Sprintf("%vv%d", mname.Inline(), i)))

				// Store the old values of variables for enum options that get
				// shadowed, and ensure they don't leak outside of the case
				oldVals := make(map[ast.VarWithType]Register)
				for k, v := range context.values {
					oldVals[k] = v
				}

				switch ev := s.Cases[i].Variable.(type) {
				case ast.EnumOption:
					// If the case was an EnumOption, it means the MatchStmt
					// variable was an enumerated data type. The index of the
					// original variable + i is the i'th parameter, so set
					// the appropriate LocalVariables in the context for
					// the case.
					val, ok := s.Condition.(ast.VarWithType)
					if !ok {
						panic("Unexpected pattern matching on non-variable")
					}
					vreg := context.Get(val)
					switch lv := vreg.(type) {
					case FuncArg:
						for j := range ev.Parameters {
							lv.Id += 1
							lv.Info = context.GetTypeInfo(s.Cases[i].LocalVariables[j].Type())
							context.SetLocalRegister(s.Cases[i].LocalVariables[j], lv)
						}

					case LocalValue:
						for j := range ev.Parameters {
							lv.Id += 1
							lv.Info = context.GetTypeInfo(s.Cases[i].LocalVariables[j].Type())
							context.SetLocalRegister(s.Cases[i].LocalVariables[j], lv)
						}
					default:
						panic(fmt.Sprintf("Expected enumeration to be a local variable or function argument: got %v", reflect.TypeOf(vreg)))
					}
				}
				body, err := compileBlock(s.Cases[i].Body, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, JMP{Label(fmt.Sprintf("%vdone", mname.Inline()))})
				context.values = oldVals

			}
			ops = append(ops, Label(fmt.Sprintf("%vdone", mname.Inline())))
		default:
			panic(fmt.Sprintf("Statement type not implemented: %v", reflect.TypeOf(s)))
		}
	}
	return ops, nil
}

// Evaluates a boolean condition. If the condition fails, jump to faillabel.
func evaluateCondition(val ast.BoolValue, context *variableLayout, faillabel Label) ([]Opcode, error) {
	var ops []Opcode
	switch c := val.(type) {
	case ast.GreaterComparison:
		body, r, err := evaluateValue(c.Left, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		body, r2, err := evaluateValue(c.Right, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		ops = append(ops, JLE{
			ConditionalJump{
				Label: faillabel,
				Src:   r,
				Dst:   r2,
			},
		})
		return ops, nil
	case ast.GreaterOrEqualComparison:
		body, r, err := evaluateValue(c.Left, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		body, r2, err := evaluateValue(c.Right, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)
		ops = append(ops, JL{
			ConditionalJump{Label: faillabel,
				Src: r,
				Dst: r2,
			},
		})
		return ops, nil
	case ast.LessThanComparison:
		body, r, err := evaluateValue(c.Left, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		body, r2, err := evaluateValue(c.Right, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)
		ops = append(ops, JGE{
			ConditionalJump{Label: faillabel,
				Src: r,
				Dst: r2,
			},
		})
		return ops, nil
	case ast.LessThanOrEqualComparison:
		body, r, err := evaluateValue(c.Left, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		body, r2, err := evaluateValue(c.Right, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)
		ops = append(ops, JG{
			ConditionalJump{Label: faillabel,
				Src: r,
				Dst: r2,
			},
		})
		return ops, nil
	case ast.EqualityComparison:
		body, r, err := evaluateValue(c.Left, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		body, r2, err := evaluateValue(c.Right, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)
		ops = append(ops, JNE{
			ConditionalJump{Label: faillabel,

				Src: r,
				Dst: r2,
			},
		})
		return ops, nil
	case ast.NotEqualsComparison:
		body, r, err := evaluateValue(c.Left, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)

		body, r2, err := evaluateValue(c.Right, context)
		if err != nil {
			return nil, err
		}
		ops = append(ops, body...)
		ops = append(ops, JE{
			ConditionalJump{Label: faillabel,
				Src: r,
				Dst: r2,
			},
		})
		return ops, nil
	case ast.BoolLiteral:
		if c {
			return ops, nil
		}
		ops = append(ops, JMP{faillabel})
		return ops, nil
	default:
		panic(fmt.Sprintf("Condition type not implemented: %v", reflect.TypeOf(c)))
	}
	return ops, nil
}

// Evaluates a value expression and returns the opcodes to evaluate it, and the
// register which contains the value evaluated.
func evaluateValue(val ast.Value, context *variableLayout) ([]Opcode, Register, error) {
	var ops []Opcode
	switch s := val.(type) {
	case ast.AdditionOperator:
		a := context.NextTempRegister()
		switch s.Left.(type) {
		case ast.VarWithType, ast.IntLiteral:
			ops = append(ops, MOV{
				Src: getRegister(s.Left, context),
				Dst: a,
			})
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator, ast.ArrayValue:
			body, r, err := evaluateValue(s.Left, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			ops = append(ops, MOV{
				Src: r,
				Dst: a,
			})
		default:
			panic(fmt.Sprintf("Unhandled left parameter in addition %v", reflect.TypeOf(s.Left)))
		}

		var r Register
		switch s.Right.(type) {
		case ast.IntLiteral, ast.VarWithType:
			// FIXME: This should validate type compatability
			r = getRegister(s.Right, context)
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator, ast.ArrayValue:
			body, r2, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			r = r2
		default:
			panic(fmt.Sprintf("Unhandled right parameter in addition: %v", reflect.TypeOf(s.Right)))

		}
		ops = append(ops, ADD{
			Src: r,
			Dst: a,
		})
		return ops, a, nil
	case ast.SubtractionOperator:
		a := context.NextTempRegister()
		switch s.Left.(type) {
		case ast.VarWithType, ast.IntLiteral:
			ops = append(ops, MOV{
				Src: getRegister(s.Left, context),
				Dst: a,
			})
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator, ast.ArrayValue:
			body, r, err := evaluateValue(s.Left, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			a = r
		default:
			panic(fmt.Sprintf("Unhandled left parameter in subtraction %v", reflect.TypeOf(s.Left)))
		}

		switch s.Right.(type) {
		case ast.IntLiteral, ast.VarWithType:
			ops = append(ops, SUB{
				Src: getRegister(s.Right, context),
				Dst: a,
			})
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator, ast.ArrayValue:
			body, r, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			ops = append(ops, SUB{
				Src: r,
				Dst: a,
			})
		default:
			panic("Unhandled right parameter in subtraction")
		}
		return ops, a, nil
	case ast.ModOperator:
		bodya, ra, err := evaluateValue(s.Left, context)
		if err != nil {
			return nil, nil, err
		}
		ops = append(ops, bodya...)

		bodyb, rb, err := evaluateValue(s.Right, context)
		if err != nil {
			return nil, nil, err
		}

		a := context.NextTempRegister()
		ops = append(ops, bodyb...)
		ops = append(ops, MOD{
			Left:  ra,
			Right: rb,
			Dst:   a,
		})
		return ops, a, nil
	case ast.MulOperator:
		a := context.NextTempRegister()
		var l, r Register
		switch s.Left.(type) {
		case ast.IntLiteral, ast.VarWithType:
			l = getRegister(s.Left, context)
		default:
			panic(fmt.Sprintf("Unhandled left parameter in mul %v", reflect.TypeOf(s.Left)))
		}

		switch s.Right.(type) {
		case ast.IntLiteral, ast.VarWithType:
			r = getRegister(s.Right, context)
		case ast.SubtractionOperator, ast.AdditionOperator, ast.MulOperator, ast.DivOperator, ast.ArrayValue:
			body, reg, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			r = reg
		}
		ops = append(ops, MUL{
			Left:  l,
			Right: r,
			Dst:   a,
		})
		return ops, a, nil
	case ast.DivOperator:
		a := context.NextTempRegister()
		var l, r Register
		switch s.Left.(type) {
		case ast.IntLiteral, ast.VarWithType:
			l = getRegister(s.Left, context)
		default:
			panic(fmt.Sprintf("Unhandled left parameter in div %v", reflect.TypeOf(s.Left)))
		}

		switch s.Right.(type) {
		case ast.IntLiteral, ast.VarWithType:
			r = getRegister(s.Right, context)
		case ast.SubtractionOperator, ast.AdditionOperator, ast.MulOperator, ast.DivOperator, ast.ArrayValue:
			body, reg, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			r = reg
		default:
			panic(fmt.Sprintf("Unhandled right parameter in div: %v", reflect.TypeOf(s.Right)))

		}
		ops = append(ops, DIV{
			Left:  l,
			Right: r,
			Dst:   a,
		})
		return ops, a, nil
	case ast.LessThanComparison, ast.LessThanOrEqualComparison,
		ast.EqualityComparison, ast.NotEqualsComparison,
		ast.GreaterComparison, ast.GreaterOrEqualComparison:
		cname := Label(fmt.Sprintf("comparison%d", loopNum))
		loopNum++

		a := context.NextTempRegister()
		bv, ok := s.(ast.BoolValue)
		if !ok {
			panic("Comparison operator doesn't implement BoolValue")
		}
		comp, err := evaluateCondition(bv, context, cname+"false")
		if err != nil {
			return nil, nil, err
		}
		ops = append(ops, comp...)
		ops = append(ops, MOV{
			Src: IntLiteral(1),
			Dst: a,
		})
		ops = append(ops, JMP{cname + "done"})
		ops = append(ops, Label(cname+"false"))
		ops = append(ops, MOV{
			Src: IntLiteral(0),
			Dst: a,
		})
		ops = append(ops, Label(cname+"done"))
		return ops, a, nil

	case ast.VarWithType, ast.IntLiteral, ast.BoolLiteral, ast.EnumOption, ast.StringLiteral:
		return nil, getRegister(s, context), nil
	case ast.ArrayValue:
		base := getRegister(s.Base, context)
		var a Register
		switch reg := base.(type) {
		case LocalValue:
			switch offset := s.Index.(type) {
			case ast.IntLiteral:
				// Special case to avoid the overhead of allocating/moving an extra register for
				// literals, we inline the multiplication..
				switch at := s.Base.Typ.(type) {
				case ast.ArrayType:
					reg.Info = context.GetTypeInfo(at.Base.Type())
				case ast.SliceType:
					reg.Id++
					reg.Info = context.GetTypeInfo(at.Base.Type())
				default:
					panic("Can only index into arrays or slices")
				}
				return nil, Offset{
					Offset: IntLiteral(int(offset) * reg.Size()),
					Base:   reg,
				}, nil
			default:
				// Evaluate the offset and look and store the value in a register.
				offsetops, offsetr, err := evaluateValue(s.Index, context)
				if err != nil {
					return nil, nil, err
				}
				ops = append(ops, offsetops...)

				// Convert the offset from index to byte offset
				realoffsetr := context.NextTempRegister()
				var tsize int
				switch at := s.Base.Typ.(type) {
				case ast.ArrayType:
					reg.Info = context.GetTypeInfo(at.Base.Type())
					tsize = reg.Info.Size
					if tsize == 0 {
						tsize = 8
					}

				case ast.SliceType:
					reg.Id++
					reg.Info = context.GetTypeInfo(at.Base.Type())
					tsize = reg.Info.Size
					if tsize == 0 {
						tsize = 8
					}
				default:
					panic("Can only index into arrays or slices")
				}

				ops = append(ops, MUL{
					Left:  IntLiteral(tsize),
					Right: offsetr,
					Dst:   realoffsetr,
				})
				a = Offset{
					Offset: realoffsetr,
					Base:   reg,
				}
			}
		case FuncArg:
			// Same as above, but Go type switches are stupid and force us to duplicate it.
			switch offset := s.Index.(type) {
			case ast.IntLiteral:
				// Special case to avoid the overhead of allocating/moving an extra register for
				// literals, we inline the multiplication..
				switch at := s.Base.Typ.(type) {
				case ast.ArrayType:
					reg.Info = context.GetTypeInfo(at.Base.Type())
				case ast.SliceType:
					reg.Id++
					reg.Info = context.GetTypeInfo(at.Base.Type())
				default:
					panic("Can only index into arrays or slices")
				}
				return nil, Offset{
					Offset: IntLiteral(int(offset) * reg.Size()),
					Base:   reg,
				}, nil
			default:
				// Evaluate the offset and look and store the value in a register.
				offsetops, offsetr, err := evaluateValue(s.Index, context)
				if err != nil {
					return nil, nil, err
				}
				ops = append(ops, offsetops...)

				// Convert the offset from index to byte offset
				realoffsetr := context.NextTempRegister() //LocalRegister(ast.VarWithType{"", ast.TypeLiteral("int"), false})
				var tsize int
				switch at := s.Base.Typ.(type) {
				case ast.ArrayType:
					reg.Info = context.GetTypeInfo(at.Base.Type())
					tsize = reg.Info.Size
					if tsize == 0 {
						tsize = 8
					}
				case ast.SliceType:
					reg.Id++
					reg.Info = context.GetTypeInfo(at.Base.Type())
					tsize = reg.Info.Size
					if tsize == 0 {
						tsize = 8
					}
				default:
					panic("Can only index into arrays or slices")
				}

				ops = append(ops, MUL{
					Left:  IntLiteral(tsize),
					Right: offsetr,
					Dst:   realoffsetr,
				})
				a = Offset{
					Offset: realoffsetr,
					Base:   reg,
				}
			}
		default:
			panic(fmt.Sprintf("Array was neither allocated in function nor passed as parameter: %v", reflect.TypeOf(base)))
		}
		return ops, a, nil
	default:
		panic(fmt.Errorf("Unhandled value type: %v", reflect.TypeOf(s)))
	}
}
