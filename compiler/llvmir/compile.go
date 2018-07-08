package llvmir

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/driusan/lang/parser/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type EnumIndexInfo struct {
	Index int
	Defn  ast.EnumTypeDefn
}

type EnumMap map[string]EnumIndexInfo
type Register struct {
	Value      value.Value
	FromMemory bool
}

func Generate(r io.Reader) (*ir.Module, error) {
	nodes, _, _, err := ast.ParseFromReader(r)
	if err != nil {
		return nil, err
	}
	m := ir.NewModule()

	ctx := NewContext(m)

	// Pass 1: Generate all the functions (and types)?
	for _, n := range nodes {
		switch f := n.(type) {
		case ast.FuncDecl:
			var ret types.Type
			switch len(f.Return) {
			case 0:
				ret = types.Void
			case 1:
				ret = getType(f.Return[0].Typ)
			default:
				panic("Returning values not implemented")
			}
			fnc := m.NewFunction(f.Name, ret)
			ctx.Funcs[f.Name] = FuncDef{fnc, f}
		case ast.EnumTypeDefn:
			for i, v := range f.Options {
				ctx.enumValues[v.Constructor] = EnumIndexInfo{i, f}
			}
		}
	}

	// Pass 2: compile the functions now that func calls can be
	for _, n := range nodes {
		switch f := n.(type) {
		case ast.FuncDecl:
			fnc := ctx.Funcs[f.Name]
			ctx.Variables = make(map[ast.VarWithType]Register)
			for _, arg := range f.Args {
				var p *types.Param
				if _, ok := arg.Type().(ast.ArrayType); ok {
					// We pass all arrays as pointers so that GEP works.
					arg.Reference = true
				}
				if arg.Reference {
					p = fnc.NewParam(string(arg.Name), types.NewPointer(getType(arg.Type())))
				} else {
					p = fnc.NewParam(string(arg.Name), getType(arg.Type()))
				}
				ctx.SetVar(arg, Register{
					Value:      p,
					FromMemory: arg.Reference,
				})
			}
			ctx.curfunc = fnc.Function
			ctx.curfuncdef = f
			body := fnc.NewBlock(f.Name + "start")
			ctx.curblock = body

			compileBlock(ctx, f.Body)

			if fnc.Sig.Ret == types.Void {
				// If we fall off the end of a void function, we
				// need to ensure there's a return
				ctx.curblock.NewRet(nil)
			}
		}
	}
	return m, nil
}

