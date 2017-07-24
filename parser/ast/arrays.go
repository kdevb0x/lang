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
	bt, ok := v.Base.Typ.(ArrayType)
	if !ok {
		panic("Attempt to index on non-array")
	}
	return bt.Base.Type()
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
