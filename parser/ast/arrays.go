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

type SliceType struct {
	Base Type
}

func (a SliceType) Type() string {
	return fmt.Sprintf("[]%v", a.Base.Type())
}

func (a SliceType) Node() Node {
	return a
}
func (a SliceType) String() string {
	return fmt.Sprintf("SliceType{[]%v}", a.Base.Type())
}