func compileBlock(ctx *Context, code ast.BlockStmt) {
	for _, stmt := range code.Stmts {
		output := ctx.curblock
		switch v := stmt.(type) {
		case ast.FuncCall:
			callFunc(ctx, v)
		case ast.ReturnStmt:
			if v.Val == nil {
				output.NewRet(nil)
			} else {
				rett := ctx.curfuncdef.Return[0].Type()
				rv := evalValue(ctx, output, v.Val, false, rett, false)
				output.NewRet(rv)
			}
		case ast.LetStmt:
			initializeLetStmt(ctx, v)
		case ast.MutStmt:
			var t ast.Type = v.Var.Type()
			switch ty := t.(type) {
			case ast.SliceType:
				t = ty.Base
			case ast.ArrayType:
				t = ty.Base
			}
			// Mutable statements allocate a memory address for the statement to
			// use as the canonical place to store the value. Initializing and
			// assigning store to the location, while reading from it loads from
			// it. The LLVM optimizer can take care of turning it from an Alloca
			// load/store to a phi function with the -mem2reg optimization, so
			// we don't bother trying to do anything fancy ourselves.
			val := evalValue(ctx, output, v.InitialValue, false, t, false)
			loc := output.NewAlloca(getType(v.Var.Type()))
			// Giving a friendly name breaks when the same function has multiple
			// mutable statements in different non-shadowed contexted in the same
			// function, because it needs to be in SSA form.
			// So we can't do "loc.SetName(v.Var.Name.String())"
			switch v.Var.Type().(type) {
			case ast.ArrayType:
				l := output.NewLoad(val)
				output.NewStore(l, loc)
			case ast.SliceType:
				var finalval value.Value
				switch v := v.InitialValue.(type) {
				case ast.ArrayLiteral:
					// If it's an array literal, get the base from element 0
					addr := ctx.curblock.NewGetElementPtr(val, constant.NewInt(0, types.I64), constant.NewInt(0, types.I64))
					ln := constant.NewInt(int64(len(v)), types.I64)

					s := constant.NewStruct(constant.NewUndef(types.NewPointer(getType(t))), ln)
					finalval = ctx.curblock.NewInsertValue(s, addr, []int64{0})
				case ast.VarWithType:
					// If it's another variable, get the base and length by extracting
					// it from the other slice.
					addr := ctx.curblock.NewExtractValue(val, []int64{0})
					ln := ctx.curblock.NewExtractValue(val, []int64{1})

					// convert it to a constant by using NewInsertValue
					s := constant.NewStruct(constant.NewUndef(types.NewPointer(getType(t))), constant.NewUndef(types.I64))
					withptr := ctx.curblock.NewInsertValue(s, addr, []int64{0})
					finalval = ctx.curblock.NewInsertValue(withptr, ln, []int64{1})
				case ast.Slice:
					finalval = val
				default:
					panic(fmt.Sprintf("Unhandled type: %v", reflect.TypeOf(v)))
				}
				output.NewStore(finalval, loc)
			default:
				output.NewStore(val, loc)
			}

			ctx.SetVar(v.Var, Register{
				Value:      loc,
				FromMemory: true,
			})
		case ast.AssignmentOperator:
			switch vr := v.Variable.(type) {
			case ast.VarWithType:
				reg := ctx.Variables[vr]
				val := evalValue(ctx, ctx.curblock, v.Value, false, nil, false)

				output.NewStore(val, reg.Value)
			case ast.ArrayValue:
				var t ast.Type = vr.Type()
				switch ty := t.(type) {
				case ast.ArrayType:
					t = ty.Base
				case ast.SliceType:
					t = ty.Base
				}
				ptr := evalValue(ctx, ctx.curblock, vr, true, nil, false)
				val := evalValue(ctx, ctx.curblock, v.Value, false, t, false)
				output.NewStore(val, ptr)

			default:
				panic(fmt.Sprintf("Unhandled assignable type %v", reflect.TypeOf(vr)))
			}
		case ast.WhileLoop:
			// Generate the new blocks for the condition, the loop, and the end label
			oblock := ctx.curblock
			ovars := ctx.cloneVars()
			loopinit := ctx.curfunc.NewBlock(fmt.Sprintf("while%dinit", ctx.loopNum))
			loopcond := ctx.curfunc.NewBlock(fmt.Sprintf("while%dcond", ctx.loopNum))
			loopbody := ctx.curfunc.NewBlock(fmt.Sprintf("while%dbody", ctx.loopNum))
			loopend := ctx.curfunc.NewBlock(fmt.Sprintf("while%dend", ctx.loopNum))
			ctx.loopNum++

			// We start by initializing any let conditions, then unconditionally checking
			// the condition.
			oblock.NewBr(loopinit)

			// The initialization unconditionally terminates by branching to checking the condition
			loopinit.NewBr(loopcond)

			// The body unconditionally terminates by branching to check the condition.
			loopbody.NewBr(loopcond)

			// Declare alloca and initialize them for any let conditionals
			ctx.curblock = loopinit
			val := evalValue(ctx, loopinit, v.Condition, false, nil, true)

			loopinit.NewCondBr(val, loopbody, loopend)
			// Hack to ensure the condition evaluation gets put into both init and cond,
			// we re-eval cond
			ctx.curblock = loopcond
			val2 := evalValue(ctx, ctx.curblock, v.Condition, true, nil, false)
			loopcond.NewCondBr(val2, loopbody, loopend)

			ctx.curblock = loopbody
			compileBlock(ctx, v.Body)
			if ctx.curblock.Term == nil {
				ctx.curblock.NewBr(loopcond)
			}
			ctx.curblock = loopend
			ctx.Variables = ovars
		case ast.IfStmt:
			oblock := ctx.curblock
			// Back up the variables to undo any thing was shadowed at the end
			ovars := ctx.cloneVars()

			ifcond := ctx.curfunc.NewBlock(fmt.Sprintf("if%dcond", ctx.ifNum))
			ifbody := ctx.curfunc.NewBlock(fmt.Sprintf("if%dbody", ctx.ifNum))
			ifelse := ctx.curfunc.NewBlock(fmt.Sprintf("if%delse", ctx.ifNum))
			ifend := ctx.curfunc.NewBlock(fmt.Sprintf("if%dend", ctx.ifNum))
			ctx.ifNum++

			if oblock.Term == nil {
				oblock.NewBr(ifcond)
			} else {
				switch t := oblock.Term.(type) {
				case *ir.TermCondBr, *ir.TermBr:
					if strings.HasPrefix(t.GetParent().GetName(), "while") ||
						strings.HasPrefix(t.GetParent().GetName(), "if") {
						// If this is an else if or an if inside of a while loop, steal
						// the term from the parent and make the parent branch here
						// instead.
						//
						// (FIXME: There should probably be a more robust way of doing this.)
						ifend.SetTerm(t)
						oblock.NewBr(ifcond)
					} else {
						panic("Attempting to jump to if statement from already terminated block")
					}
				default:
					panic(fmt.Sprintf("Attempting to jump to if statement from already terminated block %v", reflect.TypeOf(t)))
				}
			}
			ifbody.NewBr(ifend)
			ifelse.NewBr(ifend)

			ctx.curblock = ifcond
			val := evalValue(ctx, ctx.curblock, v.Condition, false, nil, false)
			ifcond.NewCondBr(val, ifbody, ifelse)

			ctx.curblock = ifbody
			compileBlock(ctx, v.Body)

			ctx.curblock = ifelse
			compileBlock(ctx, v.Else)

			ctx.curblock = ifend
			ctx.Variables = ovars
		case ast.MatchStmt:
			oblock := ctx.curblock

			ovars := ctx.cloneVars()
			val := evalValue(ctx, ctx.curblock, v.Condition, false, nil, false)
			// Cache the matchNum, so that it doesn't change under us if matches
			// are embedded in cases
			mn := ctx.matchNum
			matchend := ctx.curfunc.NewBlock(fmt.Sprintf("match%dend", mn))
			matchend.NewUnreachable()
			ctx.matchNum++

			switch c := v.Condition.(type) {
			case ast.VarWithType:
				var cases []*ir.Case
				switch e := c.Typ.(type) {
				case ast.EnumTypeDefn:
					// Destructuring an enum type
					// We switch on the variant, not the value
					//variant := val
					if len(e.Parameters) == 0 {
						// An enum with no components can act as a switch
						for i, c := range v.Cases {
							caseib := ctx.curfunc.NewBlock(fmt.Sprintf("match%dcase%d", mn, i))
							caseib.NewBr(matchend)

							ctx.curblock = caseib
							compileBlock(ctx, c.Body)
							casei := ir.NewCase(constant.NewInt(int64(i), val.Type()), caseib)
							cases = append(cases, casei)
						}
						oblock.NewSwitch(val, matchend, cases...)
					} else {
						// An enum with subcomponenents needs to be destructured as an if/else
						// chain after extracting the variant. This works similarly to the default
						// case below, but needs to extract the value for comparison first.
						cases := []struct {
							Cond *ir.BasicBlock
							Body *ir.BasicBlock
						}{}

						// Generate all the blocks first, so that the terminators can be set
						for i := range v.Cases {
							var cs struct{ Cond, Body *ir.BasicBlock }
							cs.Cond = ctx.curfunc.NewBlock(fmt.Sprintf("match%dcomp%d", mn, i))
							cs.Body = ctx.curfunc.NewBlock(fmt.Sprintf("match%dbody%d", mn, i))
							cases = append(cases, cs)
						}

						val := evalValue(ctx, ctx.curblock, v.Condition, false, nil, false)
						variant := ctx.curblock.NewExtractValue(val, []int64{0})
						// Hook up the condition with a branch to the body if it's matched, and
						// the next case condition otherwise.
						//
						// Jump to the end if we get to the end of any body.
						for i, c := range v.Cases {
							if i == 0 {
								ctx.curblock.NewBr(cases[0].Cond)
							}
							ctx.curblock = cases[i].Cond
							v2 := evalValue(ctx, ctx.curblock, c.Variable, false, nil, false)

							eq := ctx.curblock.NewICmp(ir.IntEQ, variant, v2)
							// If it's the last one and the match fails, go to the end, otherwise
							// go to the next comparison
							if i == len(cases)-1 {
								cases[i].Cond.NewCondBr(eq, cases[i].Body, matchend)
							} else {
								cases[i].Cond.NewCondBr(eq, cases[i].Body, cases[i+1].Cond)
							}
							cases[i].Body.NewBr(matchend)

							ctx.curblock = cases[i].Body
							for j, vr := range c.LocalVariables {
								// j = variant, j+1 = this local variable
								l := ctx.curblock.NewExtractValue(val, []int64{int64(j + 1)})
								ctx.SetVar(vr, Register{
									Value:      l,
									FromMemory: false,
								})
							}
							compileBlock(ctx, c.Body)
						}
					}
				case ast.SumType:
					ovars := ctx.cloneVars()
					// If it's a non-enum sum type, we still need to destructure and ensure
					// that the var locally scoped has the right type in the context.
					variant := ctx.curblock.NewExtractValue(val, []int64{int64(0)})
					valpiece := ctx.curblock.NewExtractValue(val, []int64{int64(1)})
					for i, cs := range v.Cases {
						caseib := ctx.curfunc.NewBlock(fmt.Sprintf("match%dcase%d", mn, i))
						caseib.NewBr(matchend)

						ctx.curblock = caseib
						subval := ctx.curblock.NewExtractValue(valpiece, []int64{int64(i)})
						c.Typ = cs.Variable.Type()
						ctx.SetVar(c, Register{subval, false})
						compileBlock(ctx, cs.Body)
						casei := ir.NewCase(constant.NewInt(int64(i), types.I64), caseib)

						cases = append(cases, casei)
					}
					oblock.NewSwitch(variant, matchend, cases...)
					ctx.Variables = ovars
				default:
					// Matching on a value
					for i, c := range v.Cases {
						caseib := ctx.curfunc.NewBlock(fmt.Sprintf("match%dcase%d", mn, i))
						caseib.NewBr(matchend)

						ctx.curblock = caseib
						compileBlock(ctx, c.Body)
						casei := ir.NewCase(constant.NewInt(int64(i), val.Type()), caseib)

						cases = append(cases, casei)
					}
					oblock.NewSwitch(val, matchend, cases...)
				}
			default:
				// If it's not match x {} for a variable x (ie it's a complex expression, or no
				// condition at all, or some other form, it doesn't map to an LLVM IR switch
				// statement, so instead we treat it as an if/else if chain implicitly
				// comparing the match value to the case value
				cases := []struct {
					Cond *ir.BasicBlock
					Body *ir.BasicBlock
				}{}

				// Generate all the blocks first, so that the terminators can be set
				for i := range v.Cases {
					var cs struct{ Cond, Body *ir.BasicBlock }
					cs.Cond = ctx.curfunc.NewBlock(fmt.Sprintf("match%dcomp%d", mn, i))
					cs.Body = ctx.curfunc.NewBlock(fmt.Sprintf("match%dbody%d", mn, i))
					cases = append(cases, cs)
				}

				// Hook up the condition with a branch to the body if it's matched, and
				// the next case condition otherwise.
				//
				// Jump to the end if we get to the end of any body.
				for i, c := range v.Cases {
					if i == 0 {
						ctx.curblock.NewBr(cases[0].Cond)
					}
					ctx.curblock = cases[i].Cond
					v2 := evalValue(ctx, ctx.curblock, c.Variable, false, nil, false)

					eq := ctx.curblock.NewICmp(ir.IntEQ, val, v2)
					// If it's the last one and the match fails, go to the end, otherwise
					// go to the next comparison
					if i == len(cases)-1 {
						cases[i].Cond.NewCondBr(eq, cases[i].Body, matchend)
					} else {
						cases[i].Cond.NewCondBr(eq, cases[i].Body, cases[i+1].Cond)
					}
					cases[i].Body.NewBr(matchend)

					ctx.curblock = cases[i].Body
					compileBlock(ctx, c.Body)
				}
			}
			ctx.curblock = matchend
			ctx.Variables = ovars
		case ast.Assertion:
			oblock := ctx.curblock
			// Back up the variables to undo any thing was shadowed at the end

			assertcond := ctx.curfunc.NewBlock(fmt.Sprintf("assert%dcond", ctx.assertNum))
			assertfail := ctx.curfunc.NewBlock(fmt.Sprintf("assert%dfail", ctx.assertNum))
			assertend := ctx.curfunc.NewBlock(fmt.Sprintf("assert%dend", ctx.assertNum))
			ctx.assertNum++

			assertend.Term = oblock.Term
			oblock.NewBr(assertcond)

			ctx.curblock = assertcond
			ty := inferType(v.Predicate)
			val := evalValue(ctx, ctx.curblock, v.Predicate, false, ty, false)

			ctx.curblock = assertfail
			msg := fmt.Sprintf("assertion %v failed", v.PrettyPrint(0))
			if v.Message != "" {
				msg += ": " + string(v.Message)
			}

			callFunc(ctx, ast.FuncCall{
				Name:     "Write",
				UserArgs: []ast.Value{ast.IntLiteral(2), ast.StringLiteral(msg)},
			})
			callFunc(ctx, ast.FuncCall{
				Name:     "Exit",
				UserArgs: []ast.Value{ast.IntLiteral(1)},
			})
			assertfail.NewUnreachable()

			assertcond.NewCondBr(val, assertend, assertfail)
			ctx.curblock = assertend
		default:
			panic(fmt.Sprintf("Unhandled statement type: %v", reflect.TypeOf(v)))
		}
	}
}

