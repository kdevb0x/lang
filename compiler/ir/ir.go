package ir

import (
	"fmt"
	"github.com/driusan/lang/parser/ast"
)

var debug = false

type Func struct {
	Name    string
	Body    []Opcode
	NumArgs uint
}

type Register interface {
	Size() int
	Signed() bool
}

// Registers for arguments to be passed to the
// next function call.

type FuncCallArg struct {
	Id   int
	Info ast.TypeInfo
}

func (fa FuncCallArg) String() string {
	return fmt.Sprintf("FA%d", fa.Id)
	// return fmt.Sprintf("FA%d (%v)", fa.Id, fa.Info)
}

func (fa FuncCallArg) Size() int {
	// Not sure what this should be.
	return fa.Info.Size
}

func (fa FuncCallArg) Signed() bool {
	return fa.Info.Signed
}

type FuncRetVal struct {
	Id   uint
	Info ast.TypeInfo
}

func (fa FuncRetVal) String() string {
	if debug {
		return fmt.Sprintf("FR%d (%v)", fa.Id, fa.Info)
	}
	return fmt.Sprintf("FR%d", fa.Id)
}

func (fa FuncRetVal) Size() int {
	return fa.Info.Size

}

func (fa FuncRetVal) Signed() bool {
	return fa.Info.Signed
}

// Arguments to this function.
type FuncArg struct {
	Id        uint
	Info      ast.TypeInfo
	Reference bool
}

func (fa FuncArg) String() string {
	if debug {
		return fmt.Sprintf("P%d (%v)", fa.Id, fa.Info)
	}
	return fmt.Sprintf("P%d", fa.Id)

}

func (fa FuncArg) Size() int {
	return fa.Info.Size
}
func (fa FuncArg) Signed() bool {
	return fa.Info.Signed
}

type Pointer struct {
	Register
}

func (p Pointer) String() string {
	return fmt.Sprintf("&%v", p.Register)
}

// Registers for local variables
type LocalValue struct {
	Id   uint
	Info ast.TypeInfo
}

func (lv LocalValue) String() string {
	if debug {
		return fmt.Sprintf("LV%d (%v)", lv.Id, lv.Info)
	}
	return fmt.Sprintf("LV%d", lv.Id)
}

func (lv LocalValue) Size() int {
	return lv.Info.Size
}

func (lv LocalValue) Signed() bool {
	return lv.Info.Signed
}

// A TempValue is for a temporary calculation. It lives in a register,
// but never makes it to the stack. It's mostly for intermediate calculations
// such as the "x + 1" in "let y = x + 1"
type TempValue uint

func (lv TempValue) String() string {
	return fmt.Sprintf("TV%d", lv)
}

func (lv TempValue) Size() int {
	return 8
}

func (lv TempValue) Signed() bool {
	return true
}

// An unsigned TempValue
type UTempValue uint

func (lv UTempValue) String() string {
	return fmt.Sprintf("TV%d", lv)
}

func (lv UTempValue) Size() int {
	return 8
}

func (lv UTempValue) Signed() bool {
	return false
}

type IntLiteral int

func (il IntLiteral) String() string {
	return fmt.Sprintf("$%d", il)
}

func (il IntLiteral) Size() int {
	// FIXME: What's the right value for this?
	return 0
}

func (il IntLiteral) Signed() bool {
	return true
}

type StringLiteral string

func (sl StringLiteral) String() string {
	return `$"` + string(sl) + `"`
}

func (sl StringLiteral) Size() int {
	return 0
}
func (l StringLiteral) Signed() bool {
	return false
}

// An Offset denotes a memory location which is offset from a base address.
// This is primarily for indexing into slices or arrays.
type Offset struct {
	// The register holding the offset from the base in bytes.
	Offset Register
	// The register holding the base address to be offset from.
	Base Register
}

func (o Offset) Signed() bool {
	return false
}

func (o Offset) Size() int {
	return 0
}

func (o Offset) String() string {
	return fmt.Sprintf("&(%v+%v)", o.Base, o.Offset)
}
