package irgen

import (
	"fmt"
	"reflect"

	"github.com/driusan/lang/compiler/ir"
	"github.com/driusan/lang/parser/ast"
)

type variableLayout struct {
	values   map[ast.Variable]ir.Register
	tempVars int
}

// Reserves the next available register for varname
func (c *variableLayout) NextLocalRegister(varname ast.Variable) ir.Register {
	if varname == "" {
		c.tempVars++
		return ir.LocalValue(len(c.values) + c.tempVars - 1)
	}

	c.values[varname] = ir.LocalValue(len(c.values) + c.tempVars)
	return c.values[varname]
}

// Reserves a register for a function parameter. This must be done for every
// parameter, before any LocalRegister calls are made.
func (c *variableLayout) FuncParamRegister(varname ast.Variable, i int) ir.Register {
	c.tempVars--

	c.values[varname] = ir.FuncArg(i)
	return c.values[varname]
}

// Sets a variable to refer to an existing register, without generating a new
// one.
func (c *variableLayout) SetLocalRegister(varname ast.Variable, val ir.Register) {
	c.values[varname] = val
}

// Gets the register for an existing variable. Panics on invalid variables.
func (c variableLayout) Get(varname ast.Variable) ir.Register {
	if varname == "" {
		panic("Can not get empty varname")
	}
	return c.values[varname]
}

// Gets the register for an existing variable, and a bool denoting whether
// the variable exists or not.
func (c variableLayout) SafeGet(varname ast.Variable) (ir.Register, bool) {
	v, ok := c.values[varname]
	return v, ok
}

// Compile takes an AST and writes the assembly that it compiles to to
// w.
func GenerateIR(node ast.Node) (ir.Func, error) {
	context := &variableLayout{make(map[ast.Variable]ir.Register), 0}
	switch n := node.(type) {
	case ast.ProcDecl:
		for i, arg := range n.Args {
			context.FuncParamRegister(arg.Name, i)
		}
		body, err := compileBlock(n.Body, context)
		if err != nil {
			return ir.Func{}, err
		}
		return ir.Func{Name: n.Name, Body: body, NumArgs: uint(len(n.Args))}, nil
	case ast.FuncDecl:
		for i, arg := range n.Args {
			context.FuncParamRegister(arg.Name, i)
		}
		body, err := compileBlock(n.Body, context)
		if err != nil {
			return ir.Func{}, err
		}
		return ir.Func{Name: n.Name, Body: body, NumArgs: uint(len(n.Args))}, nil
	default:
		panic(fmt.Sprintf("Unhandled Node type in compiler %v", reflect.TypeOf(n)))
	}
}

// calculate the IR to perform a function call.
func callFunc(fc ast.FuncCall, context *variableLayout) ([]ir.Opcode, error) {
	var ops []ir.Opcode
	// First handle any non-literals that need to be evaluated.
	for i, arg := range fc.Args {
		switch a := arg.(type) {
		case ast.FuncCall:
			fc, err := callFunc(a, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, fc...)
			ops = append(ops,
				ir.MOV{
					Src: ir.FuncRetVal(0),
					Dst: ir.FuncCallArg(i),
				},
			)
		case ast.AdditionOperator, ast.SubtractionOperator:
			arg, r, err := evaluateValue(a, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, arg...)
			ops = append(ops,
				ir.MOV{
					Src: r,
					Dst: ir.FuncCallArg(i),
				},
			)
		default:
			// Nothing, will handle below.
		}
	}
	for i, arg := range fc.Args {
		switch a := arg.(type) {
		case ast.StringLiteral, ast.IntLiteral, ast.BoolLiteral:
			ops = append(ops, ir.MOV{
				Src: getRegister(a, context),
				Dst: ir.FuncCallArg(i),
			})
		case ast.Variable:
			reg := context.Get(a)
			ops = append(ops, ir.MOV{
				Src: reg,
				Dst: ir.FuncCallArg(i),
			})

		case ast.FuncCall, ast.AdditionOperator, ast.SubtractionOperator:
			// handled above, do nothing.
		default:
			panic(fmt.Sprintf("Unhandled argument type in FuncCall %v", reflect.TypeOf(a)))
		}
	}

	// Perform the call.
	if fc.Name == "print" {
		ops = append(ops, ir.CALL{"printf"})
	} else {
		ops = append(ops, ir.CALL{ir.Fname(fc.Name)})
	}
	return ops, nil
}

var loopNum uint

