package ast

import (
	"fmt"
	"strings"
)

type Type interface {
	Node
	TypeName() string
	Info() TypeInfo
	Components() []Type
}

type TypeDef interface {
	TypeDefn() TypeDef
}
type Assignable interface {
	CanAssign() bool
}

type EnumOption struct {
	Constructor string
	Parameters  []string
	ParentType  Type
}

func (eo EnumOption) Node() Node {
	return eo
}

func (eo EnumOption) Value() interface{} {
	return eo.Constructor
}

func (eo EnumOption) Type() Type {
	return eo.ParentType
}

func (eo EnumOption) String() string {
	return fmt.Sprintf("EnumOption{%v, Parameters: %v ParentType: %v}", eo.Constructor, eo.Parameters, eo.ParentType)
}
func (eo EnumOption) Info() TypeInfo {
	panic("Unhandled info")
}

func (e EnumOption) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type EnumValue struct {
	Constructor EnumOption
	Parameters  []Value
}

func (ev EnumValue) Value() interface{} {
	return ev.Constructor.Type()
}

func (ev EnumValue) Node() Node {
	return ev
}
func (ev EnumValue) TypeName() string {
	base := ev.Constructor.Type().TypeName()
	for _, a := range ev.Parameters {
		base += " " + a.Type().TypeName()
	}
	return base
}

func (ev EnumValue) Type() Type {
	return ev.Constructor.Type() //UserType{TypeLiteral("int64", ev.TypeName()}
}

func (ev EnumValue) Info() TypeInfo {
	panic("Not implemented")
}
func (ev EnumValue) String() string {
	return fmt.Sprintf("EnumValue{%v, Parameters: %v}", ev.Constructor, ev.Parameters)
}

func (e EnumValue) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type VarWithType struct {
	Name      Variable
	Typ       Type
	Reference bool
}

func (vt VarWithType) CanAssign() bool {
	return true
}

func (vt VarWithType) Type() Type {
	return vt.Typ
}

func (vt VarWithType) Node() Node {
	return vt
}

func (v VarWithType) Value() interface{} {
	return v.Name
}

func (v VarWithType) BoolValue() bool {
	return true
}

func (v VarWithType) String() string {
	return fmt.Sprintf("VarWithType{%v %v %v}", v.Name, v.Typ, v.Reference)
}

func (v VarWithType) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v", nTabs(lvl), v.Name)
}

type Node interface {
	Node() Node
	PrettyPrint(lvl int) string
}

type MutStmt struct {
	Var          VarWithType
	InitialValue Value
}

type TypeDefn struct {
	Name         string
	ConcreteType Type
	Parameters   []string
}

func (t TypeDefn) Node() Node {
	return t
}

func (t TypeDefn) TypeDefn() TypeDef {
	return t
}

func (t TypeDefn) String() string {
	return fmt.Sprintf("TypeDefn{Name: %v Type: %v Parameters: %v}", t.Name, t.ConcreteType, t.Parameters)
}

func (t TypeDefn) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type EnumTypeDefn struct {
	Name           string
	Options        []EnumOption
	Parameters     []Type
	ExpectedParams int
}

func (t EnumTypeDefn) Node() Node {
	return t
}

func (t EnumTypeDefn) String() string {
	return fmt.Sprintf("EnumTypeDefn{%v, Options: %v Parameters(%d): %v}", t.Name, t.Options, t.ExpectedParams, t.Parameters)

}

func (t EnumTypeDefn) TypeDefn() TypeDef {
	return t
}

func (t EnumTypeDefn) Type() Type {
	return t
}

func (e EnumTypeDefn) Info() TypeInfo {
	t := TypeInfo{8, false}
	for _, c := range e.Parameters {
		t.Size += c.Info().Size
	}
	return t
}

func (t EnumTypeDefn) TypeName() string {
	ret := t.Name
	for _, p := range t.Parameters {
		ret += " " + p.TypeName()
	}
	return ret
}

func (t EnumTypeDefn) PrettyPrint(lvl int) string {
	ret := fmt.Sprintf("%v%v", nTabs(lvl), t.Name)
	for _, p := range t.Parameters {
		ret += " " + p.TypeName()
	}
	return ret
}

func (t EnumTypeDefn) Components() []Type {
	// 1 piece for the variant, plus each parameter
	ti := TypeLiteral("int")
	return append([]Type{ti}, t.Parameters...)
}

func (m MutStmt) Type() Type {
	return m.Var.Type()
}
func (m MutStmt) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type LetStmt struct {
	Var VarWithType
	Val Value
}

func (s LetStmt) Node() Node {
	return s
}

func (l LetStmt) Type() Type {
	return l.Var.Type()
}

