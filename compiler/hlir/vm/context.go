package vm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/ast"
)

type Pointer struct {
	r   hlir.Register
	ctx *Context
}

type Context struct {
	Funcs        map[string]hlir.Func
	Callables    ast.Callables
	RegisterData map[string]hlir.RegisterData

	// FIXME: This should probably be an io.ReadWriter instead
	stdout, stderr *strings.Builder

	localValues        map[hlir.LocalValue]interface{}
	funcRetVal         map[hlir.FuncRetVal]interface{}
	lastFuncCallRetVal map[hlir.LastFuncCallRetVal]interface{}
	tempValue          map[hlir.TempValue]interface{}
	funcArg            map[hlir.FuncArg]interface{}
	pointers           map[hlir.Pointer]Pointer
}

func (c *Context) String() string {
	return fmt.Sprintf("\tLocalValues: %v\n\tFuncRetVals: %v\n\tLastFuncCallRetVals: %v\n\tTempValues: %v\n\tFuncArgs: %v\n\tPointers: %v\n", c.localValues, c.funcRetVal, c.lastFuncCallRetVal, c.tempValue, c.funcArg, c.pointers)
}
func NewContext() *Context {
	c := &Context{}
	c.stdout = &strings.Builder{}
	c.stderr = &strings.Builder{}
	c.RegisterData = make(map[string]hlir.RegisterData)

	// One map for each register type's value in the current context
	c.localValues = make(map[hlir.LocalValue]interface{})
	c.funcRetVal = make(map[hlir.FuncRetVal]interface{})
	c.lastFuncCallRetVal = make(map[hlir.LastFuncCallRetVal]interface{})
	c.tempValue = make(map[hlir.TempValue]interface{})
	c.funcArg = make(map[hlir.FuncArg]interface{})
	c.pointers = make(map[hlir.Pointer]Pointer)
	return c
}

func (c *Context) Clone() *Context {
	return &Context{
		Funcs:        c.Funcs,
		Callables:    c.Callables,
		RegisterData: c.RegisterData,

		stdout: c.stdout,
		stderr: c.stderr,

		localValues:        c.localValues,
		funcRetVal:         c.funcRetVal,
		lastFuncCallRetVal: c.lastFuncCallRetVal,
		tempValue:          c.tempValue,
		funcArg:            c.funcArg,
		pointers:           c.pointers,
	}
}

func (c *Context) SetRegister(r hlir.Register, val interface{}) error {
	switch reg := r.(type) {
	case hlir.LocalValue:
		c.localValues[reg] = val
	case hlir.FuncRetVal:
		c.funcRetVal[reg] = val
	case hlir.TempValue:
		c.tempValue[reg] = val
	case hlir.FuncArg:
		c.funcArg[reg] = val
	case hlir.Offset:
		off, _ := resolveOffset(reg, c)
		c.SetRegister(off, val)
	default:
		panic(fmt.Sprintf("Unhandled register type: %v", reflect.TypeOf(r).Name()))
	}
	return nil
}

func (c *Context) writeStderr(msg string) {
	fmt.Fprintf(c.stderr, "%s", msg)
}
