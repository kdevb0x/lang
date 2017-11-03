package wasm

import (
	"fmt"
	"strings"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/ast"
)

type importName struct{ namespace, name string }

type Context struct {
	Imports               []Import
	Data                  []DataSection
	Functions             []Func
	stringliterals        map[string]int
	imports               map[importName]int
	callables             ast.Callables
	memTop                int
	curFuncNumArgs        uint
	curFuncNumReturns     uint
	size64                bool
	registerData          hlir.RegisterData
	typeinfo              ast.TypeInformation
	needsGlobal           bool
	curFuncRetNeedsMem    bool
	lastCallRetNeedsMem   bool
	curFuncMemVariables   map[hlir.Register]uint
	curFuncLocalVariables map[hlir.Register]uint
	curFuncMaxMem         uint
}

func encodeStrInt(s string) string {
	x := uint64(len(s))
	// `\n` is encoded as the string literal `\n`. We need to take it into account..
	x -= uint64(strings.Count(s, `\n`))
	return fmt.Sprintf(`\%0.2x\%0.2x\%0.2x\%0.2x\%0.2x\%0.2x\%0.2x\%0.2x`,
		x&0xff,
		(x&0xff00)>>8,
		(x&0xff0000)>>16,
		(x&0xff000000)>>24,
		(x&0xff00000000)>>32,
		(x&0xff0000000000)>>40,
		(x&0xff000000000000)>>48,
		(x&0xff00000000000000)>>56,
	)
}

func (c *Context) GetLiteral(s string) int {
	if i, ok := c.stringliterals[s]; ok {
		return i
	}
	sAddr := c.memTop
	c.memTop += 8
	c.memTop += len(s)
	c.stringliterals[s] = sAddr

	c.Data = append(c.Data,
		DataSection{
			int32(sAddr),
			encodeStrInt(s),
		},
		DataSection{
			int32(sAddr) + 8, // content
			s,
		})

	if align := c.memTop % 8; align != 0 {
		c.memTop += (8 - align)
	}
	return sAddr
}

func (c Context) LocalIndex(l hlir.LocalValue) uint {
	idx, ok := c.curFuncLocalVariables[l]
	if !ok {
		panic(fmt.Sprintf("Unknown LocalValue: %v", l))
	}
	return idx
}
func (c *Context) AddImport(namespace, fname string, signature Signature) int {
	if i, ok := c.imports[importName{namespace, fname}]; ok {
		return i
	}
	c.imports[importName{namespace, fname}] = len(c.imports)

	c.Imports = append(c.Imports, Import{namespace, fname, Func{
		Name:      fname,
		Signature: signature,
	},
	})
	return len(c.Imports) - 1
}

func (c *Context) GetNumArgs(fname string) uint {
	callable, ok := c.callables[fname]
	if !ok {
		panic("Could not find function")
	}
	if len(callable) != 1 {
		panic("Multiple dispatch not implemented")
	}

	return uint(len(callable[0].GetArgs()))
}

func (c *Context) GetNumReturns(fname string) uint {
	callable, ok := c.callables[fname]
	if !ok {
		panic("Could not find function")
	}
	if len(callable) != 1 {
		panic("Multiple dispatch not implemented")
	}

	return uint(len(callable[0].ReturnTuple()))
}

func (c *Context) GetSignature(fname string) Signature {
	callable, ok := c.callables[fname]
	if !ok {
		return nil
	}
	var ret Signature
	if len(callable) != 1 {
		panic("Multiple dispatch not implemented")
	}

	for _, v := range callable[0].GetArgs() {
		switch v.Typ.(type) {
		case ast.SliceType:
			ret = append(ret, Variable{i32, Param, ""}, Variable{i32, Param, ""})
		default:
			words := strings.Fields(string(v.Type()))
			for wordi, word := range words {
				typeinfo, ok := c.typeinfo[word]
				if !ok {
					panic("Could not get type info for argument:" + word)
				}

				var t VarType
				// Determine the size context based on the expected argument type
				// for IntLiterals
				if typeinfo.Size > 4 {
					t = i64
				} else {
					t = i32
				}
				name := ""
				if wordi == 0 {
					name = string(v.Name)
				}
				ret = append(ret, Variable{
					// FIXME: This should be the real type
					t,
					Param,
					name,
				})
			}
		}

	}

	for arg, v := range callable[0].ReturnTuple() {
		words := strings.Fields(string(v.Type()))
		if len(words) > 1 || arg > 1 {
			c.needsGlobal = true
			c.curFuncRetNeedsMem = true
		}
		info, ok := c.registerData[v]
		if !ok {
			panic("Could not get return type info")
		}
		var t VarType
		if info.TypeInfo.Size > 4 {
			t = i64
		} else {
			t = i32
		}
		ret = append(ret, Variable{
			t,
			Result,
			string(v.Name),
		})
	}
	return ret
}

func (c *Context) addMemoryVar(v hlir.LocalValue) {
	c.curFuncMemVariables[v] = c.curFuncMaxMem
	md, ok := c.registerData[v]
	if !ok {
		panic("Unknown register")
	}
	if md.TypeInfo.Size != 0 {
		c.curFuncMaxMem += uint(md.TypeInfo.Size)
	} else {
		c.curFuncMaxMem += 4
	}
	c.needsGlobal = true
}

func NewContext(c ast.Callables) Context {
	return Context{
		make([]Import, 0),
		make([]DataSection, 0),
		make([]Func, 0),
		make(map[string]int),
		make(map[importName]int),
		c,
		0,
		0,
		0,
		false,
		nil,
		nil,
		false,
		false,
		false,
		make(map[hlir.Register]uint),
		make(map[hlir.Register]uint),
		0,
	}
}
