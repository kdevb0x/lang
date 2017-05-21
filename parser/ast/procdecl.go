package ast

import (
	"fmt"
	"github.com/driusan/lang/parser/token"
)

// type Function should be the same as procedure, but
// until the statements are settled we're just have Procedure
type ProcDecl struct {
	Name   string
	Args   []VarWithType
	Return []VarWithType

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

func (pd *ProcDecl) PopulateName(t token.Token) error {
	switch v := t.(type) {
	case token.Whitespace:
		// an error would be fatal, so just return nil
		// to try again on the next token
		return nil
	case token.Unknown:
		pd.Name = v.String()
		return nil
	default:
		return fmt.Errorf("Invalid proc name: %v", t.String())
	}
}

func (pd ProcDecl) GetArgs() []VarWithType {
	return pd.Args
}
