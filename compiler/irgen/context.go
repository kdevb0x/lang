package irgen

import (
	"fmt"

	"github.com/driusan/lang/compiler/ir"
	"github.com/driusan/lang/parser/ast"
)

type variableLayout struct {
	values     map[ast.VarWithType]ir.Register
	tempVars   int
	typeinfo   ast.TypeInformation
	rettypes   []ast.TypeInfo
	enumvalues EnumMap
}

func (c variableLayout) GetTypeInfo(t ast.Type) ast.TypeInfo {
	return c.typeinfo[t]
}
func (c variableLayout) GetReturnTypeInfo(v uint) ast.TypeInfo {
	return c.rettypes[v]
}

// Reserves the next available register for varname
func (c *variableLayout) NextLocalRegister(varname ast.VarWithType) ir.Register {
	if varname.Type() == "" {
		panic("No type for variable.")
	}
	ti := c.typeinfo
	if varname.Name == "" {
		c.tempVars++
		return ir.LocalValue{uint(len(c.values) + c.tempVars - 1), ti[varname.Type()]}
	}

	c.values[varname] = ir.LocalValue{uint(len(c.values) + c.tempVars), ti[varname.Type()]}
	return c.values[varname]
}

// Reserves a register for a function parameter. This must be done for every
// parameter, before any LocalRegister calls are made.
func (c *variableLayout) FuncParamRegister(varname ast.VarWithType, i int) ir.Register {
	c.tempVars--
	ti := c.typeinfo
	c.values[varname] = ir.FuncArg{uint(i), ti[varname.Type()]}
	return c.values[varname]
}

// Sets a variable to refer to an existing register, without generating a new
// one.
func (c *variableLayout) SetLocalRegister(varname ast.VarWithType, val ir.Register) {
	c.values[varname] = val
}

// Gets the register for an existing variable. Panics on invalid variables.
func (c variableLayout) Get(varname ast.VarWithType) ir.Register {
	if varname.Name == "" {
		panic("Can not get empty varname")
	}
	return c.values[varname]
}

// Gets the register for an existing variable, and a bool denoting whether
// the variable exists or not.
func (c variableLayout) SafeGet(varname ast.VarWithType) (ir.Register, bool) {
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
