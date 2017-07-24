package ast

import (
	"fmt"
)

type Callable interface {
	Type
	GetArgs() []VarWithType
	ReturnTuple() Tuple
}

type Tuple []VarWithType

func (t Tuple) Type() string {
	if len(t) == 0 {
		return "(none)"
	}
	// FIXME: This should take into account all types.
	return t[0].Type()
}

// type Function should be the same as procedure, but
// until the statements are settled we're just have Funcedure
type FuncDecl struct {
	Name   string
	Args   Tuple
	Return Tuple

	Body BlockStmt
}

func (pd FuncDecl) Node() Node {
	return pd
}

func (pd FuncDecl) GetArgs() []VarWithType {
	return pd.Args
}

func (fd FuncDecl) String() string {
	return fmt.Sprintf("FuncDecl{\n\tName: %v,\n\tArgs: %v,\n\tReturn: %v,\n\tBody: %v}", fd.Name, fd.Args, fd.Return, fd.Body)
}

func (fd FuncDecl) Type() string {
	return fd.Return.Type()
}

func (fd FuncDecl) ReturnTuple() Tuple {
	return fd.Return
}
