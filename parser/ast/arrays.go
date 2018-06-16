package ast

import (
	"fmt"
)

type ArrayType struct {
	Base Type
	Size IntLiteral
}

func (a ArrayType) Type() string {
	return fmt.Sprintf("[%d]%v", a.Size, a.Base.Type())
}

func (a ArrayType) Node() Node {
	return a
}

func (a ArrayType) String() string {
	return fmt.Sprintf("ArrayType{[%d]%v}", a.Size, a.Base.Type())
}
func (a ArrayType) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v[%d]%v", nTabs(lvl), a.Size, a.Base.Type())
}

type ArrayLiteral []Value

func (v ArrayLiteral) Type() string {
	return fmt.Sprintf("[%v]%v", len(v), v[0].Type())
}

func (v ArrayLiteral) Node() Node {
	return v
}

func (v ArrayLiteral) Value() interface{} {
	return v
}

func (v ArrayLiteral) String() string {
	return fmt.Sprintf("ArrayLiteral{[%v]%v (%v)}", len(v), v[0].Type(), []Value(v))
}

func (v ArrayLiteral) PrettyPrint(lvl int) string {
	panic("PrettyPrint not implemented")
}

type ArrayValue struct {
	Base  VarWithType
	Index Value
}

func (v ArrayValue) Type() string {
	switch bt := v.Base.Typ.(type) {
	case ArrayType:
		return bt.Base.Type()
	case SliceType:
		return bt.Base.Type()
	default:
		panic("Attempt to index on non-array")
	}
}

func (v ArrayValue) String() string {
	return fmt.Sprintf("ArrayValue{%v[%v]}", v.Base, v.Index)
}
func (v ArrayValue) Node() Node {
	return v
}

func (v ArrayValue) Value() interface{} {
	return v
}

func (v ArrayValue) CanAssign() bool {
	return true
}

func (v ArrayValue) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v[%v]", nTabs(lvl), v.Base.Name, v.Index.PrettyPrint(0))
}

type SliceType struct {
	Base Type
}

func (a SliceType) Type() string {
	return fmt.Sprintf("[]%v", a.Base.Type())
}
func (a SliceType) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v[]%v", nTabs(lvl), a.Base.Type())
}

func (a SliceType) Node() Node {
	return a
}
func (a SliceType) String() string {
	return fmt.Sprintf("SliceType{[]%v}", a.Base.Type())
}
