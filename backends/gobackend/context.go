package gobackend

import (
	"fmt"

	"github.com/driusan/lang/parser/ast"
)

type Context struct {
	// Maps a VarWithType to a unique name to be used in the
	// Go source code.
	variables map[ast.VarWithType]string

	// Imports which are needed in order to compile the Go
	// source code.
	importsNeeded map[string]bool

	// Force GetVarName to create a new variable name
	newVar bool
}

func NewContext() *Context {
	return &Context{
		variables:     make(map[ast.VarWithType]string),
		importsNeeded: make(map[string]bool),
	}
}

func (c *Context) GetVarName(x ast.VarWithType) string {
	if v, ok := c.variables[x]; !c.newVar && ok {
		return v
	}
	hasname := func(x string) bool {
		for _, val := range c.variables {
			if val == x {
				return true
			}
		}
		return false
	}
	if !hasname(x.Name.String()) {
		c.variables[x] = x.Name.String()
		return x.Name.String()
	}

	i := 0
	for hasname(fmt.Sprintf("%s%d", x.Name.String(), i)) {
		i++
	}
	c.variables[x] = fmt.Sprintf("%s%d", x.Name.String(), i)
	return c.variables[x]
}