func getRegister(n ast.Node, context *variableLayout) ir.Register {
	switch v := n.(type) {
	case ast.StringLiteral:
		return ir.StringLiteral(v)
	case ast.IntLiteral:
		return ir.IntLiteral(v)
	case ast.BoolLiteral:
		if v {
			return ir.IntLiteral(1)
		}
		return ir.IntLiteral(0)
	case ast.Variable:
		return context.Get(v)
	default:
		panic(fmt.Sprintf("Unhandled type in getRegister: %v", reflect.TypeOf(v)))
	}
}

func compileBlock(block ast.BlockStmt, context *variableLayout) ([]ir.Opcode, error) {
	var ops []ir.Opcode
	for _, stmt := range block.Stmts {
		switch s := stmt.(type) {
		case ast.FuncCall:
			fc, err := callFunc(s, context)
			if err != nil {
				return nil, err
			}
			ops = append(ops, fc...)
		case ast.LetStmt:
			switch v := s.Value.(type) {
			case ast.IntLiteral, ast.StringLiteral, ast.BoolLiteral:
				reg := context.NextLocalRegister(s.Var.Name)
				ops = append(ops, ir.MOV{
					Src: getRegister(v, context),
					Dst: reg,
				})
			case ast.AdditionOperator, ast.SubtractionOperator,
				ast.DivOperator, ast.MulOperator, ast.ModOperator:
				reg := context.NextLocalRegister(s.Var.Name)
				body, r, err := evaluateValue(s.Value, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, ir.MOV{
					Src: r,
					Dst: reg,
				})
			default:
				panic("Unsupported let statement assignment")

			}
		case ast.ReturnStmt:
			switch arg := s.Val.(type) {
			case ast.FuncCall:
				fc, err := callFunc(arg, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, fc...)

				// Calling the function already will have left
				// the value in FuncRetValRegister[0]
				ops = append(ops, ir.RET{})
			default:
				ops = append(ops, ir.MOV{
					Src: getRegister(arg, context),
					Dst: ir.FuncRetVal(0),
				})
				ops = append(ops, ir.RET{})
			}
		case ast.MutStmt:
			switch v := s.InitialValue.(type) {
			case ast.IntLiteral, ast.BoolLiteral, ast.StringLiteral:
				reg := context.NextLocalRegister(s.Var.Name)
				ops = append(ops, ir.MOV{
					Src: getRegister(s.InitialValue, context),
					Dst: reg,
				})

			case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
				reg := context.NextLocalRegister(s.Var.Name)
				body, r, err := evaluateValue(s.InitialValue, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, ir.MOV{
					Src: r,
					Dst: reg,
				})

			case ast.Variable:
				reg := context.NextLocalRegister(s.Var.Name)
				val := context.Get(v)
				ops = append(ops, ir.MOV{
					Src: val,
					Dst: reg,
				})
			default:
				panic(fmt.Sprintf("Unhandled type for MutStmt assignment %v", reflect.TypeOf(s.InitialValue)))
			}
		case ast.AssignmentOperator:
			switch s.Value.(type) {
			case ast.IntLiteral, ast.BoolLiteral, ast.StringLiteral:
				ops = append(ops, ir.MOV{
					Src: getRegister(s.Value, context),
					Dst: context.Get(s.Variable),
				})
			case ast.AdditionOperator, ast.SubtractionOperator, ast.DivOperator, ast.MulOperator, ast.ModOperator:
				body, r, err := evaluateValue(s.Value, context)
				if err != nil {
					return nil, err
				}
				ops = append(ops, body...)
				ops = append(ops, ir.MOV{
					Src: r,
					Dst: context.Get(s.Variable),
				})

			default:
				panic(fmt.Sprintf("Statement type assignment not implemented: %v", reflect.TypeOf(s.Value)))

			}
		case ast.WhileLoop:
			lname := ir.Label(fmt.Sprintf("loop%dend", loopNum))
			lcond := ir.Label(fmt.Sprintf("loop%dcond", loopNum))
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

			ops = append(ops, ir.JMP{lcond})
			ops = append(ops, lname)
		case *ast.IfStmt:
			iname := ir.Label(fmt.Sprintf("if%delse", loopNum))
			dname := ir.Label(fmt.Sprintf("if%delsedone", loopNum))
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
			ops = append(ops, ir.JMP{dname})
			ops = append(ops, iname)
			if len(s.Else.Stmts) != 0 {
				body, err := compileBlock(s.Else, context)
				if err != nil {
					panic(err)
				}
				ops = append(ops, body...)
			}
			ops = append(ops, dname)
		default:
			panic(fmt.Sprintf("Statement type not implemented: %v", reflect.TypeOf(s)))
		}
	}
	return ops, nil
}

