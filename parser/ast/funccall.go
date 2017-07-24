package ast

import "fmt"

type FuncCall struct {
	Name     string
	UserArgs []Value
	Returns  Tuple
}

func (f FuncCall) Node() Node {
	return f
}

func (f FuncCall) Value() interface{} {
	return nil
}

func (f FuncCall) String() string {
	return fmt.Sprintf("FuncCall{Name: %v Args: %v}", f.Name, f.UserArgs)
}

// FIXME: This needs to be updated to work with multiple return functions
func (f FuncCall) Type() string {
	if len(f.Returns) == 0 {
		return "(none)"
	}
	return f.Returns[0].Type()
}
