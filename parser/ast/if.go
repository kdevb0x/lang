package ast

import (
	"fmt"
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
