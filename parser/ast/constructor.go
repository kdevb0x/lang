package ast

import (
	"fmt"
	//	"os"
	"reflect"
	"strings"

	"github.com/driusan/lang/parser/token"
)

type Context struct {
	Variables   map[string]VarWithType
	Mutables    map[string]VarWithType
	Functions   map[string]Callable
	Types       map[string]TypeDefn
	PureContext bool // true if inside a pure function.
	EnumOptions map[string]EnumOption
}

func NewContext() Context {
	return Context{
		Variables: make(map[string]VarWithType),
		Functions: map[string]Callable{
			// FIXME: These should be replaced by a
			// multiple dispatch Print function and/or
			// moved to a standard library, not build in.
			"PrintString": ProcDecl{
				Name: "PrintString",
				Args: []VarWithType{
					{"", TypeLiteral("string")},
				},
			},
			"PrintInt": ProcDecl{
				Name: "PrintInt",
				Args: []VarWithType{
					{"", TypeLiteral("int")},
				},
			},
			// FIXME: Remove this. Some of the invalidprogram examples depend on it now.
			"print": FuncDecl{Name: "print"},
		},
		Mutables: make(map[string]VarWithType),
		Types: map[string]TypeDefn{
			"int":    TypeDefn{TypeLiteral("int"), TypeLiteral("int"), nil},
			"uint":   TypeDefn{TypeLiteral("uint"), TypeLiteral("uint"), nil},
			"uint8":  TypeDefn{TypeLiteral("uint8"), TypeLiteral("uint8"), nil},
			"uint16": TypeDefn{TypeLiteral("uint16"), TypeLiteral("uint16"), nil},
			"uint32": TypeDefn{TypeLiteral("uint32"), TypeLiteral("uint32"), nil},
			"uint64": TypeDefn{TypeLiteral("uint64"), TypeLiteral("uint64"), nil},
			"int8":   TypeDefn{TypeLiteral("int8"), TypeLiteral("int8"), nil},
			"int16":  TypeDefn{TypeLiteral("int16"), TypeLiteral("int16"), nil},
			"int32":  TypeDefn{TypeLiteral("int32"), TypeLiteral("int32"), nil},
			"int64":  TypeDefn{TypeLiteral("int64"), TypeLiteral("int64"), nil},
			"string": TypeDefn{TypeLiteral("string"), TypeLiteral("string"), nil},
			"bool":   TypeDefn{TypeLiteral("bool"), TypeLiteral("bool"), nil},
		},
		PureContext: false,
		EnumOptions: make(map[string]EnumOption),
	}
}
func (c Context) Clone() Context {
	var c2 Context
	c2.Variables = make(map[string]VarWithType)
	c2.Functions = make(map[string]Callable)
	c2.Mutables = make(map[string]VarWithType)
	c2.Types = make(map[string]TypeDefn)
	c2.EnumOptions = make(map[string]EnumOption)
	for k, v := range c.Variables {
		c2.Variables[k] = v
	}
	for k, v := range c.Mutables {
		c2.Mutables[k] = v
	}

	for k, v := range c.Functions {
		c2.Functions[k] = v
	}
	for k, v := range c.Types {
		c2.Types[k] = v
	}
	for k, v := range c.EnumOptions {
		c2.EnumOptions[k] = v
	}
	c2.PureContext = c.PureContext
	return c2
}

func (c Context) IsVariable(s string) bool {
	for _, v := range c.Variables {
		if string(v.Name) == s {
			return true
		}
	}
	return false
}

func (c Context) IsMutable(s string) bool {
	for _, v := range c.Mutables {
		if string(v.Name) == s {
			return true
		}
	}
	return false
}

func (c Context) IsFunction(s string) bool {
	for k := range c.Functions {
		if k == s {
			return true
		}
	}
	return false
}

func (c Context) ValidType(t Type) bool {
	if _, ok := c.Types[t.Type()]; ok {
		return true
	}
	return false
}

func (c Context) EnumeratedOption(t string) *EnumOption {
	if eo, ok := c.EnumOptions[t]; ok {
		return &eo
	}
	return nil
}

func topLevelNode(T token.Token) (Node, error) {
	switch t := T.(type) {
	case token.Whitespace:
		return nil, nil
	case token.Keyword:
		switch t.String() {
		case "proc":
			return ProcDecl{}, nil
		case "func":
			return FuncDecl{}, nil
		case "type":
			return TypeDefn{}, nil
		case "data":
			return SumTypeDefn{}, nil
		}
		return nil, fmt.Errorf("Invalid top level keyword: %v", t)
	default:
		return nil, fmt.Errorf("Invalid top level token: %v %v", T, reflect.TypeOf(T))
	}
}

type TypeInfo struct {
	Size   int
	Signed bool
}

