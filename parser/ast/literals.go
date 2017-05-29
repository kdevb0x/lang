package ast

import (
	"fmt"
)

func IsLiteral(v Value) bool {
	switch v.(type) {
	case IntLiteral, BoolLiteral, StringLiteral:
		return true
	default:
		return false
	}
}

func IsCompatibleType(t TypeDefn, v Value) bool {
	switch t2 := v.(type) {
	case BoolLiteral:
		return t.ConcreteType == "bool"
	case IntLiteral:
		switch t.ConcreteType {
		case "int":
			return true
		default:
			return false
		}
	case StringLiteral:
		return t.ConcreteType == "string"
	default:
		println(t2)
		return false
	}
}

type StringLiteral string

func (v StringLiteral) Value() interface{} {
	return v
}

func (s StringLiteral) Node() Node {
	return s
}

func (s StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral(%v)", string(s))
}

func (s StringLiteral) Type() Type {
	return "string"
}

type IntLiteral int64

func (v IntLiteral) Value() interface{} {
	return v
}

func (s IntLiteral) Node() Node {
	return s
}

func (i IntLiteral) String() string {
	return fmt.Sprintf("IntLiteral(%d)", i)
}

func (i IntLiteral) Type() Type {
	return "int"
}

type BoolLiteral bool

func (v BoolLiteral) BoolValue() bool {
	return bool(v)
}
func (v BoolLiteral) Value() interface{} {
	return v
}

func (b BoolLiteral) Node() Node {
	return b
}

func (b BoolLiteral) String() string {
	if b {
		return "BoolLiteral(true)"
	}
	return "BoolLiteral(false)"
}

func (b BoolLiteral) Type() Type {
	return "bool"
}
