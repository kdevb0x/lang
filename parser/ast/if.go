package ast

import (
	"fmt"

	"github.com/driusan/lang/parser/token"
)

type IfStmt struct {
	Condition BoolValue
	Body      BlockStmt
	Else      BlockStmt
}

func (i IfStmt) String() string {
	return fmt.Sprintf("IfStmt{Condition: %v,\n\tBody: %v,\n\tElse: %v}", i.Condition, i.Body, i.Else)
}

func (i IfStmt) Node() Node {
	return i
}

func consumeIfStmt(start int, tokens []token.Token, c *Context) (int, IfStmt, error) {
	l := IfStmt{}

	if tokens[start] != token.Keyword("if") {
		return 0, IfStmt{}, fmt.Errorf("Invalid if statement")
	}
	cn, cv, err := consumeCondition(start+1, tokens, c)
	if err != nil {
		return 0, IfStmt{}, err
	}

	l.Condition = cv

	c2 := c.Clone()
	bn, block, err := consumeBlock(start+cn+1, tokens, &c2)
	if err != nil {
		return 0, IfStmt{}, err
	}

	l.Body = block
	return cn + bn + 1, l, nil
}