// Converts ast.Value to an ir.Value, if necessary adding any intermediate
// operations to fnc.
//
// If addronly is true, then evalValue will only return the register containing
// the address and not do the load for any array operations.
func evalValue(ctx *Context, fnc *ir.BasicBlock, val ast.Value, addronly bool, arrayType ast.Type, forinit bool) value.Value {
	switch t := arrayType.(type) {
	case ast.SumType:
		// If it's a sum type that this is being evaluated for, we need to add the variant and wrap it
		// in the right structure.
		literal := evalValue(ctx, fnc, val, addronly, nil, forinit)
		ltype := literal.Type()
		if ltype.Equal(getType(t)) {
			return literal
		}

		// Variant, pointer to variant being held
		var typeOptions []constant.Constant
		found := int64(-1)
		for i, st := range t {
			typePiece := getType(st)
			// Fill in everything as undef
			typeOptions = append(typeOptions, constant.NewUndef(typePiece))

			if found < 0 && ltype.Equal(typePiece) {
				found = int64(i)
			}
		}
		if found < 0 {
			panic("No compatible types for sum type")
		}
		substruct := constant.NewStruct(typeOptions...)
		iv := ctx.curblock.NewInsertValue(substruct, literal, []int64{found})

		undef := constant.NewStruct(constant.NewInt(found, types.I64), substruct)
		withval := ctx.curblock.NewInsertValue(undef, iv, []int64{1})
		return withval
	}

	switch v := val.(type) {
	case ast.StringLiteral:
		r, l := ctx.GetStringLiteral(string(v))
		addr := fnc.NewGetElementPtr(r, constant.NewInt(0, types.I64), constant.NewInt(0, types.I64))
		s := constant.NewStruct(constant.NewUndef(types.NewPointer(types.I8)), constant.NewInt(l, types.I64))
		iv := ctx.curblock.NewInsertValue(s, addr, []int64{0})
		return iv
	case ast.IntLiteral:
		if arrayType != nil {
			return constant.NewInt(int64(v), getType(arrayType))
		}
		return constant.NewInt(int64(v), types.I64)
	case ast.BoolLiteral:
		if v {
			return constant.NewInt(1, types.I1)
		} else {
			return constant.NewInt(0, types.I1)
		}
	case ast.EnumValue:
		etd := ctx.GetEnumTypeDefn(v.Constructor.Constructor)
		pieces := []constant.Constant{constant.NewInt(int64(ctx.GetEnumIndex(v.Constructor.Constructor)), types.I64)}
		for _, v := range v.Parameters {
			subpiece := evalValue(ctx, fnc, v, false, v.Type(), forinit)
			switch val := subpiece.(type) {
			case *constant.Int:
				pieces = append(pieces, val)
			default:
				panic(fmt.Sprintf("Unhandled subtype type: %v", reflect.TypeOf(val)))
			}
		}
		// If the value didn't use up all the parameters from the
		// EnumOption that it's a value for, fill it with i64 0 so
		// that all types of this EnumOption type are the same length
		expectedLen := etd.ExpectedParams
		for i := len(pieces); i <= expectedLen; i++ {
			pieces = append(pieces, constant.NewInt(0, types.I64))
		}
		if len(pieces) == 1 {
			// If there were no parameters for the variant, just return the index directly
			return pieces[0]
		}
		// If there were pieces, return it as a struct
		return constant.NewStruct(pieces...)
	case ast.VarWithType:
		vr := ctx.GetVariable(v)
		if !vr.FromMemory || addronly {
			return vr.Value
		}

		// mut var needs to be loaded to dereference the value
		// (unless explicitly told we only want the addr)
		load := fnc.NewLoad(vr.Value)
		return load
	case ast.FuncCall:
		return callFunc(ctx, v)
	case ast.SubtractionOperator:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		return ctx.curblock.NewSub(l, r)
	case ast.AdditionOperator:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		return ctx.curblock.NewAdd(l, r)
	case ast.MulOperator:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		return ctx.curblock.NewMul(l, r)
	case ast.DivOperator:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		// FIXME: This should take signedness into account
		return ctx.curblock.NewSDiv(l, r)
	case ast.ModOperator:
		// FIXME: srem is the remainder, not the modulo.
		// FIXME: This should take signedness into account
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		return ctx.curblock.NewSRem(l, r)
	case ast.EqualityComparison:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		return ctx.curblock.NewICmp(ir.IntEQ, l, r)
	case ast.NotEqualsComparison:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		return ctx.curblock.NewICmp(ir.IntNE, l, r)
	case ast.GreaterComparison:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		// FIXME: This should take into account signed/unsigned
		return ctx.curblock.NewICmp(ir.IntSGT, l, r)
	case ast.GreaterOrEqualComparison:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		// FIXME: This should take into account signed/unsigned
		return ctx.curblock.NewICmp(ir.IntSGE, l, r)
	case ast.LessThanComparison:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		// FIXME: This should take into account signed/unsigned
		return ctx.curblock.NewICmp(ir.IntSLT, l, r)
	case ast.LessThanOrEqualComparison:
		l := evalValue(ctx, ctx.curblock, v.Left, false, arrayType, forinit)
		r := evalValue(ctx, ctx.curblock, v.Right, false, arrayType, forinit)
		// FIXME: This should take into account signed/unsigned
		return ctx.curblock.NewICmp(ir.IntSLE, l, r)
	case ast.EnumOption:
		o := ctx.GetEnumIndex(v.Constructor)
		return constant.NewInt(int64(o), types.I64)
	case ast.ArrayLiteral:
		var arr []constant.Constant
		for _, e := range v {
			sv := evalValue(ctx, fnc, e, false, arrayType, forinit)
			switch val := sv.(type) {
			case *constant.Int:
				arr = append(arr, val)
			case *ir.InstInsertValue:
				// This is likely a string. We need to insert it as an undef
				// in order to declare the new array, and then insert the
				// values after, because the insertvalue isn't a constant.
				arr = append(arr, constant.NewUndef(val.Type()))
			default:
				panic(fmt.Sprintf("Unhandled subtype type: %v", reflect.TypeOf(val)))
			}
		}
		arrc := constant.NewArray(arr...)

		var val value.Value = arrc
		// Now that we've
		for i, e := range v {
			sv := evalValue(ctx, fnc, e, false, arrayType, forinit)
			switch vl := sv.(type) {
			case *constant.Int:
				// Do nothing
			case *ir.InstInsertValue:
				val = ctx.curblock.NewInsertValue(val, vl, []int64{int64(i)})
			default:
				panic(fmt.Sprintf("Unhandled subtype type: %v", reflect.TypeOf(val)))
			}
		}
		a := ctx.curblock.NewAlloca(arrc.Type())
		if addronly {
			return a
		}
		ctx.curblock.NewStore(val, a)
		return a
	case ast.ArrayValue:
		switch bt := v.Base.Type().(type) {
		case ast.ArrayType:
			idx := evalValue(ctx, fnc, v.Index, false, arrayType, forinit)
			base := ctx.GetVariable(v.Base).Value
			elptr := ctx.curblock.NewGetElementPtr(
				base,
				constant.NewInt(0, types.I64),
				idx,
			)
			if addronly {
				return elptr
			}
			return ctx.curblock.NewLoad(elptr)
		case ast.SliceType:
			// The index
			idx := evalValue(ctx, fnc, v.Index, false, arrayType, forinit)

			// The base that we're indexing from
			var extbase *ir.InstExtractValue
			base := ctx.GetVariable(v.Base).Value
			switch base.(type) {
			case *ir.InstAlloca:
				// mutable statement, needs to load through ptr before extracting
				ld := ctx.curblock.NewLoad(base)
				extbase = ctx.curblock.NewExtractValue(ld, []int64{0})
			default:
				extbase = ctx.curblock.NewExtractValue(base, []int64{0})
			}

			// Load the address
			elptr := ctx.curblock.NewGetElementPtr(
				extbase,
				idx,
			)
			if addronly {
				return elptr
			}

			// Dereference the value
			return ctx.curblock.NewLoad(elptr)
		default:
			panic(fmt.Sprintf("Unhandled indexable type: %v", reflect.TypeOf(bt)))
		}
	case ast.Cast:
		irv := evalValue(ctx, fnc, v.Val, addronly, arrayType, forinit)
		if len(v.Val.Type().Components()) == 1 {
			return ctx.convValue(irv, v.Val.Type(), v.Typ)
		}
		return irv
	case ast.Brackets:
		return evalValue(ctx, fnc, v.Val, addronly, arrayType, forinit)
	case ast.LetStmt:
		// We need to store conditional lets into a memory address as if the
		// were a mutable so that statements such as "while (let i = i + 1) < 3"
		// will work
		loc, ok := ctx.GetVariableSafe(v.Var)
		if forinit {
			// let statement in a loop conditional
			newloc := ctx.curblock.NewAlloca(getType(v.Var.Type()))
			ctx.SetVar(v.Var, Register{
				Value:      newloc,
				FromMemory: true,
			})
			val := evalValue(ctx, fnc, v.Val, false, v.Var.Type(), false)
			ctx.curblock.NewStore(val, newloc)
			return val
		} else if ok && loc.FromMemory {
			val := evalValue(ctx, fnc, v.Val, false, v.Var.Type(), forinit)
			ctx.curblock.NewStore(val, loc.Value)
			return val
		} else {
			val := evalValue(ctx, fnc, v.Val, false, v.Var.Type(), forinit)
			// let statement in an if
			ctx.SetVar(v.Var, Register{
				Value:      val,
				FromMemory: false,
			})
			return val
		}
	case ast.TupleValue:
		var subs []constant.Constant
		ty := v.Type().(ast.TupleType)
		var subvalues []value.Value
		for i, aval := range v {
			val := evalValue(ctx, fnc, aval, false, ty[i].Type(), forinit)
			subvalues = append(subvalues, val)
			if c, ok := val.(constant.Constant); ok {
				subs = append(subs, c)
			} else {
				subs = append(subs, constant.NewUndef(getType(ty[i].Type())))
			}
		}
		cstr := constant.NewStruct(subs...)

		// Now insert the ones that weren't constants
		var rv value.Value = cstr
		for i, val := range subvalues {
			if _, ok := val.(constant.Constant); !ok {
				rv = ctx.curblock.NewInsertValue(rv, val, []int64{int64(i)})
			}
		}
		return rv
	case ast.Slice:
		base := evalValue(ctx, fnc, v.Base, true, arrayType, forinit)
		cnt := constant.NewStruct(constant.NewUndef(types.NewPointer(getType(v.Base.Type()))), constant.NewInt(int64(v.Size), types.I64))
		switch p := base.(type) {
		case *ir.InstGetElementPtr:
			iv := ctx.curblock.NewInsertValue(cnt, p, []int64{0})
			return iv

		}
		panic(fmt.Sprintf("Unhandled base type for slice: %v", reflect.TypeOf(base)))
	default:
		panic(fmt.Sprintf("Unhandled ast.Value type %v", reflect.TypeOf(v)))
	}
}

