package ast

import (
	"fmt"
	"github.com/driusan/lang/parser/token"
	"reflect"
	"strconv"
)

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
	case token.Char:
		return t == "[" //|| t == "("
	}
	return false
}

func operatorPrecedence(op Value) int {
	switch op.(type) {
	case Brackets:
		return 99
	case ArrayValue, FuncCall:
		return 6
	case MulOperator, DivOperator, ModOperator:
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
	case ModOperator:
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
	case ModOperator:
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
		panic(fmt.Sprintf("Unhandled node type in getRight: %v", reflect.TypeOf(node)))
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
	case ModOperator:
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
		panic(fmt.Sprintf("Unhandled node type in setLeft: %v", reflect.TypeOf(node)))
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
	case ModOperator:
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
		panic(fmt.Sprintf("Unhandled node type in setRight: %v", reflect.TypeOf(node)))
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
	case IntLiteral, BoolLiteral, VarWithType, Brackets:
		isLit = true
	}

	switch op2.(type) {
	case IntLiteral, BoolLiteral, VarWithType, Brackets:
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
	case token.Operator("%"):
		v = ModOperator{Left: left, Right: right}
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
	case IntLiteral, VarWithType, BoolLiteral:
		return v
	}

	if operatorPrecedence(v) <= operatorPrecedence(right) {
		return v
	}
	return invertPrecedence(v)
}

func consumeValue(start int, tokens []token.Token, c *Context) (int, Value, error) {
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
				partial = c.Variables[t.String()]
			} else if eo := c.EnumeratedOption(t.String()); eo != nil {
				ev := EnumValue{Constructor: *eo}
				i += 1
				for j := 0; j < len(eo.Parameters); j++ {
					n, param, err := consumeValue(i, tokens, c)
					if err != nil {
						return 0, nil, err
					}
					ev.Parameters = append(ev.Parameters, param)
					i += n
				}
				i -= 1
				partial = ev
			} else if c.IsFunction(t.String()) {
				// if it's a function, call it.
				n, fc, err := consumeFuncCall(i, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				partial = fc
				i += n
			} else if eo := c.EnumeratedOption(t.String()); eo != nil {
				return i + 1 - start, *eo, nil
			} else {
				return 0, nil, fmt.Errorf(`Use of undefined variable "%v".`, t)
			}

			for isInfixOperator(i+1, tokens) && tokens[i+1] != token.Operator("=") {
				n, v, err := consumeInfix(i+1, tokens, c, partial)
				if err != nil {
					return 0, nil, err
				}
				partial = v
				i += n + 1
			}
			return i + 1 - start, partial, nil
		case token.Whitespace:
			continue
		case token.Char:
			switch tokens[i] {
			case token.Char(`"`):
				if tokens[i+2] != token.Char(`"`) {
					return 0, nil, fmt.Errorf("Invalid string at %v", tokens[i])
				}
				return 3, StringLiteral(tokens[i+1].String()), nil
			case token.Char(`{`):
				tn, v, err := consumeCommaSeparatedValues(i+1, tokens, c)
				al := ArrayLiteral(v)
				at := al.Type()
				if _, ok := c.Types[at]; !ok {
					typdef := ArrayType{
						Base: TypeLiteral(v[0].Type()),
						Size: IntLiteral(len(v)),
					}

					c.Types[at] = TypeDefn{
						Name:         TypeLiteral(at),
						ConcreteType: typdef,
					}
				}
				return tn + 2, al, err
			case token.Char(`[`):
				return 0, nil, fmt.Errorf("Indexing not yet implemented")
			case token.Char(`(`):
				tn, partial, err := consumeBrackets(i, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				i += tn

				for isInfixOperator(i+1, tokens) && tokens[i+1] != token.Operator("=") {
					n, v, err := consumeInfix(i+1, tokens, c, partial)
					if err != nil {
						return 0, nil, err
					}
					partial = v
					i += n + 1
				}
				return i + 1 - start, partial, nil
			default:
				return 0, nil, fmt.Errorf("Invalid character at %v", tokens[i])
			}
		case token.Operator:
			if t == token.Operator("-") {
				n, inverse, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
				}

				switch i := inverse.(type) {
				case IntLiteral:
					return n + 1, i * -1, nil
				default:
					return n + 1, MulOperator{IntLiteral(-1), inverse}, nil
				}
			}
			return 0, nil, fmt.Errorf("Invalid operator while expecting value: %v", tokens[i])
		case token.Keyword:
			if t == "let" {
				return consumeLetStmt(i, tokens, c)
			}
			return 0, nil, fmt.Errorf("Only let statements may be used inside of a value context.")
		default:
			return 0, nil, fmt.Errorf("Invalid value: %v (%v)", tokens[i], reflect.TypeOf(tokens[i]))
		}
	}
	return 0, nil, fmt.Errorf("No value")
}

func consumeCommaSeparatedValues(start int, tokens []token.Token, c *Context) (int, []Value, error) {
	var v []Value
	for i := start; i < len(tokens); i++ {
		switch tokens[i] {
		case token.Char("}"):
			return i - start, v, nil
		case token.Char(","):
			n, val, err := consumeValue(i+1, tokens, c)
			if err != nil {
				return 0, nil, err
			}
			v = append(v, val)
			i += n
		default:
			if i == start {
				n, val, err := consumeValue(i, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				v = append(v, val)
				i += n - 1
			} else {
				return 0, nil, fmt.Errorf("Unexpected token %v", tokens[i])
			}
		}

	}
	return len(tokens) - start, v, nil
}

func consumeInfix(start int, tokens []token.Token, c *Context, left Value) (int, Value, error) {
	switch tokens[start] {
	case token.Operator("+"), token.Operator("-"),
		token.Operator("*"), token.Operator("/"),
		token.Operator("%"),
		token.Operator("<"), token.Operator("<="),
		token.Operator("=="), token.Operator("!="),
		token.Operator(">"), token.Operator(">="):
		n, right, err := consumeValue(start+1, tokens, c)
		if err != nil {
			return 0, nil, err
		}

		return n, createOperatorNode(tokens[start], left, right), nil
	case token.Char("["):
		n, index, err := consumeValue(start+1, tokens, c)
		if err != nil {
			return 0, nil, err
		}
		if tokens[start+1+n] != token.Char("]") {
			return 0, nil, fmt.Errorf("Invalid index")
		}
		base, ok := left.(VarWithType)
		if !ok {
			return 0, nil, fmt.Errorf("Can only index on variables")
		}
		return n + 1, ArrayValue{base, index}, nil
	case token.Operator("="):
		return 0, left, nil
	default:
		panic(fmt.Sprintf("Unhandled infix operator %v at %v", tokens[start].String(), start))
	}
}

func consumeBrackets(start int, tokens []token.Token, c *Context) (int, Value, error) {
	switch tokens[start] {
	case token.Char("("):
		n, val, err := consumeValue(start+1, tokens, c)
		if err != nil {
			return 0, nil, err
		}
		if tokens[start+1+n] != token.Char(")") {
			return 0, nil, fmt.Errorf("Unbalanced parenthesis at %d", start)
		}
		return n + 1, Brackets{Val: val}, nil
	default:
		return 0, nil, fmt.Errorf("Brackets must start with a bracket")
	}
}