func (l LetStmt) TypeName() string {
	return l.Var.Type().TypeName()
}

func (ls LetStmt) String() string {
	return fmt.Sprintf("LetStmt{%v, Value: %v}", ls.Var, ls.Val)
}

func (ls LetStmt) Value() interface{} {
	return ls.Val.Value()
}

func (ls LetStmt) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type BlockStmt struct {
	Stmts []Node
}

func (b BlockStmt) Node() Node {
	return b
}

func (b BlockStmt) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

func (b BlockStmt) String() string {
	if len(b.Stmts) == 0 {
		return "BlockStmt{}"
	}
	ret := "BlockStmt{\n"
	for _, v := range b.Stmts {
		ret += fmt.Sprintf("%v\n", v)
	}
	return ret + "}"
}

func (ms MutStmt) String() string {
	return fmt.Sprintf("MutStmt{%v, InitialValue: %v}", ms.Var, ms.InitialValue)
}

func (ms MutStmt) Node() Node {
	return ms
}

func (ms MutStmt) TypeName() string {
	return ms.Var.Type().TypeName()
}

type BoolValue interface {
	Value
	BoolValue() bool
}

type Value interface {
	Node
	Type() Type
	Value() interface{}
}

type AssignmentOperator struct {
	Variable Assignable
	Value    Value
}

func (ao AssignmentOperator) String() string {
	return fmt.Sprintf("AssignmentOperator{%v = %v}", ao.Variable, ao.Value)
}
func (ao AssignmentOperator) Node() Node {
	return ao
}

func (ao AssignmentOperator) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type AdditionOperator struct {
	Left, Right Value
}

func (ao AdditionOperator) Node() Node {
	return ao
}

func (ao AdditionOperator) Value() interface{} {
	return true
}

func (ao AdditionOperator) String() string {
	return fmt.Sprintf("(%v + %v)", ao.Left, ao.Right)
}

func (ao AdditionOperator) Type() Type {
	return ao.Left.Type()
}

func (o AdditionOperator) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v + %v", nTabs(lvl), o.Left.PrettyPrint(0), o.Right.PrettyPrint(0))
}

type SubtractionOperator struct {
	Left, Right Value
}

func (so SubtractionOperator) Node() Node {
	return so
}

func (so SubtractionOperator) Value() interface{} {
	return true
}

func (o SubtractionOperator) String() string {
	return fmt.Sprintf("(%v - %v)", o.Left, o.Right)
}

func (o SubtractionOperator) Type() Type {
	return o.Left.Type()
}

func (o SubtractionOperator) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v - %v", nTabs(lvl), o.Left.PrettyPrint(0), o.Right.PrettyPrint(0))
}

type MulOperator struct {
	Left, Right Value
}

func (mo MulOperator) Value() interface{} {
	return 4
}

func (mo MulOperator) Node() Node {
	return mo
}

func (o MulOperator) String() string {
	return fmt.Sprintf("(%v * %v)", o.Left, o.Right)
}

func (o MulOperator) Type() Type {
	return o.Left.Type()
}

func (o MulOperator) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v * %v", nTabs(lvl), o.Left.PrettyPrint(0), o.Right.PrettyPrint(0))
}

type DivOperator struct {
	Left, Right Value
}

func (mo DivOperator) Value() interface{} {
	return 4
}

func (mo DivOperator) Node() Node {
	return mo
}
func (o DivOperator) String() string {
	return fmt.Sprintf("(%v / %v)", o.Left, o.Right)
}

func (o DivOperator) Type() Type {
	return o.Left.Type()
}

func (o DivOperator) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v / %v", nTabs(lvl), o.Left.PrettyPrint(0), o.Right.PrettyPrint(0))
}

type ModOperator struct {
	Left, Right Value
}

func (mo ModOperator) Value() interface{} {
	return 3
}

func (mo ModOperator) Node() Node {
	return mo
}

func (m ModOperator) String() string {
	return fmt.Sprintf("ModOperator{%v mod %v}", m.Left, m.Right)
}

func (m ModOperator) Type() Type {
	return m.Left.Type()
}
func (o ModOperator) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v %c %v", nTabs(lvl), o.Left.PrettyPrint(0), '%', o.Right.PrettyPrint(0))
}

type Variable string

func (v Variable) String() string {
	return string(v)
	//return fmt.Sprintf("Variable(%s)", string(v))
}

func (v Variable) Node() Node {
	return v
}

func (v Variable) PrettyPrint(lvl int) string {
	return nTabs(lvl) + string(v)
}

type TypeLiteral string

