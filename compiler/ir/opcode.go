package ir

import (
	"fmt"
)

type Opcode interface {
	String() string
	Registers() []Register
}

type RET struct{}

func (r RET) String() string {
	return "RET\n"
}

func (o RET) Registers() []Register {
	return nil
}

type CALL struct {
	FName    Fname
	Args     []Register
	TailCall bool
}

func (c CALL) String() string {
	return fmt.Sprintf("CALL %v (%v)\n", c.FName, c.Args)
}

func (o CALL) Registers() []Register {
	return o.Args
}

type Label string

func (l Label) String() string {
	return fmt.Sprintf("%v:\n", string(l))
}

func (l Label) Inline() string {
	return string(l)
}
func (o Label) Registers() []Register {
	return nil
}

type Fname string

type MOV struct {
	Src, Dst Register
}

func (m MOV) String() string {
	return fmt.Sprintf("MOV %v, %v\n", m.Src, m.Dst)
}
func (o MOV) Registers() []Register {
	return []Register{o.Src, o.Dst}
}

type JMP struct {
	Label
}

func (j JMP) String() string {
	return fmt.Sprintf("JMP %v\n", j.Label.Inline())
}

type ConditionalJump struct {
	Label
	Src, Dst Register
}
type JLE struct {
	ConditionalJump
}

func (j ConditionalJump) Registers() []Register {
	return []Register{j.Src, j.Dst}
}

func (j JLE) String() string {
	return fmt.Sprintf("JLE %v, %v, %v\n", j.ConditionalJump.Label.Inline(), j.Src, j.Dst)
}

type JGE struct {
	ConditionalJump
}

func (j JGE) String() string {
	return fmt.Sprintf("JGE %v, %v, %v\n", j.Label.Inline(), j.Src, j.Dst)
}

type JG struct {
	ConditionalJump
}

func (j JG) String() string {
	return fmt.Sprintf("JG %v, %v, %v\n", j.Label.Inline(), j.Src, j.Dst)
}

type JL struct {
	ConditionalJump
}

func (j JL) String() string {
	return fmt.Sprintf("JL %v, %v, %v\n", j.Label.Inline(), j.Src, j.Dst)
}

type JNE struct {
	ConditionalJump
}

func (j JNE) String() string {
	return fmt.Sprintf("JNE %v, %v, %v\n", j.Label.Inline(), j.Src, j.Dst)
}

type JE struct {
	ConditionalJump
}

func (j JE) String() string {
	return fmt.Sprintf("JE %v, %v, %v\n", j.Label.Inline(), j.Src, j.Dst)
}

type ADD struct {
	Src, Dst Register
}

func (o ADD) String() string {
	return fmt.Sprintf("ADD %v, %v\n", o.Src, o.Dst)
}
func (o ADD) Registers() []Register {
	return []Register{o.Src, o.Dst}
}

type SUB struct {
	Src, Dst Register
}

func (o SUB) String() string {
	return fmt.Sprintf("SUB %s, %s\n", o.Src, o.Dst)
}

func (o SUB) Registers() []Register {
	return []Register{o.Src, o.Dst}
}

type DIV struct {
	Left, Right, Dst Register
}

func (o DIV) String() string {
	return fmt.Sprintf("DIV %s, %s, %s\n", o.Left, o.Right, o.Dst)
}

func (o DIV) Registers() []Register {
	return []Register{o.Left, o.Right, o.Dst}
}

type MUL struct {
	Left, Right, Dst Register
}

func (o MUL) String() string {
	return fmt.Sprintf("MUL %s, %s, %s\n", o.Left, o.Right, o.Dst)
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
