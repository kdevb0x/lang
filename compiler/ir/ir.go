package ir

import (
	"fmt"
)

type Func struct {
	Name    string
	Body    []Opcode
	NumArgs uint
}

type Register interface{}

// Registers for arguments to be passed to the
// next function call.
type FuncCallArg uint

func (fa FuncCallArg) String() string {
	return fmt.Sprintf("FA%d", fa)
}

type FuncRetVal uint

func (fa FuncRetVal) String() string {
	return fmt.Sprintf("FR%d", fa)
}

// Arguments to this function.
type FuncArg uint

func (fa FuncArg) String() string {
	return fmt.Sprintf("P%d", fa)
}

// Registers for local variables
type LocalValue uint

func (lv LocalValue) String() string {
	return fmt.Sprintf("LV%d", lv)
}

// Register for variables that were passed to this
// function call.
type LocalArg struct{}

type IntLiteral int

func (il IntLiteral) String() string {
	return fmt.Sprintf("$%d", il)
}

type StringLiteral string

func (sl StringLiteral) String() string {
	return `$"` + string(sl) + `"`
}
