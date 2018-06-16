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

func (c EqualityComparison) Type() Type {
	return TypeLiteral("bool")
}

func (c EqualityComparison) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v == %v", nTabs(lvl), c.Left.PrettyPrint(0), c.Right.PrettyPrint(0))
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

func (n NotEqualsComparison) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v != %v", nTabs(lvl), n.Left.PrettyPrint(0), n.Right.PrettyPrint(0))
}

func (c NotEqualsComparison) Type() Type {
	return TypeLiteral("bool")
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

func (c GreaterComparison) TypeName() string {
	return "bool"
}

func (c GreaterComparison) Type() Type {
	return TypeLiteral("bool")
}

func (c GreaterComparison) String() string {
	return fmt.Sprintf("GreaterComparison{%v, %v}", c.Left, c.Right)
}
func (c GreaterComparison) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v > %v", nTabs(lvl), c.Left.PrettyPrint(0), c.Right.PrettyPrint(0))
}

type GreaterOrEqualComparison struct {
	Left, Right Value
}

func (gc GreaterOrEqualComparison) BoolValue() bool {
	// This method is mostly a sentinal, the value returned doesn't matter
	// and since left and right are interfaces, > doesn't exist.
	return true //gc.Left >= gc.Right
}

func (n GreaterOrEqualComparison) Type() Type {
	return TypeLiteral("bool")
}

func (n GreaterOrEqualComparison) Node() Node {
	return n
}

func (n GreaterOrEqualComparison) Value() interface{} {
	return n.BoolValue()
}

func (n GreaterOrEqualComparison) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v >= %v", nTabs(lvl), n.Left.PrettyPrint(0), n.Right.PrettyPrint(0))
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
	return TypeLiteral("bool")
}

func (n LessThanOrEqualComparison) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v <= %v", nTabs(lvl), n.Left.PrettyPrint(0), n.Right.PrettyPrint(0))
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
	return true
}

func (c LessThanComparison) Type() Type {
	return TypeLiteral("bool")
}

func (c LessThanComparison) String() string {
	return fmt.Sprintf("LessThanComparison{%v, %v}", c.Left, c.Right)
}

func (n LessThanComparison) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v < %v", nTabs(lvl), n.Left.PrettyPrint(0), n.Right.PrettyPrint(0))
}