func (ctx *Context) GetStringLiteral(val string) (value.Value, int64) {
	val = strings.Replace(val, `\n`, "\n", -1)
	if v, ok := ctx.StringLiterals[val]; ok {
		return v, int64(len([]byte(val)))
	}
	var chars []constant.Constant
	for _, c := range val {
		chars = append(chars, constant.NewInt(int64(c), types.I8))
	}
	arr := constant.NewArray(chars...)
	arr.CharArray = true
	r := ctx.module.NewGlobalDef(fmt.Sprintf(".str%d", len(ctx.StringLiterals)), arr)
	ctx.StringLiterals[val] = r
	return r, int64(len([]byte(val)))
}

func (ctx *Context) convValue(val value.Value, from ast.Type, to ast.Type) value.Value {
	fromc := from.Components()
	toc := to.Components()
	if len(fromc) != 1 || len(toc) != 1 {
		if getType(from).String() == getType(to).String() {
			// If they're the exact same type in the IR, bitcast it to the same thing
			// just so that there's an instruction to reference it.
			return ctx.curblock.NewBitCast(val, getType(to))
		}
		panic("Only single component values can be converted (for now.)")
	}
	if len(fromc) != len(toc) {
		panic("Can only convert values with the same number of components")
	}
	dst := toc[0].Info()
	src := fromc[0].Info()
	if src.Size == 0 {
		src.Size = 8
	}
	if dst.Size == 0 {
		dst.Size = 8
	}
	switch {
	case dst.Size < src.Size:
		return ctx.curblock.NewTrunc(val, getType(to))
	case dst.Size == src.Size:
		return ctx.curblock.NewBitCast(val, getType(to))
	case dst.Size > src.Size:
		if dst.Signed {
			return ctx.curblock.NewSExt(val, getType(to))
		} else {
			return ctx.curblock.NewZExt(val, getType(to))
		}
	default:
		panic("Unhandled conversion")
	}
}

