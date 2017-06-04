package codegen

import (
	"fmt"
	"reflect"

	"github.com/driusan/lang/compiler/ir"
	"github.com/driusan/lang/parser/ast"
)

type CPU interface {
	ToPhysical(ir.Register)
	ConvertInstruction(ir.Opcode) string
}

type amd64Registers struct {
	ax, bx, cx, dx, si, di, bp, r8, r9, r10, r11, r12, r13, r14, r15 ir.Register
}
type Amd64 struct {
	amd64Registers

	// List of what ir pseudo-register each physical register is currently
	// mapped to.
	stringLiterals map[ir.StringLiteral]PhysicalRegister

	// The number of arguments in the currently used function. Used for
	// ToPhysical to calculate where local (non-parameter) variables start.
	numArgs uint
}

func (a *amd64Registers) nextPhysicalRegister(r ir.Register, skipDX bool) (PhysicalRegister, error) {
	// Avoids AX and BP, since AX is the return register and BP is the first
	// argument to a function call.
	if a.bx == nil {
		a.bx = r
		return "BX", nil
	}
	if a.cx == nil {
		a.cx = r
		return "CX", nil
	}
	if a.dx == nil && !skipDX {
		a.dx = r
		return "DX", nil
	}
	if a.si == nil {
		a.si = r
		return "SI", nil
	}
	if a.di == nil {
		a.di = r
		return "DI", nil
	}
	if a.r8 == nil {
		a.r8 = r
		return "R8", nil
	}
	if a.r9 == nil {
		a.r9 = r
		return "R9", nil
	}
	if a.r10 == nil {
		a.r10 = r
		return "R10", nil
	}
	if a.r11 == nil {
		a.r11 = r
		return "R11", nil
	}
	if a.r12 == nil {
		a.r12 = r
		return "R12", nil
	}
	if a.r13 == nil {
		a.r13 = r
		return "R13", nil
	}
	if a.r14 == nil {
		a.r14 = r
		return "R14", nil
	}
	if a.r15 == nil {
		a.r15 = r
		return "R15", nil
	}
	return "", fmt.Errorf("No physical registers available")
}
func (a *amd64Registers) tempPhysicalRegister(skipDX bool) (PhysicalRegister, error) {
	// Avoids AX and BP, since AX is the return register and BP is the first
	// argument to a function call.
	if a.bx == nil {
		return "BX", nil
	}
	if a.cx == nil {
		return "CX", nil
	}
	if a.dx == nil && !skipDX {
		return "DX", nil
	}
	if a.si == nil {
		return "SI", nil
	}
	if a.di == nil {
		return "DI", nil
	}
	if a.r8 == nil {
		return "R8", nil
	}
	if a.r9 == nil {
		return "R9", nil
	}
	if a.r10 == nil {
		return "R10", nil
	}
	if a.r11 == nil {
		return "R11", nil
	}
	if a.r12 == nil {
		return "R12", nil
	}
	if a.r13 == nil {
		return "R13", nil
	}
	if a.r14 == nil {
		return "R14", nil
	}
	if a.r15 == nil {
		return "R15", nil
	}
	return "", fmt.Errorf("No physical registers available")
}

func (a amd64Registers) getPhysicalRegister(r ir.Register) (PhysicalRegister, error) {
	switch {
	case a.ax == r:
		return "AX", nil
	case a.bx == r:
		return "BX", nil
	case a.cx == r:
		return "CX", nil
	case a.dx == r:
		return "DX", nil
	case a.si == r:
		return "SI", nil
	case a.di == r:
		return "DI", nil
	case a.r8 == r:
		return "R8", nil
	case a.r9 == r:
		return "R9", nil
	case a.r10 == r:
		return "R10", nil
	case a.r11 == r:
		return "R11", nil
	case a.r12 == r:
		return "R12", nil
	case a.r13 == r:
		return "R13", nil
	case a.r14 == r:
		return "R14", nil
	case a.r15 == r:
		return "R15", nil
	}
	return "", fmt.Errorf("Register not mapped")
}
func (a *amd64Registers) clearRegisterMapping(f ir.Func) {
	a.ax = nil
	a.bx = nil
	a.cx = nil
	a.dx = nil
	a.si = nil
	a.di = nil
	if f.NumArgs != 0 {
		// FIXME: This shouldn't make assumptions about the type in
		// arg 0.
		a.bp = ir.FuncArg{0, ast.TypeInfo{8, true}}
	} else {
		a.bp = nil
	}
	a.r8 = nil
	a.r9 = nil
	a.r10 = nil
	a.r11 = nil
	a.r12 = nil
	a.r13 = nil
	a.r14 = nil
	a.r15 = nil
}

