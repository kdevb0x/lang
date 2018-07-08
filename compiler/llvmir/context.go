package llvmir

import (
	"fmt"

	"github.com/driusan/lang/parser/ast"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type FuncDef struct {
	*ir.Function
	ast.FuncDecl
}
type Context struct {
	// Global things
	Funcs          map[string]FuncDef
	StringLiterals map[string]*ir.Global

	enumValues EnumMap

	// Context sensitive things
	Variables map[ast.VarWithType]Register

	module     *ir.Module
	curfunc    *ir.Function
	curfuncdef ast.FuncDecl
	curblock   *ir.BasicBlock
	loopNum    uint
	ifNum      uint
	matchNum   uint
}

func NewContext(m *ir.Module) *Context {
	ctx := &Context{
		Funcs:          make(map[string]FuncDef),
		StringLiterals: make(map[string]*ir.Global),
		Variables:      make(map[ast.VarWithType]Register),
		enumValues:     make(EnumMap),
	}

	// FIXME: Load these from StdLib, don't hardcode them.
	write := m.NewFunction("Write", types.Void)
	write.NewParam("fd", types.I64)
	write.NewParam("buf", types.NewStruct(types.NewPointer(types.I8), types.I64))

	printstring := m.NewFunction("PrintString", types.Void)
	strty := types.NewStruct(types.NewPointer(types.I8), types.I64)
	printstring.NewParam("str", strty)

	printint := m.NewFunction("PrintInt", types.Void)
	printint.NewParam("n", types.I64)

	ctx.Funcs["PrintString"] = FuncDef{printstring, ast.FuncDecl{}}
	ctx.Funcs["PrintByteSlice"] = FuncDef{printstring, ast.FuncDecl{}}
	ctx.Funcs["Write"] = FuncDef{write, ast.FuncDecl{}}
	ctx.Funcs["PrintInt"] = FuncDef{printint, ast.FuncDecl{}}
	ctx.module = m

	return ctx
}

// Gets the EnumTypeDefn which created the token v for a EnumValue
// type constructor v
func (c *Context) GetEnumTypeDefn(v string) ast.EnumTypeDefn {
	val, ok := c.enumValues[v]
	if !ok {
		fmt.Printf("%v\n", c.enumValues)
		panic(fmt.Sprintf("Attempt to retrieve invalid enum option %v: ", v))
	}
	return val.Defn
}
func (c *Context) GetEnumIndex(v string) int {
	val, ok := c.enumValues[v]
	if !ok {
		fmt.Printf("%v\n", c.enumValues)
		panic(fmt.Sprintf("Attempt to retrieve invalid enum option %v: ", v))
	}
	return val.Index
}

func (ctx *Context) GetVariable(val ast.VarWithType) Register {
	if _, ok := val.Type().(ast.ArrayType); ok {
		val.Reference = true
	}
	if _, ok := ctx.Variables[hashableHack(val)]; !ok {
		fmt.Printf("\n\n%v", ctx.Variables)
		panic(fmt.Sprintf("Unknown variable %v", val))
	}
	return ctx.Variables[hashableHack(val)]
}

func (c *Context) SetVar(v ast.VarWithType, val Register) {
	if _, ok := v.Type().(ast.ArrayType); ok {
		v.Reference = true
	}
	c.Variables[hashableHack(v)] = val
}

func hashableHack(v ast.VarWithType) ast.VarWithType {
	switch t := v.Type().(type) {
	case ast.EnumTypeDefn:
		v.Typ = ast.TypeLiteral(t.Type().TypeName())
	case ast.SumType:
		v.Typ = ast.TypeLiteral(t.TypeName())
	case ast.TupleType:
		v.Typ = ast.TypeLiteral(t.TypeName())
	case ast.UserType:
		v.Typ = ast.TypeLiteral(t.TypeName())
	}
	return v
}

func (c *Context) cloneVars() map[ast.VarWithType]Register {
	rv := make(map[ast.VarWithType]Register)
	for k, v := range c.Variables {
		rv[k] = v
	}
	return rv
}