func getType(typ ast.Type) types.Type {
	ti := typ.Info()
	if typ == ast.TypeLiteral("bool") {
		return types.I1
	}
	if typ.TypeName() == "string" {
		return types.NewStruct(types.NewPointer(types.I8), types.I64)
	}
	components := typ.Components()
	if len(components) == 1 {
		switch ti.Size {
		case 0, 8:
			return types.I64
		case 4:
			return types.I32
		case 2:
			return types.I16
		case 1:
			return types.I8
		default:
			panic("Unhandled type size")
		}
	}
	switch t := typ.(type) {
	case ast.ArrayType:
		return types.NewArray(getType(t.Base), int64(t.Size))
	case ast.EnumTypeDefn:
		pieces := make([]types.Type, 0, len(components))
		for _, c := range components {
			pieces = append(pieces, getType(c))
		}
		return types.NewStruct(pieces...)
	case ast.SliceType:
		return types.NewStruct(types.NewPointer(getType(t.Base)), types.I64)
	case ast.SumType:
		// We're stupid with how we model sum types. We use a product type but
		// only set a single field in it, because LLVM IR doesn't make it very
		// easy to model unions with its type system, and this is mostly to
		// get the tests passing. This can be made more efficient later.
		var typeOptions []types.Type
		for _, st := range t {
			typeOptions = append(typeOptions, getType(st))
		}
		return types.NewStruct(types.I64, types.NewStruct(typeOptions...))
	case ast.TupleType:
		var typeOptions []types.Type
		for _, st := range t {
			typeOptions = append(typeOptions, getType(st.Type()))
		}
		return types.NewStruct(typeOptions...)
	default:
		panic(fmt.Sprintf("Unhandled type %v", reflect.TypeOf(t)))
	}
}

