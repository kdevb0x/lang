package ast

import (
	"fmt"
)

type ArrayType struct {
	Base Type
	Size IntLiteral
}

func (a ArrayType) TypeName() string {
	return fmt.Sprintf("[%d]%v", a.Size, a.Base.TypeName())
}

func (a ArrayType) Node() Node {
	return a
}

func (a ArrayType) String() string {
	return fmt.Sprintf("ArrayType{[%d]%v}", a.Size, a.Base.TypeName())
}
func (a ArrayType) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v[%d]%v", nTabs(lvl), a.Size, a.Base.TypeName())
}

func (a ArrayType) Info() TypeInfo {
	baseinfo := a.Base.Info()
	return TypeInfo{
		baseinfo.Size * int(a.Size),
		baseinfo.Signed,
	}
}
func (a ArrayType) Components() []Type {
	var v []Type
	for i := 0; i < int(a.Size); i++ {
		v = append(v, a.Base)
	}
	return v
}

type ArrayLiteral []Value

func (v ArrayLiteral) TypeName() string {
	return fmt.Sprintf("[%v]%v", len(v), v[0].Type())
}

func (v ArrayLiteral) Node() Node {
	return v
}

func (v ArrayLiteral) Value() interface{} {
	return v
}

func (v ArrayLiteral) String() string {
	if len(v) == 0 {
		return fmt.Sprintf("ArrayLiteral{[%v]nil (%v)}", len(v), []Value(v))
	}
	return fmt.Sprintf("ArrayLiteral{[%v]%v (%v)}", len(v), v[0].Type(), []Value(v))
}

func (v ArrayLiteral) PrettyPrint(lvl int) string {
	panic("PrettyPrint not implemented")
}
func (v ArrayLiteral) Type() Type {
	return ArrayType{
		Base: v[0].Type(),
		Size: IntLiteral(len(v)),
	}
}

type ArrayValue struct {
	Base  VarWithType
	Index Value
}

func (v ArrayValue) TypeName() string {
	switch bt := v.Base.Typ.(type) {
	case ArrayType:
		return bt.Base.TypeName()
	case SliceType:
		return bt.Base.TypeName()
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

func (v ArrayValue) Type() Type {
	t := v.Base.Type()
	switch t2 := t.(type) {
	case ArrayType:
		return t2.Base
	case SliceType:
		return t2.Base
	default:
		panic("Attempt to use ArrayValue for non-indexable type")
	}
}

type SliceType struct {
	Base Type
}

func (a SliceType) TypeName() string {
	return fmt.Sprintf("[]%v", a.Base.TypeName())
}
func (a SliceType) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v[]%v", nTabs(lvl), a.Base.TypeName())
}

func (a SliceType) Node() Node {
	return a
}
func (a SliceType) String() string {
	return fmt.Sprintf("SliceType{[]%v}", a.Base.TypeName())
}
func (a SliceType) Info() TypeInfo {
	baseinfo := a.Base.Info()
	return TypeInfo{
		16, // 8 for size, 8 for base pointer
		baseinfo.Signed,
	}
}

func (a SliceType) Components() []Type {
	// One int64 for the size, one int for the pointer
	return []Type{TypeLiteral("uint64"), TypeLiteral("int")}
}

type Slice struct {
	Base Value
	Size IntLiteral
}

func (a Slice) Type() Type {
	return SliceType{
		Base: a.Base.Type(),
	}
}

/*
func (a Slice) TypeName() string {
	return fmt.Sprintf("[]%v", a.Start.Typ.TypeName())
}
*/
func (a Slice) PrettyPrint(lvl int) string {
	// FIXME: This is wrong.
	return fmt.Sprintf("%v[%v:%v]", nTabs(lvl), a.Base, a.Size)
}

func (a Slice) Node() Node {
	return a
}
func (a Slice) String() string {
	return fmt.Sprintf("Slice{Start: %v, Size: %v}", a.Base, a.Size)
}

func (a Slice) Value() interface{} {
	// FIXME: this is just a stub.
	return 0
}

/*
func (a Slice) Info() TypeInfo {
	baseinfo := a.Start.Typ.Info()
	return TypeInfo{
		16, // 8 for size, 8 for base pointer
		baseinfo.Signed,
	}
}

func (a Slice) Components() []Type {
	// One int64 for the size, one int for the pointer
	return []Type{TypeLiteral("uint64"), TypeLiteral("int")}
}
*/
