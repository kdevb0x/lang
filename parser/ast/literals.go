package ast

import (
	"fmt"
)

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

type BoolLiteral bool

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
