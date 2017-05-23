package ast

import (
	"fmt"
	"github.com/driusan/lang/parser/token"
	"reflect"
	"strconv"
)

// FIXME: This whole consumeValue flow is badly designed.
// This is going to have to be completely reworked.
func isInfixOperator(pos int, tokens []token.Token) bool {
	if pos >= len(tokens) {
		return false
	}
	switch t := tokens[pos].(type) {
	case token.Operator:
		switch t {
		case "+", "-", "*", "/", "%",
			"<=", "<", "==", ">", ">=", "!=", "=":
			return true
		}
	}
	return false
}

func consumeValue(start int, tokens []token.Token, c Context) (int, Value, error) {
	// FIXME: Implement this properly
	return consumeBoolValue(start, tokens, c)
}

func consumeIntValue(start int, tokens []token.Token, c Context) (int, Value, error) {
	var partial Value
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			if n, err := strconv.Atoi(t.String()); err == nil {
				partial = IntLiteral(n)
			} else if c.IsVariable(t.String()) {
				// If it's a defined variable, that takes
				// precedence as a symbol over a function
				// name.
				partial = Variable(t.String())
			} else if c.IsFunction(t.String()) {
				// if it's a function, call it.
				n, fc, err := consumeFuncCall(i, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				partial = fc
				i += n
			} else {
				// FIXME: Otherwise, it may still be a parameter.
				// Validate this.
				partial = Variable(t.String())
			}

			switch tokens[i+1] {
			case token.Operator("+"):
				n, right, err := consumeIntValue(i+2, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				i += n - 1
				partial = AdditionOperator{
					Left:  partial,
					Right: right,
				}
			case token.Operator("-"):
				n, right, err := consumeIntValue(i+2, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				i += n - 1

				partial = SubtractionOperator{
					Left:  partial,
					Right: right,
				}
			case token.Operator("%"):
				n, right, err := consumeIntValue(i+2, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				i += n - 1
				partial = ModOperation{
					A: partial,
					B: right,
				}
			}
		default:
			return i + 2 - start, partial, nil
		}
	}
	if partial == nil {
		return 0, nil, fmt.Errorf("No value")
	}
	panic("Don't know how many tokens were consumed")
}

func operatorPrecedence(op Value) int {
	switch op.(type) {
	case MulOperator, DivOperator, ModOperation:
		return 5
	case AdditionOperator, SubtractionOperator:
		return 4
	case EqualityComparison, NotEqualsComparison, GreaterOrEqualComparison, GreaterComparison, LessThanComparison, LessThanOrEqualComparison:
		return 3
	default:
		panic(fmt.Sprintf("Unhandled precedence %v", reflect.TypeOf(op)))
	}
}

func getLeft(node Value) Value {
	switch v := node.(type) {
	case AdditionOperator:
		return v.Left
	case SubtractionOperator:
		return v.Left
	case DivOperator:
		return v.Left
	case MulOperator:
		return v.Left
	case NotEqualsComparison:
		return v.Left
	case EqualityComparison:
		return v.Left
	case GreaterOrEqualComparison:
		return v.Left
	case GreaterComparison:
		return v.Left
	case LessThanOrEqualComparison:
		return v.Left
	case LessThanComparison:
		return v.Left
	default:
		panic(fmt.Sprintf("Unhandled node type in getLeft: %v", reflect.TypeOf(node)))
	}
}
func getRight(node Value) Value {
	switch v := node.(type) {
	case AdditionOperator:
		return v.Right
	case SubtractionOperator:
		return v.Right
	case DivOperator:
		return v.Right
	case MulOperator:
		return v.Right
	case NotEqualsComparison:
		return v.Right
	case EqualityComparison:
		return v.Right
	case GreaterOrEqualComparison:
		return v.Right
	case GreaterComparison:
		return v.Right
	case LessThanOrEqualComparison:
		return v.Right
	case LessThanComparison:
		return v.Right
	default:
		panic("Unhandled node type in getRight")
	}
}

func setLeft(node, value Value) Value {
	switch v := node.(type) {
	case MulOperator:
		v.Left = value
		return v
	case DivOperator:
		v.Left = value
		return v
	case SubtractionOperator:
		v.Left = value
		return v
	case AdditionOperator:
		v.Left = value
		return v
	case NotEqualsComparison:
		v.Left = value
		return v
	case EqualityComparison:
		v.Left = value
		return v
	case GreaterOrEqualComparison:
		v.Left = value
		return v
	case GreaterComparison:
		v.Left = value
		return v
	case LessThanOrEqualComparison:
		v.Left = value
		return v
	case LessThanComparison:
		v.Left = value
		return v
	default:
		panic(fmt.Sprintf("Unhandled node type in setLeft", reflect.TypeOf(node)))
	}
}

func setRight(node, value Value) Value {
	switch v := node.(type) {
	case MulOperator:
		v.Right = value
		return v
	case DivOperator:
		v.Right = value
		return v
	case SubtractionOperator:
		v.Right = value
		return v
	case AdditionOperator:
		v.Right = value
		return v
	case NotEqualsComparison:
		v.Right = value
		return v
	case EqualityComparison:
		v.Right = value
		return v
	case GreaterOrEqualComparison:
		v.Right = value
		return v
	case GreaterComparison:
		v.Right = value
		return v
	case LessThanOrEqualComparison:
		v.Right = value
		return v
	case LessThanComparison:
		v.Right = value
		return v
	default:
		panic(fmt.Sprintf("Unhandled node type in setRight", reflect.TypeOf(node)))
	}
}

func invertPrecedence(node Value) Value {
	// This is mostly black magic, but it passes all the tests that currently
	// exist.
	// Suppose there's an expression 2 * 4 - 6 / 3.
	// As the expression is being calculated, it gets constructed as,
	// (2 * (4 - (6 / 3))) and this is called because the * and - precedence
	// is wrong. It needs to be converted to ((2*4) - (6 / 3)).
	//
	// On the other hand, suppose there's an expression 1 / 2 + 3 - 4.
	// When it gets here, it's been parsed left to right as (1 / (2 + (3 - 4)))
	// It should be ((1 / 2) + 3) - 4).
	//
	// I suspect there's other edge cases in longer/deeper expressions since
	// this code comes from test driven development and not logical
	// reasoning, but it covers enough it'll do for now.
	op1 := node
	op2 := getRight(node)
	op3 := getRight(op2)

	var isLit bool
	switch op3.(type) {
	case IntLiteral, BoolLiteral, Variable:
		isLit = true
	}

	switch op2.(type) {
	case IntLiteral, BoolLiteral, Variable:
		isLit = true
	}
	if isLit || operatorPrecedence(op3) > operatorPrecedence(op2) {
		// Handle the (2 * (4 - (6 / 3))) => ((2*4) - (6 / 3))
		lLeft := getLeft(op1)
		lRight := getLeft(op2)

		lOper := setRight(op1, lRight)
		lOper = setLeft(lOper, lLeft)

		mainOper := setRight(op2, op3)
		mainOper = setLeft(mainOper, lOper)
		return mainOper
	}
	// Handle the (1 / (2 + (3 - 4))) => ((1 / 2) + 3) - 4) case

	lLeft := getLeft(op1)
	lRight := getLeft(op2)

	v2 := getLeft(op3)
	//	v3 := getRight(op3)

	lOper := setRight(op1, lRight)
	lOper = setLeft(lOper, lLeft)

	op2 = setLeft(op2, lOper)
	op2 = setRight(op2, v2)

	op3 = setLeft(op3, op2)
	return op3
}

// Creates the AST node for an operator, taking precedence into account.
func createOperatorNode(op token.Token, left, right Value) Value {
	var v Value
	switch op {
	case token.Operator("+"):
		v = AdditionOperator{Left: left, Right: right}
	case token.Operator("-"):
		v = SubtractionOperator{Left: left, Right: right}
	case token.Operator("*"):
		v = MulOperator{Left: left, Right: right}
	case token.Operator("/"):
		v = DivOperator{Left: left, Right: right}
	case token.Operator("=="):
		v = EqualityComparison{Left: left, Right: right}
	case token.Operator("!="):
		v = NotEqualsComparison{Left: left, Right: right}
	case token.Operator(">"):
		v = GreaterComparison{Left: left, Right: right}
	case token.Operator(">="):
		v = GreaterOrEqualComparison{Left: left, Right: right}
	case token.Operator("<"):
		v = LessThanComparison{Left: left, Right: right}
	case token.Operator("<="):
		v = LessThanOrEqualComparison{Left: left, Right: right}
	default:
		panic("Unhandled operator type in createOperatorNode")
	}

	switch right.(type) {
	case IntLiteral, Variable, BoolLiteral:
		return v
	}

	if operatorPrecedence(v) <= operatorPrecedence(right) {
		return v
	}
	return invertPrecedence(v)
}
func consumeBoolValue(start int, tokens []token.Token, c Context) (int, Value, error) {
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			var partial Value
			if n, err := strconv.Atoi(t.String()); err == nil {
				partial = IntLiteral(n)
			} else if t.String() == "true" {
				partial = BoolLiteral(true)
			} else if t.String() == "false" {
				partial = BoolLiteral(false)
			} else if c.IsVariable(t.String()) {
				// If it's a defined variable, that takes
				// precedence as a symbol over a function
				// name.
				partial = Variable(t.String())
			} else if c.IsFunction(t.String()) {
				// if it's a function, call it.
				n, fc, err := consumeFuncCall(i, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				partial = fc
				i += n
			} else {
				// FIXME: Otherwise, it may still be a parameter.
				// Validate this.
				partial = Variable(t.String())
			}

			for isInfixOperator(i+1, tokens) {
				switch tokens[i+1] {
				case token.Operator("+"), token.Operator("-"),
					token.Operator("*"), token.Operator("/"),
					token.Operator("<"), token.Operator("<="),
					token.Operator("=="), token.Operator("!="),
					token.Operator(">"), token.Operator(">="):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}

					finalPos := i + 2 - start + n
					return finalPos, createOperatorNode(tokens[i+1], partial, right), nil
				case token.Operator("%"):
					// FIXME: The MOD operator should be normalized to work the same way
					// as the other operators
					n, right, err := consumeIntValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					i += n - 1
					partial = ModOperation{
						A: partial,
						B: right,
					}
				default:
					panic(fmt.Sprintf("Unhandled infix operator %v at %v", tokens[i].String(), i))
				}
			}
			return i + 1 - start, partial, nil
		case token.Whitespace:
			continue
		default:
			return 0, nil, fmt.Errorf("Invalid value: %v", tokens[i])
		}
	}
	return 0, nil, fmt.Errorf("No value")
}
