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
	var eln int = 0
	if tokens[start+cn+bn+1] == token.Keyword("else") {
		en, elseblock, err := consumeElseStmt(start+cn+bn+1, tokens, c)
		if err != nil {
			return 0, IfStmt{}, err
		}
		l.Else = elseblock
		eln = en
	}
	return cn + bn + eln + 1, l, nil
}

func consumeElseStmt(start int, tokens []token.Token, c *Context) (int, BlockStmt, error) {
	switch tokens[start+1] {
	case token.Keyword("if"):
		bn, block, err := consumeIfStmt(start+1, tokens, c)
		if err != nil {
			return 0, BlockStmt{}, err
		}
		return bn + 1, BlockStmt{[]Node{block}}, nil
	case token.Char("{"):
		bn, block, err := consumeBlock(start+1, tokens, c)
		if err != nil {
			return 0, BlockStmt{}, err
		}
		return bn + 1, block, nil
	default:
		return 0, BlockStmt{}, fmt.Errorf("Invalid else statement")
	}
}