func getSliceLength(base ast.Value) constant.Constant {
	switch t := base.(type) {
	case ast.ArrayLiteral:
		return constant.NewInt(int64(len(t)), types.I64)
	default:
		panic(fmt.Sprintf("Don't know how to extract slice length from %v", reflect.TypeOf(t)))
	}
}

func initializeLetStmt(ctx *Context, v ast.LetStmt) {
	// Let statements go into a register and never need to be loaded
	// from memory, because they're immutable
	switch t := v.Var.Type().(type) {
	case ast.SliceType:
		val := evalValue(ctx, ctx.curblock, v.Val, false, t.Base, false)
		var finalval value.Value
		switch v := v.Val.(type) {
		case ast.ArrayLiteral:
			// If it's an array literal, get the base from element 0
			addr := ctx.curblock.NewGetElementPtr(val, constant.NewInt(0, types.I64), constant.NewInt(0, types.I64))
			ln := constant.NewInt(int64(len(v)), types.I64)

			s := constant.NewStruct(constant.NewUndef(types.NewPointer(getType(t.Base))), ln)
			finalval = ctx.curblock.NewInsertValue(s, addr, []int64{0})
		case ast.VarWithType:
			// If it's another variable, get the base and length by extracting
			// it from the other slice.
			addr := ctx.curblock.NewExtractValue(val, []int64{0})
			ln := ctx.curblock.NewExtractValue(val, []int64{1})

			// convert it to a constant by using NewInsertValue
			s := constant.NewStruct(constant.NewUndef(types.NewPointer(getType(t.Base))), constant.NewUndef(types.I64))
			withptr := ctx.curblock.NewInsertValue(s, addr, []int64{0})
			finalval = ctx.curblock.NewInsertValue(withptr, ln, []int64{1})
		case ast.Slice:
			finalval = val
		default:
			panic(fmt.Sprintf("Unhandled type: %v", reflect.TypeOf(v)))
		}
		ctx.SetVar(v.Var, Register{
			Value:      finalval,
			FromMemory: false,
		})
	case ast.ArrayType:
		val := evalValue(ctx, ctx.curblock, v.Val, false, t.Base, false)
		ctx.SetVar(v.Var, Register{
			Value:      val,
			FromMemory: false,
		})
	case ast.TupleType:
		val := evalValue(ctx, ctx.curblock, v.Val, false, t, false)
		ctx.SetVar(v.Var, Register{
			Value:      val,
			FromMemory: false,
		})
		basename := v.Var.Name
		ty := v.Var.Type().(ast.TupleType)
		for i := range ty {
			// Set variables for accessors too.
			v.Var.Name = basename + "." + ty[i].Name
			v.Var.Typ = ty[i].Type()
			ev := ctx.curblock.NewExtractValue(val, []int64{int64(i)})
			ctx.SetVar(v.Var, Register{
				Value:      ev,
				FromMemory: false,
			})
		}
	case ast.UserType:
		// Strip user types to their underlying type and try again.
		ot := v.Var.Typ
		v.Var.Typ = t.Typ

		initializeLetStmt(ctx, v)

		// The Variables should be stored with the user type, not the
		// underlying type.
		reg := ctx.GetVariable(v.Var)
		delete(ctx.Variables, hashableHack(v.Var))
		v.Var.Typ = ot

		ctx.SetVar(v.Var, reg)
	default:
		val := evalValue(ctx, ctx.curblock, v.Val, false, t, false)
		ctx.SetVar(v.Var, Register{
			Value:      val,
			FromMemory: false,
		})
	}
}

