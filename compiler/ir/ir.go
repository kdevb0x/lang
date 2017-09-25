package ir

import (
	"fmt"
	"github.com/driusan/lang/parser/ast"
)

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
	return fmt.Sprintf("FR%d", fa.Id)
	//	return fmt.Sprintf("FR%d (%v)", fa.Id, fa.Info)
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
	return fmt.Sprintf("P%d", fa.Id)
	//return fmt.Sprintf("P%d (%v)", fa.Id, fa.Info)

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
	return fmt.Sprintf("LV%d", lv.Id)
	// return fmt.Sprintf("LV%d (%v)", lv.Id, lv.Info)
}

func (lv LocalValue) Size() int {
	return lv.Info.Size
}

func (lv LocalValue) Signed() bool {
	return lv.Info.Signed
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
