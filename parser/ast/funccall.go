package ast

import "fmt"

type FuncCall struct {
	Name string
	Args []Value
}

func (f FuncCall) Node() Node {
	return f
}

func (f FuncCall) Value() interface{} {
	return nil
}

func (f FuncCall) String() string {
	return fmt.Sprintf("FuncCall{Name: %v Args: %v}", f.Name, f.Args)
}
