package hlir

import (
	"fmt"

	"github.com/driusan/lang/parser/ast"
)

type RegisterData map[Register]RegisterInfo

type RegisterInfo struct {
	Name     string
	TypeInfo ast.TypeInfo
	Variable ast.VarWithType

	// Only valid for slices
	SliceSize uint
	Creator   ast.VarWithType
	//Creator ast.Value

	Type ast.Type
}

type variableLayout struct {
	values       map[ast.VarWithType]Register
	sliceBase    map[ast.VarWithType]Register
	tempVars     int
	tempRegs     uint
	typeinfo     ast.TypeInformation
	funcargs     []ast.VarWithType
	rettypes     []ast.TypeInfo
	retsumtypes  map[int]ast.Type
	enumvalues   EnumMap
	callables    ast.Callables
	numLocals    uint
	registerInfo RegisterData
	loopCond     bool
	loop         *LOOP
}

func (c variableLayout) GetTypeInfo(t string) ast.TypeInfo {
	if t == "string" {
		return ast.TypeInfo{
			Size:   16,
			Signed: false,
		}
	}
	ti, ok := c.typeinfo[t]
	if !ok {
		panic("Could not get type info for " + string(t))
	}
	return ti
}

func (c *variableLayout) NextTempRegister() Register {
	r := TempValue(c.tempRegs)
	c.tempRegs++
	return r
}

// Reserves the next available register for varname
func (c *variableLayout) NextLocalRegister(varname ast.VarWithType) Register {
	if varname.Type() == nil || varname.Type().TypeName() == "" {
		panic("No type for variable " + varname.Name + ".")
	}

	if varname.Name == "" {
		panic("No name for variable")
	}

	c.numLocals++

	switch t := varname.Type().(type) {
	case ast.SumType, ast.EnumTypeDefn, ast.TupleType, ast.UserType:
		// Hack, since SumType is unhashable and can't
		// be used as a key for c.values
		varname.Typ = ast.TypeLiteral(t.TypeName())
	}
	// If this variable is shadowing another variable, increase tempVars to
	// make sure the next calls increment the LocalVariable number and don't
	// reuse the same variable.
	_, postInc := c.values[varname]
	lv := LocalValue(uint(len(c.values) + c.tempVars))
	c.values[varname] = lv
	if postInc {
		c.tempVars++
	}
	c.registerInfo[lv] = RegisterInfo{
		string(varname.Name),
		c.typeinfo[varname.Type().TypeName()],
		varname,
		0,
		ast.VarWithType{},
		varname.Typ,
	}
	return lv
}

// Reserves a register for a function parameter. This must be done for every
// parameter, before any LocalRegister calls are made.
func (c *variableLayout) FuncParamRegister(varname ast.VarWithType, i int) Register {
	c.tempVars--
	fa := FuncArg{uint(i), varname.Reference}
	switch t := varname.Type().(type) {
	case ast.SumType:
		// Hack, since SumType is unhashable and can't
		// be used as a key for c.values
		varname.Typ = ast.TypeLiteral(t.TypeName())
	}
	c.values[hashableHack(varname)] = fa
	c.registerInfo[fa] = RegisterInfo{
		string(varname.Name),
		c.typeinfo[varname.Type().TypeName()],
		varname,
		0,
		ast.VarWithType{},
		varname.Typ,
	}
	return c.values[hashableHack(varname)]
}

// Sets a variable to refer to an existing register, without generating a new
// one.
func (c *variableLayout) SetLocalRegister(varname ast.VarWithType, val Register) {
	c.values[hashableHack(varname)] = val
}

// Gets the register for an existing variable. Panics on invalid variables.
func (c variableLayout) Get(varname ast.VarWithType) Register {
	if varname.Name == "" {
		panic("Can not get empty varname")
	}
	val, ok := c.values[hashableHack(varname)]
	if !ok {
		fmt.Printf("%v", c.values)
		panic("Could not get variable named " + varname.Name)
	}
	return val
}

// Gets the register for an existing variable, and a bool denoting whether
// the variable exists or not.
func (c variableLayout) SafeGet(varname ast.VarWithType) (Register, bool) {
	v, ok := c.values[hashableHack(varname)]
	return v, ok
}

func (c variableLayout) GetEnumIndex(v string) int {
	val, ok := c.enumvalues[v]
	if !ok {
		fmt.Printf("%v\n", c.enumvalues)
		panic(fmt.Sprintf("Attempt to retrieve invalid enum option %v: ", v))
	}
	return val
}

func (c variableLayout) CloneValues() map[ast.VarWithType]Register {
	v := make(map[ast.VarWithType]Register)
	for k, vl := range c.values {
		v[k] = vl
	}
	return v
}

func hashableHack(varname ast.VarWithType) ast.VarWithType {
	switch t := varname.Type().(type) {
	case ast.SumType, ast.EnumTypeDefn, ast.TupleType, ast.UserType:
		// Hack, since SumType is unhashable and can't
		// be used as a key for c.values
		varname.Typ = ast.TypeLiteral(t.TypeName())
	}
	return varname
}
