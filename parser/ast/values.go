package ast

import (
	"fmt"
	"github.com/driusan/lang/parser/token"
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
				// FIXME: Only operators that are used by sample
				// programs are implemented.
				switch tokens[i+1] {
				case token.Operator("+"):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					return i + 2 - start + n, AdditionOperator{
						Left:  partial,
						Right: right,
					}, nil
				case token.Operator("-"):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					return i + 2 - start + n, SubtractionOperator{
						Left:  partial,
						Right: right,
					}, nil
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
				case token.Operator(">"):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					return i + 2 - start + n, GreaterComparison{
						Left:  partial,
						Right: right,
					}, nil
				case token.Operator(">="):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					return i + 2 - start + n, GreaterOrEqualComparison{
						Left:  partial,
						Right: right,
					}, nil
				case token.Operator("=="):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					return i + 2 - start + n, EqualityComparison{
						Left:  partial,
						Right: right,
					}, nil
				case token.Operator("!="):
					n, right, err := consumeValue(i+2, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					return i + 2 - start + n, NotEqualsComparison{
						Left:  partial,
						Right: right,
					}, nil
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
