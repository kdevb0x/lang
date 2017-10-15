package ir

import (
	"fmt"
	"strings"

	"github.com/driusan/lang/parser/ast"
)

type variableLayout struct {
	values      map[ast.VarWithType]Register
	tempVars    int
	tempRegs    uint
	typeinfo    ast.TypeInformation
	funcargs    []ast.VarWithType
	rettypes    []ast.TypeInfo
	enumvalues  EnumMap
	callables   ast.Callables
	numLocals   uint
	maxFuncCall uint
}

func (c variableLayout) GetTypeInfo(t string) ast.TypeInfo {
	ti, ok := c.typeinfo[t]
	if !ok {
		panic("Could not get type info for " + string(t))
	}
	return ti
}
func (c variableLayout) GetReturnTypeInfo(v uint) ast.TypeInfo {
	return c.rettypes[v]
}

func (c *variableLayout) NextTempRegister() Register {
	r := TempValue(c.tempRegs)
	c.tempRegs++
	return r
}

// Reserves the next available register for varname
func (c *variableLayout) NextLocalRegister(varname ast.VarWithType) Register {
	if varname.Type() == "" {
		panic("No type for variable " + varname.Name + ".")
	}
	ti := c.typeinfo
	typ := varname.Type()
	firstType := strings.Fields(string(typ))[0]
	c.numLocals++

	if varname.Name == "" {
		c.tempVars++
		return LocalValue{uint(len(c.values) + c.tempVars - 1), ti[firstType]}
	}

	// If this variable is shadowing another variable, increase tempVars to
	// make sure the next calls increment the LocalVariable number and don't
	// reuse the same variable.
	_, postInc := c.values[varname]
	c.values[varname] = LocalValue{uint(len(c.values) + c.tempVars), ti[firstType]}
	if postInc {
		c.tempVars++
	}
	return c.values[varname]
}

// Reserves a register for a function parameter. This must be done for every
// parameter, before any LocalRegister calls are made.
func (c *variableLayout) FuncParamRegister(varname ast.VarWithType, i int) Register {
	c.tempVars--
	ti := c.typeinfo
	c.values[varname] = FuncArg{uint(i), ti[varname.Type()], varname.Reference}
	return c.values[varname]
}

// Sets a variable to refer to an existing register, without generating a new
// one.
func (c *variableLayout) SetLocalRegister(varname ast.VarWithType, val Register) {
	if _, ok := c.values[varname]; !ok {
		c.numLocals++
	}
	c.values[varname] = val
}

// Gets the register for an existing variable. Panics on invalid variables.
func (c variableLayout) Get(varname ast.VarWithType) Register {
	if varname.Name == "" {
		panic("Can not get empty varname")
	}
	val, ok := c.values[varname]
	if !ok {
		panic("Could not get variable named " + varname.Name)
	}
	return val
}

// Gets the register for an existing variable, and a bool denoting whether
// the variable exists or not.
func (c variableLayout) SafeGet(varname ast.VarWithType) (Register, bool) {
	v, ok := c.values[varname]
	return v, ok
}

func (c variableLayout) GetEnumIndex(v string) int {
	val, ok := c.enumvalues[v]
	if !ok {
		panic(fmt.Sprintf("Attempt to retrieve invalid enum option %v: ", v))
	}
	return val
}
