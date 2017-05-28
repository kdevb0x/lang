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
func (c EqualityComparison) Type() Type {
	return "bool"
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

func (c NotEqualsComparison) Type() Type {
	return "bool"
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

func (c GreaterComparison) Type() Type {
	return "bool"
}

type GreaterOrEqualComparison struct {
	Left, Right Value
}

func (gc GreaterOrEqualComparison) BoolValue() bool {
	// This method is mostly a sentinal, the value returned doesn't matter
	// and since left and right are interfaces, > doesn't exist.
	return true //gc.Left >= gc.Right
}

func (c GreaterOrEqualComparison) Type() Type {
	return "bool"
}

func (n GreaterOrEqualComparison) Node() Node {
	return n
}

func (n GreaterOrEqualComparison) Value() interface{} {
	return n.BoolValue()
}

type LessThanOrEqualComparison struct {
	Left, Right Value
}

func (n LessThanOrEqualComparison) Node() Node {
	return n
}

func (n LessThanOrEqualComparison) Value() interface{} {
	return n.BoolValue()
}

func (gc LessThanOrEqualComparison) BoolValue() bool {
	// This method is mostly a sentinal, the value returned doesn't matter
	// and since left and right are interfaces, > doesn't exist.
	return true //gc.Left >= gc.Right
}

func (c LessThanOrEqualComparison) Type() Type {
	return "bool"
}

type LessThanComparison struct {
	Left, Right Value
}

func (n LessThanComparison) Node() Node {
	return n
}

func (n LessThanComparison) Value() interface{} {
	return n.BoolValue()
}

func (gc LessThanComparison) BoolValue() bool {
	// This method is mostly a sentinal, the value returned doesn't matter
	// and since left and right are interfaces, > doesn't exist.
	return true //gc.Left >= gc.Right
}

func (c LessThanComparison) Type() Type {
	return "bool"
}
