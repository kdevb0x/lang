package ast

import (
	"fmt"
)

type Callable interface {
	GetArgs() []VarWithType
}

// type Function should be the same as procedure, but
// until the statements are settled we're just have Funcedure
type FuncDecl struct {
	Name   string
	Args   []VarWithType
	Return []VarWithType

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
