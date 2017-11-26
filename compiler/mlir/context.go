package mlir

import (
	"fmt"
	//	"strings"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/ast"
)

type Context struct {
	Callables    ast.Callables
	RegisterData hlir.RegisterData
	curFunc      *Func
}

func NewContext(c ast.Callables, rd hlir.RegisterData) *Context {
	return &Context{
		Callables:    c,
		RegisterData: rd,
	}
}
func (c Context) GetFuncInfo(name string) ast.Callable {
	options := c.Callables[name]
	if len(options) == 0 {
		panic("Func " + name + " not found")
	} else if len(options) > 1 {
		panic("Multiple dispatch not implemented")
	}
	return options[0]

}

func (c Context) GetTypeInfo(r hlir.Register) ast.TypeInfo {
	ti, ok := c.RegisterData[r]
	if !ok {
		panic(fmt.Sprintf("Could not get metadata for %v", r))
	}
	if ti.TypeInfo.Size == 0 {
		ti.TypeInfo.Size = 8
	}
	return ti.TypeInfo
}
