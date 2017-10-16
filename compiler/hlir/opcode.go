package hlir

import (
	"fmt"
)

type Opcode interface {
	String() string
}

type RET struct{}

func (r RET) String() string {
	return "RET\n"
}

func (o RET) Registers() []Register {
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

type MOV struct {
	Src, Dst Register
}

func (m MOV) String() string {
	return fmt.Sprintf("MOV %v, %v\n", m.Src, m.Dst)
}
func (o MOV) Registers() []Register {
	return []Register{o.Src, o.Dst}
}

type ADD struct {
	Left, Right, Dst Register
}

func (o ADD) String() string {
	return fmt.Sprintf("ADD %v + %v => %v\n", o.Left, o.Right, o.Dst)
}

type SUB struct {
	Left, Right, Dst Register
}

func (o SUB) String() string {
	return fmt.Sprintf("SUB %s - %s => %v\n", o.Left, o.Right, o.Dst)
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

type MUL struct {
	Left, Right, Dst Register
}

func (o MUL) String() string {
	return fmt.Sprintf("MUL %s * %s => %s", o.Left, o.Right, o.Dst)
}

func (o MUL) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

type MOD struct {
	Left, Right, Dst Register
}

func (o MOD) String() string {
	return fmt.Sprintf("MOD %s, %s, %s\n", o.Left, o.Right, o.Dst)
}

func (o MOD) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

type ControlFlow struct {
	Condition
	Body []Opcode
}

type IF struct {
	ControlFlow
	ElseBody []Opcode
}

func (o IF) String() string {
	return fmt.Sprintf("IF %s (%s) ELSE (%s)\n", o.Condition, o.Body, o.ElseBody)
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

type LOOP ControlFlow

func (o LOOP) String() string {
	return fmt.Sprintf("LOOP %v (%s)\n", o.Condition, o.Body)
}

type JumpTable []ControlFlow

func (jt JumpTable) String() string {
	v := ""
	for _, s := range jt {
		v += s.String() + fmt.Sprintf(" %v\n", s.Body)

	}
	return fmt.Sprintf("JumpTable (\n%v)", v)
}

type EQ struct {
	Left, Right, Dst Register
}

func (o EQ) String() string {
	return fmt.Sprintf("EQ %s, %s, %s", o.Left, o.Right, o.Dst)
}

type NEQ struct {
	Left, Right, Dst Register
}

func (o NEQ) String() string {
	return fmt.Sprintf("NEQ %s, %s, %s", o.Left, o.Right, o.Dst)
}

type GEQ struct {
	Left, Right, Dst Register
}

func (o GEQ) String() string {
	return fmt.Sprintf("GEQ %s, %s, %s", o.Left, o.Right, o.Dst)
}

type GT struct {
	Left, Right, Dst Register
}

func (o GT) String() string {
	return fmt.Sprintf("GT %s, %s, %s", o.Left, o.Right, o.Dst)
}

type LT struct {
	Left, Right, Dst Register
}

func (o LT) String() string {
	return fmt.Sprintf("LT %s, %s, %s", o.Left, o.Right, o.Dst)
}
