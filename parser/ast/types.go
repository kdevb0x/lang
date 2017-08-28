package ast

import (
	"fmt"
)

// FIXME: This should be eliminated
type Type interface {
	Node
	Type() string
}

type TypeDef interface {
	TypeDefn() TypeDef
}

type EnumOption struct {
	Constructor string
	Parameters  []Type
	ParentType  Type
}

func (eo EnumOption) Node() Node {
	return eo
}

func (eo EnumOption) Value() interface{} {
	return eo.Constructor
}
func (eo EnumOption) Type() string {
	return eo.ParentType.Type()
}

func (eo EnumOption) String() string {
	return fmt.Sprintf("EnumOption{%v, Parameters: %v ParentType: %v}", eo.Constructor, eo.Parameters, eo.ParentType)
}

type EnumValue struct {
	Constructor EnumOption
	Parameters  []Value
}

func (ev EnumValue) Value() interface{} {
	return ev.Constructor
}

func (ev EnumValue) Node() Node {
	return ev
}
func (ev EnumValue) Type() string {
	return ev.Constructor.Type()
}

func (ev EnumValue) String() string {
	return fmt.Sprintf("EnumValue{%v, Parameters: %v}", ev.Constructor, ev.Parameters)
}

type VarWithType struct {
	Name      Variable
	Typ       Type
	Reference bool
}

func (vt VarWithType) Type() string {
	if vt.Typ == nil {
		return ""
	}
	return vt.Typ.Type()
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

type Node interface {
	Node() Node
}

type MutStmt struct {
	Var          VarWithType
	InitialValue Value
}

type TypeDefn struct {
	Name         Type
	ConcreteType Type
	Parameters   []Type
}

func (t TypeDefn) Node() Node {
	return t
}

func (t TypeDefn) TypeDefn() TypeDef {
	return t
}

type SumTypeDefn struct {
	Name       Type
	Options    []EnumOption
	Parameters []Type
}

func (t SumTypeDefn) Node() Node {
	return t
}

func (t SumTypeDefn) String() string {
	return fmt.Sprintf("SumTypeDefn{%v, Options: %v}", t.Name, t.Options)

}

func (t SumTypeDefn) TypeDefn() TypeDef {
	return t
}

func (m MutStmt) Type() string {
	return m.Var.Type()
}

type LetStmt struct {
	Var   VarWithType
	Value Value
}

func (s LetStmt) Node() Node {
	return s
}

func (l LetStmt) Type() string {
	return l.Var.Type()
}

func (ls LetStmt) String() string {
	return fmt.Sprintf("LetStmt{%v, Value: %v}", ls.Var, ls.Value)
}

type BlockStmt struct {
	Stmts []Node
}

func (b BlockStmt) Node() Node {
	return b
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

type BoolValue interface {
	Value
	BoolValue() bool
}

type Value interface {
	Type
	Value() interface{}
}

type AssignmentOperator struct {
	Variable VarWithType
	Value    Value
}

func (ao AssignmentOperator) String() string {
	return fmt.Sprintf("AssignmentOperator{%v = %v}", ao.Variable, ao.Value)
}
func (ao AssignmentOperator) Node() Node {
	return ao
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

func (ao AdditionOperator) Type() string {
	return ao.Left.Type()
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

func (o SubtractionOperator) Type() string {
	return o.Left.Type()
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

func (o MulOperator) Type() string {
	if lt := o.Left.Type(); lt != "" {
		return lt
	}
	return o.Right.Type()
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

func (o DivOperator) Type() string {
	return o.Left.Type()
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

func (m ModOperator) Type() string {
	return m.Left.Type()
}

type Variable string

func (v Variable) String() string {
	return string(v)
	//return fmt.Sprintf("Variable(%s)", string(v))
}

func (v Variable) Node() Node {
	return v
}

type TypeLiteral string

func (tl TypeLiteral) Type() string {
	return string(tl)
}

func (tl TypeLiteral) Node() Node {
	return tl
}