func callFunc(ctx *Context, v ast.FuncCall) value.Value {
	if v.Name == "len" {
		v := evalValue(ctx, ctx.curblock, v.UserArgs[0], false, ast.TypeLiteral("uint64"), false)
		ln := ctx.curblock.NewExtractValue(v, []int64{1})
		return ln
	}

	fnc := ctx.Funcs[v.Name]
	if fnc.Function == nil {
		panic("Unknown function" + v.Name)
	}
	argdefs := fnc.FuncDecl.Args

	var args []value.Value
	for i, arg := range v.UserArgs {
		switch v.Name {
		case "PrintInt":
			// Hack until there's a stdlib. PrintInt needs to handle
			// every type of int, so we always convert it to int64
			val := evalValue(ctx, ctx.curblock, arg, false, nil, false)
			args = append(args, ctx.convValue(val, arg.Type(), ast.TypeLiteral("int64")))
		case "PrintString", "PrintByteSlice":
			// Hack because PrintString takes a i64 and an i8* as parameters,
			// but variable strings are typed as { i64, i64 }
			val := evalValue(ctx, ctx.curblock, arg, false, ast.TypeLiteral("uint8"), false)
			args = append(args, val)
		case "Create", "Open":
			val := evalValue(ctx, ctx.curblock, arg, false, ast.TypeLiteral("string"), false)
			args = append(args, val)
		case "Write", "Read":
			val := evalValue(ctx, ctx.curblock, arg, false, nil, false)
			args = append(args, val)
		case "Close", "Exit":
			val := evalValue(ctx, ctx.curblock, arg, false, nil, false)
			args = append(args, val)
		case "len":
			panic("len is a pseudo-function")
		default:
			switch t := argdefs[i].Type().(type) {
			case ast.ArrayType:
				base := evalValue(ctx, ctx.curblock, arg, true, nil, false)
				args = append(args, base)
			case ast.SliceType:
				base := evalValue(ctx, ctx.curblock, arg, false, nil, false)
				args = append(args, base)
			case ast.SumType:
				base := evalValue(ctx, ctx.curblock, arg, false, t, false)
				args = append(args, base)
			default:
				if argdefs[i].Reference {
					args = append(args, evalValue(ctx, ctx.curblock, arg, true, nil, false))
				} else {
					args = append(args, evalValue(ctx, ctx.curblock, arg, false, nil, false))
				}
			}
		}
	}
	f := ctx.curblock.NewCall(fnc.Function, args...)
	return f
}

