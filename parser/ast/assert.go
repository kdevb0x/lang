package ast

import (
	"fmt"

	"github.com/driusan/lang/parser/token"
)

type Assertion struct {
	Message   string
	Predicate Value
}

func (a Assertion) Node() Node {
	return a
}

func (a Assertion) String() string {
	return fmt.Sprintf("Assertion{ %v }", a.Predicate)
}

func (a Assertion) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

func consumeAssertStmt(start int, tokens []token.Token, c *Context) (int, Assertion, error) {
	a := Assertion{}
	if tokens[start] != token.Keyword("assert") {
		return 0, Assertion{}, fmt.Errorf("Invalid assertion statement")
	}
	if tokens[start+1] != token.Char("(") {
		return 0, Assertion{}, fmt.Errorf("Invalid assertion statement: Missing opening bracket")
	}
	vn, v, err := consumeValue(start+2, tokens, c)
	if err != nil {
		return 0, Assertion{}, err
	}
	if v.Type() != "bool" {
		return 0, Assertion{}, fmt.Errorf("Assertion predicate must be a bool")
	}
	a.Predicate = v
	switch tokens[start+vn+2] {
	case token.Char(")"):
		return vn + 3, a, nil
	case token.Char(","):
		mn, m, err := consumeValue(start+vn+3, tokens, c)
		if err != nil {
			return 0, Assertion{}, err
		}
		ms, ok := m.(StringLiteral)
		if !ok {
			return 0, Assertion{}, fmt.Errorf("Invalid assertion statement: second argument must be a string literal")
		}
		a.Message = string(ms)
		if tokens[start+vn+mn+3] != token.Char(")") {
			return 0, Assertion{}, fmt.Errorf("Invalid assertion statement: missing closing bracket")
		}
		return vn + mn + 4, a, nil

	default:
		return 0, Assertion{}, fmt.Errorf("Invalid assertion statement: expecting ',' or ')'")
	}

}