func (a *Amd64) ToPhysical(r ir.Register) PhysicalRegister {
	switch v := r.(type) {
	case ir.StringLiteral:
		return PhysicalRegister("$" + string(a.stringLiterals[v]) + "+0(SB)")
	case ir.IntLiteral:
		return PhysicalRegister(fmt.Sprintf("$%d", v))
	case ir.FuncCallArg:
		if v.Id == 0 {
			return "BP"
		}
		return PhysicalRegister(fmt.Sprintf("%d(SP)", 8*v.Id))
	case ir.LocalValue:
		return PhysicalRegister(fmt.Sprintf("%v+%d(FP)", v.String(), (int(v.Id-a.numArgs)*8)+(int(a.numArgs)*8)))
	case ir.FuncRetVal:
		if v.Id == 0 {
			return "AX"
		}
		panic("Multiple return values not yet implemented.")
	case ir.FuncArg:
		// First check if the arg is already in a register.
		r, err := a.getPhysicalRegister(v)
		if err == nil {
			return r
		}

		// Otherwise, the first arg goes in BP, and the rest are on
		// the stack.
		if v.Id == 0 {
			return "BP"
		}
		// FIXME: The prefix of this is supposed to be the variable name,
		// not the IR register name..
		return PhysicalRegister(fmt.Sprintf("%v+%d(FP)", v.String(), int(v.Id)*8))
	default:
		panic(fmt.Sprintf("Unhandled register type %v", reflect.TypeOf(v)))
	}
}

func isFA0(v ir.Register) bool {
	if r, ok := v.(ir.FuncArg); ok && r.Id == 0 {
		return true
	}
	return false
}
func checkBPUsed(i int, ops []ir.Opcode) bool {
	// we're likely about to wipe out ir.FuncArg(0).
	// Check if we care.
	for j := i + 1; j < len(ops); j++ {
		for _, v := range ops[j].Registers() {
			if isFA0(v) {
				return true
			}
		}
	}
	return false
}