func (tl TypeLiteral) TypeName() string {
	return string(tl)
}
func (tl TypeLiteral) Info() TypeInfo {
	switch tl {
	case "bool":
		return TypeInfo{1, false}
	case "byte", "uint8":
		return TypeInfo{1, false}
	case "int8":
		return TypeInfo{1, true}
	case "uint16":
		return TypeInfo{2, false}
	case "int16":
		return TypeInfo{2, true}
	case "uint32":
		return TypeInfo{4, false}
	case "int32":
		return TypeInfo{4, true}
	case "uint64":
		return TypeInfo{8, false}
	case "int64":
		return TypeInfo{8, true}
	case "uint":
		return TypeInfo{0, false}
	case "int":
		return TypeInfo{0, true}
	case "string":
		// 8 for length, 8 for ptr
		return TypeInfo{16, false}
	default:
		panic("Unhandled type literal " + string(tl))
	}
}

func (tl TypeLiteral) Node() Node {
	return tl
}

func (tl TypeLiteral) PrettyPrint(lvl int) string {
	return nTabs(lvl) + string(tl)
}

func (tl TypeLiteral) Components() []Type {
	if tl == "string" {
		// One length, one pointer
		return []Type{TypeLiteral("uint64"), TypeLiteral("uint")}
	}
	return []Type{tl}
}

type Brackets struct {
	Val Value
}

func (b Brackets) Value() interface{} {
	return b.Val
}
func (b Brackets) Node() Node {
	return b
}

func (b Brackets) Type() Type {
	return b.Val.Type()
}
func (b Brackets) String() string {
	return fmt.Sprintf("Brackets{%v}", b.Val)
}

func (b Brackets) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v(%v)", nTabs(lvl), b.Val.PrettyPrint(0))
}

type Cast struct {
	Val Value
	Typ Type
}

func (c Cast) Value() interface{} {
	return c.Val
}
func (c Cast) Node() Node {
	return c
}

func (c Cast) Type() Type {
	return c.Typ
}

func (c Cast) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%vcast (%v) as %v", nTabs(lvl), c.Val.PrettyPrint(0), c.Typ.PrettyPrint(0))
}

type SumType []Type

func (s SumType) TypeName() string {
	if len(s) == 0 {
		return ""
	}
	var ret strings.Builder
	for i, o := range s {
		if i == 0 {
			fmt.Fprintf(&ret, "%v", o.TypeName())
		} else {
			fmt.Fprintf(&ret, " | %v", o.TypeName())
		}
	}
	return ret.String()
}

func (s SumType) Node() Node {
	return s
}

func (s SumType) String() string {
	return fmt.Sprintf("SumType{%v}", s.TypeName())
}

func (s SumType) PrettyPrint(lvl int) string {
	return nTabs(lvl) + s.TypeName()
}

func (s SumType) Info() TypeInfo {
	// Sum types are the size of the largest possible type, plus 1 word to hold the variant
	// that's currently being worked with.
	ti := TypeInfo{8, false}
	var max TypeInfo
	for _, v := range s {
		subtype := v.Info()
		if subtype.Size > max.Size {
			max = subtype
		}
	}
	if max.Size == 0 {
		max.Size = 8
	}
	ti.Size += max.Size
	return ti
}

func (s SumType) Components() []Type {
	var possible []Type
	for _, subtype := range s {
		if sub := subtype.Components(); len(sub) >= len(possible) {
			possible = sub
		}
	}
	// One for the variant stored, then the biggest possible variant
	return append([]Type{TypeLiteral("uint64")}, possible...)
}

type TupleValue []Value

func (tv TupleValue) Node() Node {
	return tv
}

func (tv TupleValue) PrettyPrint(lvl int) string {
	rv := nTabs(lvl) + "("
	for i, v := range tv {
		rv += v.PrettyPrint(0)
		if i != len(tv)-1 {
			rv += ", "
		}
	}
	return rv + ")"
}

func (tv TupleValue) Type() Type {
	var rv TupleType
	for _, c := range tv {
		rv = append(rv, VarWithType{"", c.Type(), false})
	}
	return rv
}
func (tv TupleValue) Value() interface{} {
	return nil
}

type UserType struct {
	Typ  Type
	Name string
}

func (ut UserType) String() string {
	return fmt.Sprintf("UserType{Name: %v, Type: %v}", ut.Name, ut.Typ)
}

func (ut UserType) TypeName() string {
	return ut.Name
}

func (ut UserType) Type() Type {
	return ut
}

func (ut UserType) Node() Node {
	return ut
}

func (ut UserType) Components() []Type {
	return ut.Typ.Components()
}

func (ut UserType) Info() TypeInfo {
	return ut.Typ.Info()
}

func (ut UserType) PrettyPrint(lvl int) string {
	return nTabs(lvl) + ut.Name
}