type TypeInformation map[string]TypeInfo

func Parse(val string) ([]Node, TypeInformation, error) {
	tokens, err := token.Tokenize(strings.NewReader(val))
	if err != nil {
		return nil, nil, err
	}
	return Construct(tokens)
}

func stripWhitespace(tokens []token.Token) []token.Token {
	// Remove all whitespace tokens to simplify the parsing. We don't care
	// about it anymore now that we've finished splitting into tokens.
	t2 := make([]token.Token, 0, len(tokens))
	for i := 0; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Whitespace:
			continue
		default:
			t2 = append(t2, t)
		}
	}
	return t2
}

// Construct constructs the top level ASTNodes for a file.
func Construct(tokens []token.Token) ([]Node, TypeInformation, error) {
	var nodes []Node
	ti := TypeInformation{
		("int"):     TypeInfo{8, true},
		("uint"):    TypeInfo{8, false},
		("int8"):    TypeInfo{1, true},
		("uint8"):   TypeInfo{1, false},
		("int16"):   TypeInfo{2, true},
		("uint16"):  TypeInfo{2, false},
		("int32"):   TypeInfo{4, true},
		("uint32"):  TypeInfo{4, false},
		("int64"):   TypeInfo{8, true},
		("uint64"):  TypeInfo{8, false},
		("bool"):    TypeInfo{1, false},
		("string"):  TypeInfo{0, false},
		("sumtype"): TypeInfo{8, false},
	}

	c := NewContext()

	tokens = stripWhitespace(tokens)
	// For debugging only.
	/*
		for i := 0; i < len(tokens); i++ {
			fmt.Fprintf(os.Stderr, "%d: '%v'\n", i, tokens[i].String())
		}
	*/

	err := extractPrototypes(tokens, &c)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(tokens); i++ {
		// Parse the top level "func" or "proc" keyword
		cn, err := topLevelNode(tokens[i])
		if err != nil {
			return nil, nil, err
		}

		switch cur := cn.(type) {
		case ProcDecl:
			// move past the "proc" keyword and reset the local
			// variables and mutables, since we're in a new function.
			c.Variables = make(map[string]VarWithType)
			c.Mutables = make(map[string]VarWithType)
			c.PureContext = false
			i++

			// FIXME: This should check that the name is valid and
			// i isn't out of bounds.
			cur.Name = tokens[i].String()
			i++

			n, a, r, err := consumePrototype(i, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			cur.Args = a
			cur.Return = r
			i += n

			for _, v := range cur.Args {
				c.Variables[string(v.Name)] = v
			}
			n, block, err := consumeBlock(i, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			cur.Body = block

			i += n

			nodes = append(nodes, cur)
		case FuncDecl:
			// move past the "func" keyword and reset the local
			// variables and mutables, since we're in a new function.
			c.Variables = make(map[string]VarWithType)
			c.Mutables = make(map[string]VarWithType)
			c.PureContext = true

			i++

			// FIXME: This should check that the name is valid and
			// i isn't out of bounds.
			cur.Name = tokens[i].String()
			i++

			n, a, r, err := consumePrototype(i, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			cur.Args = a
			cur.Return = r
			i += n
			for _, v := range cur.Args {
				c.Variables[string(v.Name)] = v
			}

			n, block, err := consumeBlock(i, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			cur.Body = block

			i += n
			nodes = append(nodes, cur)
		case TypeDefn:
			typeName := tokens[i+1].String()
			nodes = append(nodes, c.Types[typeName])
			switch concrete := c.Types[typeName].ConcreteType.Type(); concrete {
			case "int", "uint", "int8", "uint8", "int16", "uint16",
				"int32", "uint32", "int64", "uint64", "bool", "string",
				"sumtype":
				ti[typeName] = ti[c.Types[typeName].ConcreteType.Type()]
			default:
				panic("Unhandled concrete type: " + string(c.Types[typeName].ConcreteType.Type()))
			}
			i += 2
		case SumTypeDefn:
			n, typeNames, err := consumeIdentifiersUntilEquals(i+1, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			i += n + 1

			cur.Name = TypeLiteral(typeNames[0].String())
			n, options, err := consumeSumTypeList(i+1, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			for _, constructor := range options {
				constructor.ParentType = cur.Name
				cur.Options = append(cur.Options, constructor)
			}

			ti[cur.Name.Type()] = TypeInfo{8, false}

			i += n
			nodes = append(nodes, cur)
		}
	}
	return nodes, ti, nil
}

func consumePrototype(start int, tokens []token.Token, c *Context) (n int, args []VarWithType, retn []VarWithType, err error) {
	n, argsDefn, err := consumeArgs(start, tokens, c)
	if err != nil {
		return 0, nil, nil, err
	}

	n2, retDefn, err := consumeTypeList(start+n, tokens, *c)
	if err != nil {
		return 0, nil, nil, err
	}
	return n + n2, argsDefn, retDefn, nil
}

func extractPrototypes(tokens []token.Token, c *Context) error {
	// First pass: extract all the types, so that we can get the type
	// signatures on the second pass.
	for i := 0; i < len(tokens); i++ {
		// Parse the top level "func" or "proc" keyword
		cn, err := topLevelNode(tokens[i])
		if err != nil {
			return err
		}

		switch cur := cn.(type) {
		case ProcDecl:
			i++

			cur.Name = tokens[i].String()
			i++

			n, err := skipPrototype(i, tokens, c)
			if err != nil {
				return err
			}
			i += n

			n, err = skipBlock(i, tokens, c)
			if err != nil {
				return err
			}
			i += n

			// Now that we know the function is valid, add it to
			// the context's list of functions
			c.Functions[cur.Name] = cur
		case FuncDecl:
			i++

			cur.Name = tokens[i].String()
			i++

			n, err := skipPrototype(i, tokens, c)
			if err != nil {
				return err
			}
			i += n
			n, err = skipBlock(i, tokens, c)
			if err != nil {
				return err
			}
			i += n

			c.Functions[cur.Name] = cur
		case TypeDefn:
			i++

			cur.Name = TypeLiteral(tokens[i].String())
			i++
			cur.ConcreteType = TypeLiteral(tokens[i].String())
			c.Types[cur.Name.Type()] = cur
		case SumTypeDefn:
			n, typeNames, err := consumeIdentifiersUntilEquals(i+1, tokens, c)
			if err != nil {
				return err
			}
			i += n + 1

			cur.Name = TypeLiteral(typeNames[0].String())

			n, _, err = consumeSumTypeList(i+1, tokens, c)
			if err != nil {
				return err
			}

			// FIXME: This TypeLiteral is stupid and inaccurate now that Type() is an interface
			c.Types[cur.Name.Type()] = TypeDefn{ConcreteType: TypeLiteral("sumtype")}

			i += n
		}
	}

	// Second pass, extract the parameter lists of the functions, so that
	// we have all the information we need to validate function calls.
	for i := 0; i < len(tokens); i++ {
		// Parse the top level "func" or "proc" keyword
		cn, err := topLevelNode(tokens[i])
		if err != nil {
			return err
		}

		switch cur := cn.(type) {
		case ProcDecl:
			i++

			cur.Name = tokens[i].String()
			i++

			n, a, r, err := consumePrototype(i, tokens, c)
			if err != nil {
				return err
			}
			cur.Args = a
			cur.Return = r
			i += n

			n, err = skipBlock(i, tokens, c)
			if err != nil {
				return err
			}
			i += n

			// Now that we know the function is valid, add it to
			// the context's list of functions
			c.Functions[cur.Name] = cur
		case FuncDecl:
			i++

			cur.Name = tokens[i].String()
			i++

			n, a, r, err := consumePrototype(i, tokens, c)
			if err != nil {
				return err
			}
			cur.Args = a
			cur.Return = r
			i += n

			n, err = skipBlock(i, tokens, c)
			if err != nil {
				return err
			}
			i += n

			c.Functions[cur.Name] = cur
		case TypeDefn:
			i++

			cur.Name = TypeLiteral(tokens[i].String())
			i++
			cur.ConcreteType = TypeLiteral(tokens[i].String())
			c.Types[cur.Name.Type()] = cur
		case SumTypeDefn:
			n, typeNames, err := consumeIdentifiersUntilEquals(i+1, tokens, c)
			if err != nil {
				return err
			}
			i += n + 1

			cur.Name = TypeLiteral(typeNames[0].String())
			var pv []Type
			for _, param := range typeNames[1:] {
				pv = append(pv, TypeLiteral(param.String()))
			}
			n, options, err := consumeSumTypeList(i+1, tokens, c)
			if err != nil {
				return err
			}

			for _, o := range options {
				o.ParentType = cur.Name
				c.EnumOptions[o.Constructor] = o
			}

			c.Types[cur.Name.Type()] = TypeDefn{
				Name:         cur.Name,
				ConcreteType: TypeLiteral("sumtype"),
				Parameters:   pv,
			}

			i += n
		}
	}
	return nil
}

// consumeBlock consumes a balanced number of tokens delimited by a balanced
// number of "{" and "}" characters. It returns the ASTNode for the block, and
// the number of tokens that were consumed.
func consumeBlock(start int, tokens []token.Token, c *Context) (int, BlockStmt, error) {
	if tokens[start] != token.Char("{") {
		return 0, BlockStmt{}, fmt.Errorf("Invalid block. (%v)", tokens[start])
	}
	var blockStmt BlockStmt

	for i := start + 1; i < len(tokens); i++ {
		// First handle open or close brackets
		if tokens[i] == token.Char("{") {
			c2 := c.Clone()
			n, subblock, err := consumeBlock(i, tokens, &c2)
			if err != nil {
				return 0, BlockStmt{}, err
			}
			blockStmt.Stmts = append(blockStmt.Stmts, subblock)
			i += n
			continue
		} else if tokens[i] == token.Char("}") {
			return i - start, blockStmt, nil
		}
		switch t := tokens[i].(type) {
		case token.Unknown:
			if i+1 >= len(tokens) {
				return 0, BlockStmt{}, fmt.Errorf("Invalid token at end of file.")
			}
			switch tokens[i+1] {
			case token.Char("("):
				if c.IsFunction(tokens[i].String()) {
					n, funcCall, err := consumeFuncCall(i, tokens, c)
					if err != nil {
						return 0, BlockStmt{}, err
					}
					blockStmt.Stmts = append(blockStmt.Stmts, funcCall)
					i += n
				} else {
					return 0, BlockStmt{}, fmt.Errorf("Call to undefined function: %v", tokens[i])
				}
			case token.Char("["):
				// We're indexing into an array (probably)
				// use consumeValue to get the ArrayValue for the index
				// that we're checking.
				n, v, err := consumeValue(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}

				// After indexing into an array, the only operation that makes
				// sense is assignment.
				if tokens[n+i] != token.Operator("=") {
					return 0, BlockStmt{}, fmt.Errorf("Invalid variable assignment.")
				}

				av, ok := v.(ArrayValue)
				if !ok {
					return 0, BlockStmt{}, fmt.Errorf("Can not index non-array value")
				}

				basetype, ok := av.Base.Typ.(ArrayType)
				if !ok {
					return 0, BlockStmt{}, fmt.Errorf("Array type must be ArrayType")
				}

				// Adjust i to take into account what we consumed and the
				i += n + 1

				// Get the value that's being assigned.
				n, val, err := consumeValue(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				nm := Variable(fmt.Sprintf("%v[%d]", av.Base.Name, av.Index))
				blockStmt.Stmts = append(blockStmt.Stmts, AssignmentOperator{
					Variable: VarWithType{nm, basetype.Base},
					Value:    val,
				})
				i += n - 1
			case token.Operator("="):
				if !c.IsVariable(t.String()) {
					return 0, BlockStmt{}, fmt.Errorf("Invalid variable for assignment: %v", tokens[i])
				}

				if !c.IsMutable(t.String()) {
					return 0, BlockStmt{}, fmt.Errorf(`Can not assign to immutable let variable "%v".`, tokens[i])
				}
				n, val, err := consumeValue(i+2, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}

				blockStmt.Stmts = append(blockStmt.Stmts, AssignmentOperator{
					Variable: c.Variables[tokens[i].String()],
					Value:    val,
				})
				i += n + 1
			default:
				return 0, BlockStmt{}, fmt.Errorf("Don't know how to handle token: %v at token %d", tokens[i+1], i+1)
			}
		case token.Keyword:
			switch t {
			case "let":
				n, letstmt, err := consumeLetStmt(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				blockStmt.Stmts = append(blockStmt.Stmts, letstmt)
				i += n
			case "mutable":
				n, mutstmt, err := consumeMutStmt(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				blockStmt.Stmts = append(blockStmt.Stmts, mutstmt)
				i += n
			case "return":
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				blockStmt.Stmts = append(
					blockStmt.Stmts,
					ReturnStmt{
						Val: v,
					},
				)
				i += n
			case "while":
				n, v, err := consumeWhileLoop(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				blockStmt.Stmts = append(blockStmt.Stmts, v)
				i += n
			case "if":
				n, v, err := consumeIfStmt(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				i += n

				// This is more complicated than it should be.
				var firstIf *IfStmt = &v
				var lastIf *IfStmt = &v
				for i < len(tokens) && tokens[i+1] == token.Keyword("else") {
					if i < len(tokens)-2 && tokens[i+2] == token.Keyword("if") {
						n, nextIf, err := consumeIfStmt(i+2, tokens, c)

						if err != nil {
							return 0, BlockStmt{}, err
						}
						lastIf.Else = BlockStmt{
							[]Node{
								&nextIf,
							},
						}

						i += n + 2
						lastIf = &nextIf
					} else {
						n, elseB, err := consumeBlock(i+2, tokens, c)
						if err != nil {
							return 0, BlockStmt{}, err
						}
						lastIf.Else = elseB
						i += n + 2
					}
				}
				blockStmt.Stmts = append(blockStmt.Stmts, firstIf)
			case "match":
				n, v, err := consumeMatchStmt(i, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				blockStmt.Stmts = append(blockStmt.Stmts, v)
				i += n
			default:
				panic(fmt.Sprintf("Unimplemented keyword: %v at %v", tokens[i], i))
			}
		default:
			return 0, BlockStmt{}, fmt.Errorf("Unhandled token type in block %v for token %v", tokens[i].String(), i)
		}

	}
	return 0, BlockStmt{}, fmt.Errorf("Unterminated block statement")
}

func consumeStmt(start int, tokens []token.Token, c *Context) (int, Node, error) {
	switch t := tokens[start].(type) {
	case token.Unknown:
		if start+1 >= len(tokens) {
			return 0, BlockStmt{}, fmt.Errorf("Invalid token at end of file.")
		}
		switch tokens[start+1] {
		case token.Char("("):
			if c.IsFunction(tokens[start].String()) {
				n, fc, err := consumeFuncCall(start, tokens, c)
				if err == nil {
					return n + 1, fc, nil
				}
				return 0, nil, err
			} else {
				return 0, nil, fmt.Errorf("Call to undefined function: %v", tokens[start])
			}
		case token.Char("["):
			n, v, err := consumeValue(start, tokens, c)
			if err != nil {
				return 0, nil, err
			}
			panic(fmt.Sprintf("%v : %v", n, v))
		case token.Operator("="):
			if !c.IsVariable(t.String()) {
				return 0, nil, fmt.Errorf("Invalid variable for assignment: %v", tokens[start])
			}

			if !c.IsMutable(t.String()) {
				return 0, nil, fmt.Errorf(`Can not assign to immutable let variable "%v".`, tokens[start])
			}
			n, val, err := consumeValue(start+2, tokens, c)
			if err != nil {
				return 0, nil, err
			}

			return n, AssignmentOperator{
				Variable: c.Variables[tokens[start].String()],
				Value:    val,
			}, nil
		default:
			return 0, nil, fmt.Errorf("Don't know how to handle token: %v(%v) at token %d", reflect.TypeOf(tokens[start+1]), tokens[start+1], start+1)
		}
	case token.Keyword:
		switch t {
		case "let":
			return consumeLetStmt(start, tokens, c)
		case "mutable":
			return consumeMutStmt(start, tokens, c)
		case "return":
			n, nd, err := consumeValue(start+1, tokens, c)
			return n + 1, ReturnStmt{Val: nd}, err
		case "while":
			return consumeWhileLoop(start, tokens, c)
			/*
				case "if":
					n, v, err := consumeIfStmt(start, tokens, c)
					if err != nil {
						return 0, BlockStmt{}, err
					}
					start += n

					// This is more complicated than it should be.
					var firstIf *IfStmt = &v
					var lastIf *IfStmt = &v
					for i < len(tokens) && tokens[i+1] == token.Keyword("else") {
						if i < len(tokens)-2 && tokens[i+2] == token.Keyword("if") {
							n, nextIf, err := consumeIfStmt(i+2, tokens, c)
								if err != nil {
								return 0, BlockStmt{}, err
							}
							lastIf.Else = BlockStmt{
								[]Node{
											&nextIf,
										},
									}

									i += n + 2
									lastIf = &nextIf
								} else {
									n, elseB, err := consumeBlock(i+2, tokens, c)
									if err != nil {
										return 0, BlockStmt{}, err
									}
									lastIf.Else = elseB
									i += n + 2
								}
							}
							blockStmt.Stmts = append(blockStmt.Stmts, firstIf)
						case "match":
							n, v, err := consumeMatchStmt(i, tokens, c)
							if err != nil {
								return 0, BlockStmt{}, err
							}
							blockStmt.Stmts = append(blockStmt.Stmts, v)
							i += n
			*/
		default:
			panic(fmt.Sprintf("Unimplemented keyword: %v at %v", tokens[start], start))
		}
	default:
		panic(fmt.Sprintf("Unhandled token type in block %v for token %v", tokens[start].String(), start))
	}

}
func consumeFuncCall(start int, tokens []token.Token, c *Context) (int, FuncCall, error) {
	name := tokens[start].String()
	f := FuncCall{
		Name: name,
	}

	decl, ok := c.Functions[name]
	if !ok {
		return 0, FuncCall{}, fmt.Errorf("Undefined function: %v", name)
	}
	if c.PureContext {
		if _, ok := decl.(ProcDecl); ok {
			return 0, FuncCall{}, fmt.Errorf("Can not call procedure from pure function.")
		}
	}
	// FIXME: Hack because printf is variadic, and variadic functions aren't
	// implemented, but print is required for pretty much every test.
	if name == "print" {
		consumingString := false
		var strLit StringLiteral
		for i := start + 1; i < len(tokens); i++ {
			switch t := tokens[i].(type) {
			case token.Char:
				if t == ")" {
					return i - start, f, nil
				} else if t == `"` {
					if consumingString {
						f.UserArgs = append(f.UserArgs, strLit)
						strLit = ""
						consumingString = false
					} else {
						consumingString = true
					}
					continue
				} else if t == "," || t == "(" {
					continue
				} else {
					panic("Unexpected Char in function call" + t.String())
				}
			case token.Unknown:
				if c.IsVariable(t.String()) {
					f.UserArgs = append(f.UserArgs, c.Variables[t.String()])
				} else if c.IsFunction(t.String()) {
					n, subcall, err := consumeFuncCall(i, tokens, c)
					if err != nil {
						return 0, FuncCall{}, err
					}
					f.UserArgs = append(f.UserArgs, subcall)
					i += n
				} else {
					return 0, FuncCall{}, fmt.Errorf(`Use of undefined variable "%v".`, t)
				}
			case token.String:
			default:
				panic("Unexpected token consuming func call." + t.String())
			}
			if consumingString {
				strLit += StringLiteral(tokens[i].String())
			}
		}
	}

	// FIXME: This should support variadic functions, too.
	args := decl.GetArgs()

	f.Returns = decl.ReturnTuple()

	if tokens[start+1] == token.Char("(") && tokens[start+2] == token.Char(")") {
		if len(args) == 0 {
			return 2, f, nil
		}
		return 0, FuncCall{}, fmt.Errorf("Unexpected number of parameters to %v: got 0 want %v.", tokens[start], len(args))
	}

	argStart := start + 2 // start = name, +1 = "(", +2 = the first param..
argLoop:
	for {
		n, val, err := consumeValue(argStart, tokens, c)
		if err != nil {
			return 0, FuncCall{}, err
		}
		argStart += n
		switch tokens[argStart] {
		case token.Char(")"):
			argStart++
			f.UserArgs = append(f.UserArgs, val)
			break argLoop
		case token.Char(","):
			argStart++
			f.UserArgs = append(f.UserArgs, val)
		default:
			return 0, FuncCall{}, fmt.Errorf("Invalid token in %v. Expecting ')' or ',' in function argument list, got %v", name, tokens[argStart])
		}
	}
	if len(args) != len(f.UserArgs) {
		return 0, FuncCall{}, fmt.Errorf("Unexpected number of parameters to %v: got %v want %v.", tokens[start], len(f.UserArgs), len(args))
	}
	return argStart - start - 1, f, nil
}

func consumeLetStmt(start int, tokens []token.Token, c *Context) (int, Node, error) {
	l := LetStmt{}

	defer func() {
		c.Variables[l.Var.Name.String()] = l.Var
	}()
	if tokens[start] != token.Keyword("let") {
		return 0, nil, fmt.Errorf("Invalid let statement")
	}
	for i := start + 1; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			if l.Var.Name == "" {
				l.Var.Name = Variable(t.String())
				if _, ok := c.Mutables[t.String()]; ok {
					return 0, nil, fmt.Errorf("Can not shadow mutable variable \"%v\".", t.String())
				}
			} else if l.Var.Typ == nil {
				l.Var.Typ = TypeLiteral(t.String())
				if !c.ValidType(l.Var.Typ) {
					return 0, nil, fmt.Errorf("Invalid type: %v", t.String())
				}

			} else {
				return 0, nil, fmt.Errorf("Invalid name for let statement")
			}
		case token.Type:
			if l.Var.Typ == nil {
				l.Var.Typ = TypeLiteral(t.String())
			} else {
				return 0, nil, fmt.Errorf("Unexpected type in let statement")

			}
		case token.Char:
			if l.Var.Typ == nil && t == "[" {
				n, ty, err := consumeType(i, tokens, *c)
				if err != nil {
					return 0, nil, err
				}
				l.Var.Typ = ty
				i += n - 1
			} else {
				return 0, nil, fmt.Errorf("Unexpected character in let statement: %v(%v)", reflect.TypeOf(t), t)
			}
		case token.Operator:
			if t == token.Operator("=") {
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				if l.Var.Typ == nil {
					td := c.Types[v.Type()]
					switch td.ConcreteType.(type) {
					case ArrayType:
						l.Var.Typ = td.ConcreteType
					default:
						l.Var.Typ = TypeLiteral(v.Type())
					}
				}

				if IsLiteral(v) {
					if err := IsCompatibleType(c.Types[l.Type()], v); err != nil {
						return 0, nil, fmt.Errorf(`Incompatible assignment for variable "%v": %v.`, l.Var.Name, err)
					}
				} else {
					if v.Type() != l.Type() {
						return 0, nil, fmt.Errorf(`Incompatible assignment for variable "%v": can not assign %v to %v.`, l.Var.Name, v.Type(), l.Type())
					}
				}
				l.Value = v
				return i + n - start, l, nil
			}
			return 0, nil, fmt.Errorf("Invalid let statement")
		case token.Whitespace:
		default:
			return 0, nil, fmt.Errorf("Invalid let statement: %v", t)
		}

	}
	return 0, nil, fmt.Errorf("Invalid let statement")
}

func consumeMutStmt(start int, tokens []token.Token, c *Context) (int, Node, error) {
	l := MutStmt{}

	defer func() {
		c.Variables[l.Var.Name.String()] = l.Var
		c.Mutables[l.Var.Name.String()] = l.Var
	}()

	if tokens[start] != token.Keyword("mutable") {
		return 0, nil, fmt.Errorf("Invalid mutable variable statement")
	}
	for i := start + 1; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			if l.Var.Name == "" {
				l.Var.Name = Variable(t.String())
				if _, ok := c.Mutables[t.String()]; ok {
					return 0, nil, fmt.Errorf("Can not shadow mutable variable \"%v\".", t.String())
				}
			} else if l.Var.Typ == nil {
				l.Var.Typ = TypeLiteral(t.String())
				if !c.ValidType(l.Var.Typ) {
					return 0, nil, fmt.Errorf("Invalid type: %v", t.String())
				}
			} else {
				return 0, nil, fmt.Errorf("Invalid name for mutable declaration")
			}
		case token.Type:
			if l.Var.Typ == nil {
				l.Var.Typ = TypeLiteral(t.String())
			} else {
				return 0, nil, fmt.Errorf("Unexpected type in mutable declaration")

			}
		case token.Char:
			if l.Var.Typ == nil && t == "[" {
				n, ty, err := consumeType(i, tokens, *c)
				if err != nil {
					return 0, nil, err
				}
				l.Var.Typ = ty
				i += n - 1
			} else {
				return 0, nil, fmt.Errorf("Unexpected char in mutable declaration: %v(%v)", reflect.TypeOf(t), t)
			}
		case token.Operator:
			if t == token.Operator("=") {
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				if l.Var.Typ == nil {
					td := c.Types[v.Type()]
					switch td.ConcreteType.(type) {
					case ArrayType:
						l.Var.Typ = td.ConcreteType
					default:
						l.Var.Typ = TypeLiteral(v.Type())
					}
				}

				if IsLiteral(v) {
					if err := IsCompatibleType(c.Types[string(l.Type())], v); err != nil {
						return 0, nil, fmt.Errorf(`Incompatible assignment for variable "%v": %v.`, l.Var.Name, err)
					}
				} else {
					if v.Type() != l.Type() {
						return 0, nil, fmt.Errorf(`Incompatible assignment for variable "%v": can not assign %v to %v.`, l.Var.Name, v.Type(), l.Type())
					}
				}

				l.InitialValue = v
				return i + n - start, l, nil
			}
			return 0, nil, fmt.Errorf("Invalid mutable declaration")
		case token.Whitespace:
		default:
			return 0, nil, fmt.Errorf("Invalid mutable declaration: %v", t)
		}

	}
	return 0, nil, fmt.Errorf("Invalid mutable declaration")
}

func consumeArgs(start int, tokens []token.Token, c *Context) (int, []VarWithType, error) {
	var args []VarWithType
	started := false
	parsingNames := true
	var names []Variable
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Char:
			if t == "," {
				parsingNames = true
			} else if t == "(" {
				started = true
			} else if t == ")" {
				return i + 1 - start, args, nil
			} else {
				return 0, nil, fmt.Errorf("Invalid argument")
			}
			continue
		case token.Unknown:
			if parsingNames {
				names = append(names, Variable(t.String()))
				parsingNames = false
			} else {
				n, tk, err := consumeType(i, tokens, *c)
				if err != nil {
					return 0, nil, err
				}
				i += n - 1

				for _, n := range names {
					args = append(args, VarWithType{
						Name: n,
						Typ:  tk,
					})
				}

				names = nil
				parsingNames = true
			}
		case token.Type:
			if parsingNames {
				return 0, nil, fmt.Errorf("Expected name, got type")
			}
			for _, n := range names {
				args = append(args, VarWithType{
					Name: n,
					Typ:  TypeLiteral(t.String()),
				})
			}
			names = nil
			parsingNames = true
			continue
		default:
			return 0, nil, fmt.Errorf("Invalid token in argument list: %v %v", t, started)
		}
	}
	return 0, nil, fmt.Errorf("Could not parse arguments")
}

func consumeTypeList(start int, tokens []token.Token, c Context) (int, []VarWithType, error) {
	var args []VarWithType
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Char:
			if t == "," {
			} else if t == "(" {
			} else if t == ")" {
				return i + 1 - start, args, nil
			} else {
				return 0, nil, fmt.Errorf("Invalid argument")
			}
			continue
		case token.Unknown, token.Type:
			n, tk, err := consumeType(i, tokens, c)
			if err != nil {
				return 0, nil, err
			}
			i += n - 1
			args = append(args, VarWithType{"", tk})
			continue
		default:
			return 0, nil, fmt.Errorf("Invalid token in argument list: %v", t)
		}
	}
	return 0, nil, fmt.Errorf("Could not parse arguments")
}

func consumeType(start int, tokens []token.Token, c Context) (int, Type, error) {
	nm := tokens[start].String()
	if nm == "[" {
		if tokens[start+1] == token.Char("]") {
			panic("Variable length slices not implemented")
		}
		in, size, err := consumeValue(start+1, tokens, &c)
		if err != nil {
			return 0, nil, err
		}
		sz, ok := size.(IntLiteral)
		if !ok {
			return 0, nil, fmt.Errorf("Array size must be an int literal")
		}
		if tokens[start+in+1] != token.Char("]") {
			return 0, nil, fmt.Errorf("Unexpected token: %v", tokens[start+in+1])
		}
		n, base, err := consumeType(start+in+2, tokens, c)
		if err != nil {
			return 0, nil, err
		}
		tn := fmt.Sprintf("[%d]%v", sz, base)
		t := ArrayType{
			Base: base,
			Size: sz,
		}
		c.Types[tn] = TypeDefn{t, t, nil}
		return n + in + 2, t, nil
	}
	consumed := 1
	typedef := c.Types[nm]
	rv := TypeLiteral(nm)
	for range typedef.Parameters {
		n, t, err := consumeType(start+consumed, tokens, c)
		if err != nil {
			return 0, nil, err
		}
		consumed += n
		rv += TypeLiteral(" ") + TypeLiteral(t.Type())
	}
	return consumed, rv, nil
}

func consumeIdentifiersUntilEquals(start int, tokens []token.Token, c *Context) (int, []token.Token, error) {
	var vals []token.Token
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			vals = append(vals, t)
		case token.Operator:
			if t == "=" {
				return i - start, vals, nil
			}
			return 0, nil, fmt.Errorf("Unexpected operator: %v", t)
		default:
			return 0, nil, fmt.Errorf("Invalid token: %v", t)
		}
	}
	return 0, nil, fmt.Errorf("Could not parse identifiers")
}

func consumeSumTypeList(start int, tokens []token.Token, c *Context) (int, []EnumOption, error) {
	var vals []EnumOption
	var val EnumOption
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			if val.Constructor != "" {
				val.Parameters = append(val.Parameters, TypeLiteral(t.String()))
			} else {
				val = EnumOption{Constructor: t.String()}
			}

			if i+1 < len(tokens) && tokens[i+1] == token.Operator("|") {
				vals = append(vals, val)
				val = EnumOption{}
				i += 1
			}
		case token.Keyword:
			vals = append(vals, val)
			return i - start, vals, nil
		default:
			return 0, nil, fmt.Errorf("Invalid token in sumtype: %v", t)
		}
	}
	vals = append(vals, val)
	return len(tokens) - start, vals, nil
}