func (a Amd64) opSuffix(sizeInBytes int) string {
	switch sizeInBytes {
	case 1:
		return "BLSX"
	case 2:
		return "WLSX"
	case 4:
		return "L"
	case 0, 8:
		return "Q"
	default:
		panic("Unhandled register size in MOV")
	}
}
func (a *Amd64) ConvertInstruction(i int, ops []ir.Opcode) string {
	op := ops[i]
	switch o := op.(type) {
	case ir.Label:
		return o.String()
	case ir.MOV:
		dst := a.ToPhysical(o.Dst)
		v := ""
		if dst == "BP" && !isFA0(o.Dst) && checkBPUsed(i, ops) {
			// Move it to the next free register.
			fpreserve, err := a.nextPhysicalRegister(a.bp, false)
			if err != nil {
				panic(err)
			}
			v += fmt.Sprintf("\tMOV%v BP, %v\n\t", a.opSuffix(a.bp.Size()), fpreserve)
		}
		var src PhysicalRegister
		switch val := o.Src.(type) {
		case ir.LocalValue, ir.FuncArg:
			// First check if the arg is already in a register.
			r, err := a.getPhysicalRegister(val)
			if err == nil {
				src = r
				break
			}
			src, err = a.tempPhysicalRegister(false)
			if err != nil {
				panic(err)
			}
			suffix := a.opSuffix(val.Size())
			v += fmt.Sprintf("\tMOV%v %v, %v\n\t", suffix, a.ToPhysical(val), src)

		default:
			src = a.ToPhysical(val)
		}

		v += fmt.Sprintf("MOV%v %v, %v", a.opSuffix(o.Src.Size()), src, dst)
		return v
	case ir.CALL:
		var v string
		if len(o.Args) > 0 && a.bp != nil && checkBPUsed(i-1, ops) {
			// Move it to the next free register.
			fpreserve, err := a.nextPhysicalRegister(a.bp, false)
			if err != nil {
				panic(err)
			}
			v += fmt.Sprintf("\tMOVQ BP, %v\n\t", fpreserve)
		}
		for i, arg := range o.Args {
			var fa PhysicalRegister
			if o.TailCall {
				// If it's a tail call, the dst should get optimized
				// to the same location as this call's.
				fa = a.ToPhysical(ir.FuncArg{uint(i), ast.TypeInfo{arg.Size(), arg.Signed()}})
			} else {
				fa = a.ToPhysical(ir.FuncCallArg{i, ast.TypeInfo{arg.Size(), arg.Signed()}})
			}
			var physArg PhysicalRegister
			switch arg.(type) {
			case ir.LocalValue, ir.FuncArg:
				// First check if the arg is already in a register.
				r, err := a.getPhysicalRegister(arg)
				if err == nil {
					physArg = r
					break
				}
				physArg, err = a.tempPhysicalRegister(false)
				if err != nil {
					panic(err)
				}
				suffix := a.opSuffix(arg.Size())
				v += fmt.Sprintf("\tMOV%v %v, %v\n\t", suffix, a.ToPhysical(arg), physArg)
			default:
				physArg = a.ToPhysical(arg)
			}

			v += fmt.Sprintf("MOVQ %v, %v\n\t", physArg, fa)
		}
		if o.TailCall {
			// Optimize the call away to a JMP and reuse the stack
			// frame.
			tmp, err := a.tempPhysicalRegister(false)
			if err != nil {
				panic(err)
			}
			// Jump 1 instruction past the start of this symbol.
			// The first instruction is SUBQ $k, SP which the linker
			// inserts. We want to reuse the stack, not push to it.
			//
			// FIXME: This needs to manually adjust SP if the stack
			// space reserved for the new symbol isn't the same as
			// the current function.
			v += fmt.Sprintf("MOVQ $%v+4(SB), %v\n\t", o.FName, tmp)
			return v + fmt.Sprintf("JMP %v", tmp)
		} else {
			return v + fmt.Sprintf("CALL %v+0(SB)", o.FName)
		}
	case ir.RET:
		return fmt.Sprintf("RET")
	case ir.ADD:
		dst, err := a.getPhysicalRegister(o.Dst)
		v := ""
		if err != nil {
			dst, err = a.nextPhysicalRegister(o.Dst, false)
			if err != nil {
				panic(err)
			}
			// This is the first time using this register for this
			// value, so just MOV the value into it and trash whatever
			// was there before.
			return fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), dst)
		}
		var src PhysicalRegister
		switch val := o.Src.(type) {
		case ir.LocalValue:
			// First check if the arg is already in a register.
			r, err := a.getPhysicalRegister(val)
			if err == nil {
				src = r
				break
			}
			src, err = a.tempPhysicalRegister(false)
			if err != nil {
				panic(err)
			}
			v += fmt.Sprintf("\tMOVQ %v, %v\n\t", a.ToPhysical(val), src)
		default:
			src = a.ToPhysical(val)
		}

		v += fmt.Sprintf("ADDQ %v, %v", src, dst)
		return v
	case ir.SUB:
		// Special cases: 1, 0, and -1
		if o.Src == ir.IntLiteral(0) {
			// Subtracting 0 from something is stupid.
			return ""
		} else if o.Src == ir.IntLiteral(1) {
			return fmt.Sprintf("DECQ %v", a.ToPhysical(o.Dst))
		} else if o.Src == ir.IntLiteral(-1) {
			return fmt.Sprintf("INCQ %v", a.ToPhysical(o.Dst))
		}
		// Normal subtraction.
		// FIXME: This is only required if o.Right isn't really a register,
		// but a fake register like "$15". This also should have a better
		// way to keep track of what the register is and delete it.
		r, err := a.nextPhysicalRegister(ir.FuncArg{9999, ast.TypeInfo{8, true}}, true)
		if err != nil {
			panic(err)
		}

		v := fmt.Sprintf("MOVQ %v, %v\n\t", a.ToPhysical(o.Src), r)

		return v + fmt.Sprintf("SUBQ %v, %v", r, a.ToPhysical(o.Dst))
	case ir.MOD:
		v := ""
		// DIV clobbers DX with the result of the MOD, so if there's
		// anything there preserve it
		popax, popdx := false, false
		if a.ax != nil && a.ax != o.Left {
			v += "PUSHQ AX\n\t"
			popax = true
		}

		if a.dx != nil && a.dx != o.Dst {
			v += "PUSHQ DX\n\t"
			popdx = true
		}
		v += "MOVQ $0, DX\n\t"
		v += fmt.Sprintf("MOVQ %v, AX // %v\n\t", a.ToPhysical(o.Left), o.Left)

		// FIXME: This is only required if o.Right isn't really a register,
		// but a fake register like "$15". This also should have a better
		// way to keep track of what the register is and delete it.
		r, err := a.nextPhysicalRegister(ir.FuncArg{9999, ast.TypeInfo{8, true}}, true)
		if err != nil {
			panic(err)
		}

		v += fmt.Sprintf("MOVQ %v, %v\n\t", a.ToPhysical(o.Right), r)
		v += fmt.Sprintf("DIVQ %v\n\t", r)
		v += fmt.Sprintf("MOVQ DX, %v", a.ToPhysical(o.Dst))
		if popdx {
			v += "\n\tPOPQ DX"
		}
		if popax {
			v += "\n\tPOPQ AX"
		}
		return v
	case ir.DIV:
		v := ""
		// DIV clobbers DX with the result of the MOD, so if there's
		// anything there preserve it
		popax, popdx := false, false
		if a.ax != nil && a.ax != o.Left {
			v += "PUSHQ AX\n\t"
			popax = true
		}

		if a.dx != nil && a.dx != o.Dst {
			v += "PUSHQ DX\n\t"
			popdx = true
		}
		v += "MOVQ $0, DX\n\t"
		v += fmt.Sprintf("MOVQ %v, AX // %v\n\t", a.ToPhysical(o.Left), o.Left)

		// FIXME: This is only required if o.Right isn't really a register,
		// but a fake register like "$15". This also should have a better
		// way to keep track of what the register is and delete it.
		r, err := a.nextPhysicalRegister(ir.FuncArg{9999, ast.TypeInfo{8, true}}, true)
		if err != nil {
			panic(err)
		}

		v += fmt.Sprintf("MOVQ %v, %v\n\t", a.ToPhysical(o.Right), r)
		v += fmt.Sprintf("DIVQ %v\n\t", r)
		v += fmt.Sprintf("MOVQ AX, %v", a.ToPhysical(o.Dst))
		if popdx {
			v += "\n\tPOPQ DX"
		}
		if popax {
			v += "\n\tPOPQ AX"
		}
		return v
	case ir.MUL:
		v := ""
		// MUL multiples AX by the operand, and puts the overflow in
		// DX, so preserve them if there's anything there.
		popax, popdx := false, false
		if a.ax != nil && a.ax != o.Left {
			v += "PUSHQ AX\n\t"
			popax = true
		}

		if a.dx != nil && a.dx != o.Dst {
			v += "PUSHQ DX\n\t"
			popdx = true
		}
		//		v += "MOVQ $0, DX\n\t"
		v += fmt.Sprintf("MOVQ %v, AX // %v\n\t", a.ToPhysical(o.Left), o.Left)

		// FIXME: This is only required if o.Right isn't really a register,
		// but a fake register like "$15". This also should have a better
		// way to keep track of what the register is and delete it.
		r, err := a.nextPhysicalRegister(ir.FuncArg{9999, ast.TypeInfo{8, true}}, true)
		if err != nil {
			panic(err)
		}

		v += fmt.Sprintf("MOVQ %v, %v\n\t", a.ToPhysical(o.Right), r)
		v += fmt.Sprintf("MULQ %v\n\t", r)
		v += fmt.Sprintf("MOVQ AX, %v", a.ToPhysical(o.Dst))
		if popdx {
			v += "\n\tPOPQ DX"
		}
		if popax {
			v += "\n\tPOPQ AX"
		}
		return v
	case ir.JMP:
		return fmt.Sprintf("JMP %v", o.Label.Inline())
	case ir.JE:
		// FIXME: Only required if both src and dst are not really registers
		src, err := a.tempPhysicalRegister(false)
		if err != nil {
			panic(err)
		}
		v := fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), src)
		return v + fmt.Sprintf("\n\tCMPQ %v, %v\n\tJE %v", src, a.ToPhysical(o.Dst), o.Label.Inline())
	case ir.JL:
		// FIXME: Only required if both src and dst are not really registers
		src, err := a.tempPhysicalRegister(false)
		if err != nil {
			panic(err)
		}
		v := fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), src)
		return v + fmt.Sprintf("\n\tCMPQ %v, %v\n\tJL %v", src, a.ToPhysical(o.Dst), o.Label.Inline())
	case ir.JLE:
		// FIXME: Only required if both src and dst are not really registers
		src, err := a.tempPhysicalRegister(false)
		if err != nil {
			panic(err)
		}
		v := fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), src)
		return v + fmt.Sprintf("\n\tCMPQ %v, %v\n\tJLE %v", src, a.ToPhysical(o.Dst), o.Label.Inline())
	case ir.JNE:
		// FIXME: Only required if both src and dst are not really registers
		src, err := a.tempPhysicalRegister(false)
		if err != nil {
			panic(err)
		}
		v := fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), src)
		return v + fmt.Sprintf("\n\tCMPQ %v, %v\n\tJNE %v", src, a.ToPhysical(o.Dst), o.Label.Inline())
	case ir.JGE:
		// FIXME: Only required if both src and dst are not really registers
		src, err := a.tempPhysicalRegister(false)
		if err != nil {
			panic(err)
		}
		v := fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), src)
		return v + fmt.Sprintf("\n\tCMPQ %v, %v\n\tJGE %v", src, a.ToPhysical(o.Dst), o.Label.Inline())
	case ir.JG:
		// FIXME: Only required if both src and dst are not really registers
		src, err := a.tempPhysicalRegister(false)
		if err != nil {
			panic(err)
		}
		v := fmt.Sprintf("MOVQ %v, %v", a.ToPhysical(o.Src), src)
		return v + fmt.Sprintf("\n\tCMPQ %v, %v\n\tJG %v", src, a.ToPhysical(o.Dst), o.Label.Inline())
	default:
		panic(fmt.Sprintf("Unhandled instruction in AMD64 code generation %v", reflect.TypeOf(o)))
	}
}
