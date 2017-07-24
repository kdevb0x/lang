package ast

import (
	"fmt"
)

// type Function should be the same as procedure, but
// until the statements are settled we're just have Procedure
type ProcDecl struct {
	Name   string
	Args   Tuple
	Return Tuple

	Body BlockStmt
}

func (pd ProcDecl) String() string {
	return fmt.Sprintf(`ProcDecl(Name: %v,
	Args: %v
	Return: %v

	Body: %v)`, pd.Name, pd.Args, pd.Return, pd.Body)
}

func (pd ProcDecl) Node() Node {
	return pd
}

func (pd ProcDecl) GetArgs() []VarWithType {
	return pd.Args
}

func (pd ProcDecl) Type() string {
	return pd.Return.Type()
}

func (pd ProcDecl) ReturnTuple() Tuple {
	return pd.Return
}
