package ast

import "fmt"

type ReturnStmt struct {
	Val Value
}

func (rs ReturnStmt) Node() Node {
	return rs
}

func (rs ReturnStmt) String() string {
	return fmt.Sprintf("ReturnStmt{ %v }", rs.Val)
}
