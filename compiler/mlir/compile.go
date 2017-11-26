package mlir

import (
	"fmt"
	"reflect"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/ast"
)

var branchNum uint = 0

type jumpType byte

const (
	notComparison = jumpType(iota)
	jumpSuccess
	jumpFailure
)

// Compile takes an AST and writes the assembly that it compiles to to
// w.
func Generate(node ast.Node, typeInfo ast.TypeInformation, callables ast.Callables, enums hlir.EnumMap) (Func, hlir.EnumMap, error) {
	hlirfunc, newenums, regData, err := hlir.Generate(node, typeInfo, callables, enums)
	if err != nil {
		return Func{}, nil, err
	}
	if enums == nil {
		enums = make(hlir.EnumMap)
	}
	for k, v := range newenums {
		enums[k] = v
	}
	branchNum = 0
	ctx := NewContext(callables, regData)

	f := Func{Name: hlirfunc.Name, NumArgs: hlirfunc.NumArgs, NumLocals: hlirfunc.NumLocals}

	ctx.curFunc = &f
	for _, op := range hlirfunc.Body {
		f.Body = append(f.Body, ctx.convertOp(op, "", notComparison)...)
	}
	return f, enums, nil
}

func (ctx *Context) convertOp(op hlir.Opcode, conditionLabel Label, jt jumpType) []Opcode {
	switch o := op.(type) {
	case hlir.CALL:
		c := CALL{
			FName:    Fname(o.FName),
			Args:     make([]Register, 0, len(o.Args)),
			TailCall: o.TailCall,
		}
		for _, r := range o.Args {
			c.Args = append(c.Args, ctx.convertRegister(r))
		}
		fi := ctx.GetFuncInfo(string(o.FName))
		largestFuncCall := uint(len(c.Args))
		if rsize := uint(len(fi.ReturnTuple())); rsize > 0 {
			largestFuncCall += rsize + 1
		}

		if largestFuncCall > ctx.curFunc.LargestFuncCall {
			ctx.curFunc.LargestFuncCall = largestFuncCall
		}
		return []Opcode{c}
	case hlir.MOV:
		return []Opcode{
			MOV{
				Src: ctx.convertRegister(o.Src),
				Dst: ctx.convertRegister(o.Dst),
			},
		}

	case hlir.RET:
		return []Opcode{RET{}}
	case hlir.ADD:
		return []Opcode{
			MOV{
				Src: ctx.convertRegister(o.Left),
				Dst: ctx.convertRegister(o.Dst),
			},
			ADD{
				Src: ctx.convertRegister(o.Right),
				Dst: ctx.convertRegister(o.Dst),
			},
		}
	case hlir.SUB:
		return []Opcode{
			MOV{
				Src: ctx.convertRegister(o.Left),
				Dst: ctx.convertRegister(o.Dst),
			},
			SUB{
				Src: ctx.convertRegister(o.Right),
				Dst: ctx.convertRegister(o.Dst),
			},
		}
	case hlir.MUL:
		return []Opcode{
			MUL{
				ctx.convertRegister(o.Left),
				ctx.convertRegister(o.Right),
				ctx.convertRegister(o.Dst),
			},
		}
	case hlir.DIV:
		return []Opcode{
			DIV{
				ctx.convertRegister(o.Left),
				ctx.convertRegister(o.Right),
				ctx.convertRegister(o.Dst),
			},
		}
	case hlir.MOD:
		return []Opcode{
			MOD{
				ctx.convertRegister(o.Left),
				ctx.convertRegister(o.Right),
				ctx.convertRegister(o.Dst),
			},
		}

	case hlir.LOOP:
		var ops []Opcode
		cond := Label(fmt.Sprintf("loop%dcond", branchNum))
		end := Label(fmt.Sprintf("loop%dend", branchNum))
		branchNum++
		ops = append(ops, cond)
		for _, op := range o.Condition.Body {
			ops = append(ops, ctx.convertOp(op, end, jumpFailure)...)
		}

		for _, op := range o.Body {
			ops = append(ops, ctx.convertOp(op, end, notComparison)...)
		}
		ops = append(ops, JMP{cond})
		ops = append(ops, end)
		return ops
	case hlir.IF:
		var ops []Opcode
		elselabel := Label(fmt.Sprintf("if%delse", branchNum))
		end := Label(fmt.Sprintf("if%delsedone", branchNum))
		branchNum++
		for _, op := range o.ControlFlow.Condition.Body {
			ops = append(ops, ctx.convertOp(op, elselabel, jumpFailure)...)
		}
		for _, op := range o.Body {
			ops = append(ops, ctx.convertOp(op, elselabel, notComparison)...)
		}
		ops = append(ops, JMP{end})
		ops = append(ops, elselabel)
		for _, op := range o.ElseBody {
			ops = append(ops, ctx.convertOp(op, end, notComparison)...)
		}
		ops = append(ops, end)
		return ops
	case hlir.JumpTable:
		var ops []Opcode

		// First generate the labels for each case and an end label, in case on of
		// them has an embedded branch and changes branchNum
		var labels []Label = make([]Label, 0, len(o)+1)
		for i := range o {
			labels = append(labels, Label(fmt.Sprintf("match%dv%d", branchNum, i)))
		}
		end := Label(fmt.Sprintf("match%ddone", branchNum))

		for i, c := range o {
			caselabel := labels[i]
			for _, op := range c.Condition.Body {
				ops = append(ops, ctx.convertOp(op, caselabel, jumpSuccess)...)
			}
		}
		ops = append(ops, JMP{end})
		for i, c := range o {
			ops = append(ops, labels[i])
			for _, op := range c.Body {
				ops = append(ops, ctx.convertOp(op, "", notComparison)...)
			}
			ops = append(ops, JMP{end})
		}
		ops = append(ops, end)
		branchNum++
		return ops
	case hlir.EQ:
		switch jt {
		case jumpSuccess:
			return []Opcode{JE{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		case jumpFailure:
			return []Opcode{JNE{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		default:
			panic("Equals used outside of a comparison context")
		}
	case hlir.NEQ:
		return []Opcode{JE{
			ConditionalJump{
				conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
		},
		}
	case hlir.GT:
		switch jt {
		case jumpSuccess:
			return []Opcode{JG{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		case jumpFailure:
			return []Opcode{JLE{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		default:
			panic("Equals used outside of a comparison context")
		}
	case hlir.GEQ:
		switch jt {
		case jumpSuccess:
			return []Opcode{JGE{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		case jumpFailure:
			return []Opcode{JL{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		default:
			panic("Equals used outside of a comparison context")
		}
	case hlir.LTE:
		switch jt {
		case jumpSuccess:
			return []Opcode{JLE{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		case jumpFailure:
			return []Opcode{JG{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		default:
			panic("Equals used outside of a comparison context")
		}
	case hlir.LT:
		switch jt {
		case jumpSuccess:
			return []Opcode{JL{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		case jumpFailure:
			return []Opcode{JGE{
				ConditionalJump{
					conditionLabel, ctx.convertRegister(o.Left), ctx.convertRegister(o.Right)},
			},
			}
		default:
			panic("Equals used outside of a comparison context")
		}
	default:
		// Uncomment this after the tests all pass.
		panic(fmt.Sprintf("Unhandled op type %v", reflect.TypeOf(op)))
	}
}

func (ctx Context) convertRegister(reg hlir.Register) Register {
	switch r := reg.(type) {
	case hlir.IntLiteral:
		return IntLiteral(r)
	case hlir.StringLiteral:
		return StringLiteral(r)
	case hlir.LocalValue:
		ti := ctx.GetTypeInfo(reg)
		return LocalValue{uint(r), ti}
	case hlir.FuncRetVal:
		ti := ctx.GetTypeInfo(reg)
		return FuncRetVal{uint(r), ti}
	case hlir.TempValue:
		return TempValue(r)
	case hlir.FuncArg:
		return FuncArg{uint(r.Id), ctx.GetTypeInfo(reg), r.Reference}
	case hlir.Offset:
		scale := uint(r.Scale)
		if scale == 0 {
			scale = 8
		}
		return Offset{
			Offset: ctx.convertRegister(r.Offset),
			Scale:  scale,
			Base:   ctx.convertRegister(r.Base),
		}
	case hlir.Pointer:
		return Pointer{
			ctx.convertRegister(r.Register),
		}
	case hlir.LastFuncCallRetVal:
		return FuncRetVal{r.RetNum, ctx.GetTypeInfo(reg)}
	default:
		panic(fmt.Sprintf("Unhandled register type %v", reflect.TypeOf(r)))
	}
}
