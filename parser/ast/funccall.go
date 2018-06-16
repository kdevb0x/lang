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
func (f FuncCall) Type() Type {
	if len(f.Returns) == 1 {
		return f.Returns[0].Type()
	}
	return f.Returns.Type()
}

func (f FuncCall) PrettyPrint(lvl int) string {
	ret := fmt.Sprintf("%v%v(", nTabs(lvl), f.Name)
	for i, v := range f.UserArgs {
		if i == 0 {
			ret += v.PrettyPrint(0)
		} else {
			ret += ", " + v.PrettyPrint(0)
		}

	}
	return ret + ")"
}
