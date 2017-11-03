package hlir

import (
	"fmt"
)

type Opcode interface {
	Registers() []Register
	ModifiedRegisters() []Register
	String() string
}

type RET struct{}

func (r RET) String() string {
	return "RET\n"
}

func (o RET) Registers() []Register {
	return nil
}

func (o RET) ModifiedRegisters() []Register {
	return nil
}

type FName string
type CALL struct {
	FName    FName
	Args     []Register
	TailCall bool
}

func (c CALL) String() string {
	return fmt.Sprintf("CALL %v (%v)", c.FName, c.Args)
}

func (o CALL) Registers() []Register {
	return o.Args
}

func (o CALL) ModifiedRegisters() []Register {
	return nil
}

type MOV struct {
	Src, Dst Register
}

func (m MOV) String() string {
	return fmt.Sprintf("MOV %v, %v\n", m.Src, m.Dst)
}

func (o MOV) Registers() []Register {
	return []Register{o.Src, o.Dst}
}

func (o MOV) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type ADD struct {
	Left, Right, Dst Register
}

func (o ADD) String() string {
	return fmt.Sprintf("ADD %v + %v => %v\n", o.Left, o.Right, o.Dst)
}

func (o ADD) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o ADD) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type SUB struct {
	Left, Right, Dst Register
}

func (o SUB) String() string {
	return fmt.Sprintf("SUB %s - %s => %v\n", o.Left, o.Right, o.Dst)
}

func (o SUB) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o SUB) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type DIV struct {
	Left, Right, Dst Register
}

func (o DIV) String() string {
	return fmt.Sprintf("DIV %s / %s => %s", o.Left, o.Right, o.Dst)
}

func (o DIV) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o DIV) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type MUL struct {
	Left, Right, Dst Register
}

func (o MUL) String() string {
	return fmt.Sprintf("MUL %s * %s => %s", o.Left, o.Right, o.Dst)
}

func (o MUL) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o MUL) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type MOD struct {
	Left, Right, Dst Register
}

func (o MOD) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o MOD) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

func (o MOD) String() string {
	return fmt.Sprintf("MOD %s, %s, %s\n", o.Left, o.Right, o.Dst)
}

type ControlFlow struct {
	Condition
	Body []Opcode
}

func (o ControlFlow) Registers() []Register {
	r := o.Condition.Registers()
	for _, op := range o.Body {
		r = append(r, op.Registers()...)
	}
	return r

}
func (o ControlFlow) ModifiedRegisters() []Register {
	r := o.Condition.ModifiedRegisters()
	for _, op := range o.Body {
		r = append(r, op.ModifiedRegisters()...)
	}
	return r

}

type IF struct {
	ControlFlow
	ElseBody []Opcode
}

func (o IF) String() string {
	return fmt.Sprintf("IF %s (%s) ELSE (%s)\n", o.Condition, o.Body, o.ElseBody)
}

func (o IF) Registers() []Register {
	r := o.ControlFlow.Registers()
	for _, op := range o.ElseBody {
		r = append(r, op.Registers()...)
	}
	return r
}

func (o IF) ModifiedRegisters() []Register {
	r := o.ControlFlow.ModifiedRegisters()
	for _, op := range o.ElseBody {
		r = append(r, op.ModifiedRegisters()...)
	}
	return r
}

type Condition struct {
	// The opcodes to run to evaluate the condition. This happens at the start of
	Body []Opcode
	// The Register to compare after evaluating the condition body. The operation
	// is implicitly comparing the Register to 0.
	Register Register
}

func (c Condition) String() string {
	return fmt.Sprintf("(%v , %v != 0)", c.Body, c.Register)
}

func (c Condition) Registers() []Register {
	r := []Register{c.Register}
	for _, op := range c.Body {
		r = append(r, op.Registers()...)
	}
	return r

}
func (c Condition) ModifiedRegisters() []Register {
	var r []Register
	for _, op := range c.Body {
		r = append(r, op.ModifiedRegisters()...)
	}
	return r

}

type LOOP ControlFlow

func (o LOOP) String() string {
	return fmt.Sprintf("LOOP %v (%s)\n", o.Condition, o.Body)
}

func (o LOOP) Registers() []Register {
	return ControlFlow(o).Registers()
}

func (o LOOP) ModifiedRegisters() []Register {
	return ControlFlow(o).ModifiedRegisters()
}

type JumpTable []ControlFlow

func (jt JumpTable) String() string {
	v := ""
	for _, s := range jt {
		v += s.String() + fmt.Sprintf(" %v\n", s.Body)

	}
	return fmt.Sprintf("JumpTable (\n%v)", v)
}

func (jt JumpTable) Registers() []Register {
	var r []Register
	for _, s := range jt {
		r = append(r, s.Registers()...)

	}
	return r
}

func (jt JumpTable) ModifiedRegisters() []Register {
	var r []Register
	for _, s := range jt {
		r = append(r, s.ModifiedRegisters()...)

	}
	return r
}

type EQ struct {
	Left, Right, Dst Register
}

func (o EQ) String() string {
	return fmt.Sprintf("EQ %s, %s, %s", o.Left, o.Right, o.Dst)
}

func (o EQ) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o EQ) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type NEQ struct {
	Left, Right, Dst Register
}

func (o NEQ) String() string {
	return fmt.Sprintf("NEQ %s, %s, %s", o.Left, o.Right, o.Dst)
}
func (o NEQ) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o NEQ) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type GEQ struct {
	Left, Right, Dst Register
}

func (o GEQ) String() string {
	return fmt.Sprintf("GEQ %s, %s, %s", o.Left, o.Right, o.Dst)
}

func (o GEQ) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o GEQ) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type GT struct {
	Left, Right, Dst Register
}

func (o GT) String() string {
	return fmt.Sprintf("GT %s, %s, %s", o.Left, o.Right, o.Dst)
}

func (o GT) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o GT) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type LT struct {
	Left, Right, Dst Register
}

func (o LT) String() string {
	return fmt.Sprintf("LT %s, %s, %s", o.Left, o.Right, o.Dst)
}

func (o LT) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o LT) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}

type LTE struct {
	Left, Right, Dst Register
}

func (o LTE) String() string {
	return fmt.Sprintf("LEQ %s, %s, %s", o.Left, o.Right, o.Dst)
}
func (o LTE) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

func (o LTE) ModifiedRegisters() []Register {
	return []Register{o.Dst}
}