// Evaluates a boolean condition. If the condition fails, jump to faillabel.
func evaluateCondition(val ast.BoolValue, context *variableLayout, faillabel ir.Label) ([]ir.Opcode, error) {
	var ops []ir.Opcode
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

		ops = append(ops, ir.JLE{
			ir.ConditionalJump{
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
		ops = append(ops, ir.JL{
			ir.ConditionalJump{Label: faillabel,
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
		ops = append(ops, ir.JNE{
			ir.ConditionalJump{Label: faillabel,

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
		ops = append(ops, ir.JE{
			ir.ConditionalJump{Label: faillabel,
				Src: r,
				Dst: r2,
			},
		})
		return ops, nil
	default:
		panic(fmt.Sprintf("Condition type not implemented: %v", reflect.TypeOf(c)))
	}
	return ops, nil
}

// Evaluates a value expression and returns the opcodes to evaluate it, and the
// register which contains the value evaluated.
func evaluateValue(val ast.Value, context *variableLayout) ([]ir.Opcode, ir.Register, error) {
	var ops []ir.Opcode
	switch s := val.(type) {
	case ast.AdditionOperator:
		a := context.NextLocalRegister("")
		switch s.Left.(type) {
		case ast.IntLiteral, ast.Variable:
			ops = append(ops, ir.ADD{
				Src: getRegister(s.Left, context),
				Dst: a,
			})
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
			body, r, err := evaluateValue(s.Left, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			a = r
		default:
			panic(fmt.Sprintf("Unhandled left parameter in addition %v", reflect.TypeOf(s.Left)))
		}

		var r ir.Register
		switch s.Right.(type) {
		case ast.IntLiteral, ast.Variable:
			r = getRegister(s.Right, context)
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
			body, r2, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			r = r2
		default:
			panic(fmt.Sprintf("Unhandled right parameter in addition: %v", reflect.TypeOf(s.Right)))

		}
		ops = append(ops, ir.ADD{
			Src: r,
			Dst: a,
		})
		return ops, a, nil
	case ast.SubtractionOperator:
		a := context.NextLocalRegister("")
		switch s.Left.(type) {
		case ast.IntLiteral, ast.Variable:
			ops = append(ops, ir.MOV{
				Src: getRegister(s.Left, context),
				Dst: a,
			})
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
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
		case ast.IntLiteral, ast.Variable:
			ops = append(ops, ir.SUB{
				Src: getRegister(s.Right, context),
				Dst: a,
			})
		case ast.AdditionOperator, ast.SubtractionOperator, ast.MulOperator, ast.DivOperator, ast.ModOperator:
			body, r, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			ops = append(ops, ir.SUB{
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

		a := context.NextLocalRegister("")
		ops = append(ops, bodyb...)
		ops = append(ops, ir.MOD{
			Left:  ra,
			Right: rb,
			Dst:   a,
		})
		return ops, a, nil
	case ast.MulOperator:
		a := context.NextLocalRegister("")
		var l, r ir.Register
		switch s.Left.(type) {
		case ast.IntLiteral, ast.Variable:
			l = getRegister(s.Left, context)
		default:
			panic(fmt.Sprintf("Unhandled left parameter in mul %v", reflect.TypeOf(s.Left)))
		}

		switch s.Right.(type) {
		case ast.IntLiteral, ast.Variable:
			r = getRegister(s.Right, context)
		case ast.SubtractionOperator, ast.AdditionOperator, ast.MulOperator, ast.DivOperator:
			body, reg, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			r = reg
		}
		ops = append(ops, ir.MUL{
			Left:  l,
			Right: r,
			Dst:   a,
		})
		return ops, a, nil
	case ast.DivOperator:
		a := context.NextLocalRegister("")
		var l, r ir.Register
		switch s.Left.(type) {
		case ast.IntLiteral, ast.Variable:
			l = getRegister(s.Left, context)
		default:
			panic(fmt.Sprintf("Unhandled left parameter in div %v", reflect.TypeOf(s.Left)))
		}

		switch s.Right.(type) {
		case ast.IntLiteral, ast.Variable:
			r = getRegister(s.Right, context)
		case ast.SubtractionOperator, ast.AdditionOperator, ast.MulOperator, ast.DivOperator:
			body, reg, err := evaluateValue(s.Right, context)
			if err != nil {
				return nil, nil, err
			}
			ops = append(ops, body...)
			r = reg
		default:
			panic(fmt.Sprintf("Unhandled right parameter in div: %v", reflect.TypeOf(s.Right)))

		}
		ops = append(ops, ir.DIV{
			Left:  l,
			Right: r,
			Dst:   a,
		})
		return ops, a, nil
	case ast.Variable, ast.IntLiteral, ast.BoolLiteral:
		return nil, getRegister(s, context), nil
	default:
		panic(fmt.Errorf("Unhandled value type: %v", reflect.TypeOf(s)))
	}
}
