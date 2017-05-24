package ast

import (
	"fmt"
)

type VarWithType struct {
	Name Variable
	Type string
}

func (vt VarWithType) Node() Node {
	return vt
}

type Node interface {
	Node() Node
	// String() string
}

type MutStmt struct {
	Var          VarWithType
	InitialValue Value
}
type LetStmt struct {
	Var   VarWithType
	Value Value
}

func (s LetStmt) Node() Node {
	return s
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
	return fmt.Sprintf("MutStmt{Name: %v, InitialValue: %v}", ms.Var.Name.String(), ms.InitialValue)
}

func (ms MutStmt) Node() Node {
	return ms
}

type BoolValue interface {
	Value
	BoolValue() bool
}

type Value interface {
	Node
	Value() interface{}
}

type AssignmentOperator struct {
	Variable Variable
	Value    Value
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

type Variable string

func (v Variable) Value() interface{} {
	return v
}

func (v Variable) String() string {
	return fmt.Sprintf("Variable(%s)", string(v))
}

func (v Variable) Node() Node {
	return v
}
