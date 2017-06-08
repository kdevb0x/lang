package ast

import (
	"github.com/driusan/lang/parser/token"

	"fmt"
)

type MatchCase struct {
	Variable Value
	Body     BlockStmt
}

func (i MatchCase) Node() Node {
	return i
}

type MatchStmt struct {
	Condition Value
	Cases     []MatchCase
}

func (i MatchStmt) String() string {
	return fmt.Sprintf("MatchStmt{Condition: %v,\n\tBody: %v}", i.Condition, i.Cases)
}

func (i MatchStmt) Node() Node {
	return i
}

func consumeMatchStmt(start int, tokens []token.Token, c *Context) (int, MatchStmt, error) {
	l := MatchStmt{}

	if tokens[start] != token.Keyword("match") {
		return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
	}
	cn, cv, err := consumeValue(start+1, tokens, c)
	if err != nil {
		return 0, MatchStmt{}, err
	}
	l.Condition = cv

	if tokens[start+cn+1] != token.Char("{") {
		return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
	}

	for i := start + cn + 2; i < len(tokens); {
		c2 := c.Clone()
		n, cs, err := consumeCase(i, tokens, &c2)
		if err != nil {
			return 0, MatchStmt{}, err
		}
		l.Cases = append(l.Cases, cs)
		i += n
		if tokens[i+1] == token.Char("}") {
			return i - start, l, nil
		}
	}
	return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
}

func consumeCase(start int, tokens []token.Token, c *Context) (int, MatchCase, error) {
	l := MatchCase{}

	if tokens[start] != token.Keyword("case") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement. Unexpected '%v' at %d", tokens[start], start)
	}
	n, v, err := consumeValue(start+1, tokens, c)
	if err != nil {
		return 0, MatchCase{}, err
	}
	l.Variable = v
	if tokens[start+n+1] != token.Char(":") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement at token %v. Expected ':', not '%v'", start, tokens[start+n+1])
	}
	for i := start + n + 2; i < len(tokens); {
		if tokens[i] == token.Keyword("case") || tokens[i] == token.Char("}") {
			return i - start, l, nil
		}
		n, stmt, err := consumeStmt(i, tokens, c)
		if err != nil {
			return 0, MatchCase{}, err
		}
		l.Body.Stmts = append(l.Body.Stmts, stmt)
		i += n
	}
	return 0, MatchCase{}, fmt.Errorf("Unterminated case statement")
}