func inferType(val ast.Value) ast.Type {
	switch v := val.(type) {
	case ast.LessThanOrEqualComparison:
		if ast.IsLiteral(v.Left) && ast.IsLiteral(v.Right) {
			return ast.TypeLiteral("int")
		} else if ast.IsLiteral(v.Left) {
			return v.Right.Type()
		} else {
			return v.Left.Type()
		}
	case ast.LessThanComparison:
		if ast.IsLiteral(v.Left) && ast.IsLiteral(v.Right) {
			return ast.TypeLiteral("int")
		} else if ast.IsLiteral(v.Left) {
			return v.Right.Type()
		} else {
			return v.Left.Type()
		}
	case ast.NotEqualsComparison:
		if ast.IsLiteral(v.Left) && ast.IsLiteral(v.Right) {
			return ast.TypeLiteral("int")
		} else if ast.IsLiteral(v.Left) {
			return v.Right.Type()
		} else {
			return v.Left.Type()
		}
	case ast.EqualityComparison:
		if ast.IsLiteral(v.Left) && ast.IsLiteral(v.Right) {
			return ast.TypeLiteral("int")
		} else if ast.IsLiteral(v.Left) {
			return v.Right.Type()
		} else {
			return v.Left.Type()
		}
	case ast.GreaterComparison:
		if ast.IsLiteral(v.Left) && ast.IsLiteral(v.Right) {
			return ast.TypeLiteral("int")
		} else if ast.IsLiteral(v.Left) {
			return v.Right.Type()
		} else {
			return v.Left.Type()
		}
	case ast.GreaterOrEqualComparison:
		if ast.IsLiteral(v.Left) && ast.IsLiteral(v.Right) {
			return ast.TypeLiteral("int")
		} else if ast.IsLiteral(v.Left) {
			return v.Right.Type()
		} else {
			return v.Left.Type()
		}
	case ast.BoolLiteral:
		return ast.TypeLiteral("bool")
	default:
		panic(fmt.Sprintf("Unhandled comparison type inference %v", reflect.TypeOf(v)))
	}
}
