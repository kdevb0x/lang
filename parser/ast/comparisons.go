package ast

import (
	"fmt"
)

type EqualityComparison struct {
	Left, Right Value
}

func (ec EqualityComparison) BoolValue() bool {
	return ec.Left == ec.Right
}

func (ec EqualityComparison) Value() interface{} {
	return ec.BoolValue()
}

func (n EqualityComparison) Node() Node {
	return n
}

func (n EqualityComparison) String() string {
	return fmt.Sprintf("EqualityComparison{%v == %v}", n.Left, n.Right)
}

type NotEqualsComparison struct {
	Left, Right Value
}

func (ec NotEqualsComparison) BoolValue() bool {
	return ec.Left != ec.Right
}

func (ec NotEqualsComparison) Value() interface{} {
	return ec.BoolValue()
}

func (n NotEqualsComparison) Node() Node {
	return n
}

func (n NotEqualsComparison) String() string {
	return fmt.Sprintf("NotEqualsComparison{%v == %v}", n.Left, n.Right)
}

type GreaterComparison struct {
	Left, Right Value
}

func (gc GreaterComparison) BoolValue() bool {
	// This method is mostly a sentinal, the value returned doesn't matter
	// and since left and right are interfaces, > doesn't exist.
	return true //gc.Left > gc.Right
}

func (n GreaterComparison) Node() Node {
	return n
}

func (n GreaterComparison) Value() interface{} {
	return n.BoolValue()
}

type GreaterOrEqualComparison struct {
	Left, Right Value
}

func (gc GreaterOrEqualComparison) BoolValue() bool {
	// This method is mostly a sentinal, the value returned doesn't matter
	// and since left and right are interfaces, > doesn't exist.
	return true //gc.Left >= gc.Right
}

func (n GreaterOrEqualComparison) Node() Node {
	return n
}

func (n GreaterOrEqualComparison) Value() interface{} {
	return n.BoolValue()
}
