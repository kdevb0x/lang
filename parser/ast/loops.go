package ast

import (
	"fmt"
)

type WhileLoop struct {
	Condition BoolValue
	Body      BlockStmt
}

func (l WhileLoop) Node() Node {
	return l
}

func (l WhileLoop) String() string {
	return fmt.Sprintf("WhileLoop{\n\tCondition: %v\n\tBody: %v\n}", l.Condition, l.Body)
}

func (l WhileLoop) PrettyPrint(lvl int) string {
	panic("Not implemented")
}
