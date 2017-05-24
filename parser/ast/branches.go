package ast

import (
	"fmt"
	"github.com/driusan/lang/parser/token"
	"reflect"
)

func consumeWhileLoop(start int, tokens []token.Token, c Context) (int, Node, error) {
	l := WhileLoop{}

	if tokens[start] != token.Keyword("while") {
		return 0, nil, fmt.Errorf("Invalid while loop")
	}
	cn, cv, err := consumeCondition(start+1, tokens, c)
	if err != nil {
		return 0, nil, err
	}

	l.Condition = cv

	bn, block, err := consumeBlock(start+cn+1, tokens, c)
	if err != nil {
		return 0, nil, err
	}

	l.Body = block
	return cn + bn + 1, l, nil
}

func consumeIfStmt(start int, tokens []token.Token, c Context) (int, IfStmt, error) {
	l := IfStmt{}

	if tokens[start] != token.Keyword("if") {
		return 0, IfStmt{}, fmt.Errorf("Invalid if statement")
	}
	cn, cv, err := consumeCondition(start+1, tokens, c)
	if err != nil {
		return 0, IfStmt{}, err
	}

	l.Condition = cv

	bn, block, err := consumeBlock(start+cn+1, tokens, c)
	if err != nil {
		return 0, IfStmt{}, err
	}

	l.Body = block
	return cn + bn + 1, l, nil
}

func consumeCondition(start int, tokens []token.Token, c Context) (int, BoolValue, error) {
	n, cond, err := consumeBoolValue(start, tokens, c)
	if err != nil {
		return 0, nil, err
	}
	switch val := cond.(type) {
	case GreaterComparison:
		return n, val, nil
	case EqualityComparison:
		return n, val, nil
	case NotEqualsComparison:
		return n, val, nil
	case GreaterOrEqualComparison:
		return n, val, nil
	case LessThanComparison:
		return n, val, nil
	case LessThanOrEqualComparison:
		return n, val, nil
	default:
		return 0, nil, fmt.Errorf("Unsupported comparison %s", reflect.TypeOf(val))
	}
}
