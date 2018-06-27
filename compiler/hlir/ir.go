package hlir

import (
	"fmt"

	"github.com/driusan/lang/parser/ast"
)

var debug = true

type Func struct {
	Name      string
	Body      []Opcode
	NumArgs   uint
	NumLocals uint
}

type Register interface{}

// Registers for arguments to be passed to the
// next function call.
type FuncCallArg uint

func (fa FuncCallArg) String() string {
	return fmt.Sprintf("FA%d", fa)
	// return fmt.Sprintf("FA%d (%v)", fa.Id, fa.Info)
}

// Denotes the return value of this function.
type FuncRetVal uint

func (fa FuncRetVal) String() string {
	return fmt.Sprintf("FR%d", fa)
}

// Denotes the return of the last function call.
type LastFuncCallRetVal struct {
	CallNum uint
	RetNum  uint
}

func (fa LastFuncCallRetVal) String() string {
	return fmt.Sprintf("CV(%d,%d)", fa.CallNum, fa.RetNum)
}

// Arguments to this function.
type FuncArg struct {
	Id        uint
	Reference bool
}

func (fa FuncArg) String() string {
	if debug {
		return fmt.Sprintf("P%d (%v)", fa.Id, fa.Reference)
	}
	return fmt.Sprintf("P%d", fa.Id)

}

type Pointer struct {
	Register
}

func (p Pointer) String() string {
	return fmt.Sprintf("&%v", p.Register)
}

// Registers for local variables
type LocalValue uint

func (lv LocalValue) String() string {
	return fmt.Sprintf("LV%d", lv)
}

// A TempValue is for a temporary calculation. It lives in a register,
// but never makes it to the stack. It's mostly for intermediate calculations
// such as the "x + 1" in "let y = x + 1"
type TempValue uint

func (lv TempValue) String() string {
	return fmt.Sprintf("TV%d", lv)
}

type IntLiteral int

func (il IntLiteral) String() string {
	return fmt.Sprintf("$%d", il)
}

type StringLiteral string

func (sl StringLiteral) String() string {
	return `$"` + string(sl) + `"`
}

// An Offset denotes a memory location which is offset from a base address.
// This is primarily for indexing into slices or arrays.
type Offset struct {
	// The register holding the offset from the base in bytes.
	Offset Register
	// The size of the type being offset.
	Scale IntLiteral
	// The register holding the base address to be offset from.
	Base Register

	Container ast.VarWithType
}

func (o Offset) String() string {
	if debug {
		return fmt.Sprintf("&(%v+%v*%v (%v))", o.Base, o.Offset, o.Scale, o.Container)
	}
	return fmt.Sprintf("&(%v+%v*%v)", o.Base, o.Offset, o.Scale)
}

type SliceBasePointer struct {
	Register
}

func (o SliceBasePointer) String() string {
	return fmt.Sprintf("*(%v)", o.Register)
}