func skipBlock(start int, tokens []token.Token, c *Context) (int, error) {
	i := start
	if tokens[i] != token.Char("{") {
		return 0, fmt.Errorf("Can not skip block. Not a block start.")
	}

	// Skip over the block, we're only trying to
	// extract the prototype.
	blockLevel := 1
	for i++; blockLevel > 0; i++ {
		if i >= len(tokens) {
			return 0, fmt.Errorf("Missing closing bracket for block (1): %v >= %v", i, len(tokens))
		}
		switch tokens[i] {
		case token.Char("{"):
			blockLevel++
		case token.Char("}"):
			blockLevel--
		}
	}
	if blockLevel > 0 {
		return 0, fmt.Errorf("Missing closing bracket for block (2)")
	}
	return i - 1 - start, nil
}

func skipPrototype(start int, tokens []token.Token, c *Context) (int, error) {
	n, err := skipTuple(start, tokens, c)
	if err != nil {
		return 0, err
	}
	n2, err := skipTuple(start+n, tokens, c)
	if err != nil {
		return 0, err
	}
	return n + n2, nil
}

func skipTuple(start int, tokens []token.Token, c *Context) (int, error) {
	i := start
	if tokens[i] != token.Char("(") {
		return 0, fmt.Errorf("Can not skip tuple. Expecting '(', not %v", tokens[i])
	}

	for ; i < len(tokens); i++ {
		if tokens[i] == token.Char(")") {
			return i + 1 - start, nil
		}
	}
	return 0, fmt.Errorf("Missing closing ')' for tuple.")
}
