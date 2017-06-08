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
			"print": FuncDecl{Name: "print"},
		},
		Mutables: make(map[string]VarWithType),
		Types: map[string]TypeDefn{
			"int":    TypeDefn{"int", "int"},
			"uint":   TypeDefn{"uint", "uint"},
			"uint8":  TypeDefn{"uint8", "uint8"},
			"uint16": TypeDefn{"uint16", "uint16"},
			"uint32": TypeDefn{"uint32", "uint32"},
			"uint64": TypeDefn{"uint64", "uint64"},
			"int8":   TypeDefn{"int8", "int8"},
			"int16":  TypeDefn{"int16", "int16"},
			"int32":  TypeDefn{"int32", "int32"},
			"int64":  TypeDefn{"int64", "int64"},
			"string": TypeDefn{"string", "string"},
			"bool":   TypeDefn{"bool", "bool"},
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
	if _, ok := c.Types[string(t)]; ok {
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

type TypeInformation map[Type]TypeInfo

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
		"int":    TypeInfo{8, true},
		"uint":   TypeInfo{8, false},
		"int8":   TypeInfo{1, true},
		"uint8":  TypeInfo{1, false},
		"int16":  TypeInfo{2, true},
		"uint16": TypeInfo{2, false},
		"int32":  TypeInfo{4, true},
		"uint32": TypeInfo{4, false},
		"int64":  TypeInfo{8, true},
		"uint64": TypeInfo{8, false},
		"bool":   TypeInfo{1, false},
		"string": TypeInfo{0, false},
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
			switch concrete := c.Types[typeName].ConcreteType; concrete {
			case "int", "uint", "int8", "uint8", "int16", "uint16",
				"int32", "uint32", "int64", "uint64", "bool", "string":
				ti[Type(typeName)] = ti[concrete]
			default:
				panic("Unhandled concrete type: " + string(c.Types[typeName].ConcreteType))
			}
			i += 2
		case SumTypeDefn:
			n, typeNames, err := consumeIdentifiersUntilEquals(i+1, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			if len(typeNames) != 1 {
				panic("Generic sum types not yet implemented.")
			}
			i += n + 1

			cur.Name = Type(typeNames[0].String())
			n, options, err := consumeSumTypeList(i+1, tokens, &c)
			if err != nil {
				return nil, nil, err
			}
			for _, constructor := range options {
				cur.Options = append(cur.Options, EnumOption{constructor, cur.Name})
			}
			// FIXME: This type info is completely wrong.
			ti[cur.Name] = TypeInfo{8, false}

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

	n2, retDefn, err := consumeTypeList(start+n, tokens)
	if err != nil {
		return 0, nil, nil, err
	}
	return n + n2, argsDefn, retDefn, nil
}

func extractPrototypes(tokens []token.Token, c *Context) error {
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

			if tokens[i] != token.Char("{") {
				return fmt.Errorf("Invalid syntax for %v. Missing body.", cur.Name)
			}

			// Skip over the block, we're only trying to
			// extract the prototype.
			blockLevel := 1
			for i++; blockLevel > 0; i++ {
				if i >= len(tokens) {
					return fmt.Errorf("Missing closing bracket for %v", cur.Name)
				}
				switch tokens[i] {
				case token.Char("{"):
					blockLevel++
				case token.Char("}"):
					blockLevel--
				}
			}
			// the for loop overshot by one token before doing the
			// comparison.
			i--

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

			if tokens[i] != token.Char("{") {
				return fmt.Errorf("Invalid syntax for %v. Missing body.", cur.Name)
			}

			// Skip over the block, we're only trying to
			// extract the prototype.
			blockLevel := 1
			for i++; blockLevel > 0; i++ {
				if i >= len(tokens) {
					return fmt.Errorf("Missing closing bracket for %v", cur.Name)
				}
				switch tokens[i] {
				case token.Char("{"):
					blockLevel++
				case token.Char("}"):
					blockLevel--
				}
			}
			i--
			c.Functions[cur.Name] = cur
		case TypeDefn:
			i++

			cur.Name = Type(tokens[i].String())
			i++
			cur.ConcreteType = Type(tokens[i].String())
			c.Types[string(cur.Name)] = cur
		case SumTypeDefn:
			n, typeNames, err := consumeIdentifiersUntilEquals(i+1, tokens, c)
			if err != nil {
				return err
			}
			if len(typeNames) != 1 {
				panic("Generic sum types not yet implemented.")
			}
			i += n + 1

			cur.Name = Type(typeNames[0].String())

			n, options, err := consumeSumTypeList(i+1, tokens, c)
			if err != nil {
				return err
			}

			for _, o := range options {
				c.EnumOptions[o] = EnumOption{o, cur.Name}
			}
			// FIXME: There needs to be a type interface so that
			// SumTypeDefn and TypeDefn can be in the same map.
			c.Types[string(cur.Name)] = TypeDefn{Name: "sumtype"}

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
				panic(fmt.Sprintf("Don't know how to handle token: %v at token %d", tokens[i+1], i+1))
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
			case "mut":
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
			panic(fmt.Sprintf("Unhandled token type in block %v for token %v", tokens[i].String(), i))
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
			panic(fmt.Sprintf("Don't know how to handle token: %v at token %d", tokens[start+1], start+1))
		}
	case token.Keyword:
		switch t {
		case "let":
			return consumeLetStmt(start, tokens, c)
		case "mut":
			return consumeMutStmt(start, tokens, c)
		case "return":
			return consumeValue(start+1, tokens, c)
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
			return 0, FuncCall{}, fmt.Errorf("Invalid token in %v. Expecting ')' or ',' in function argument list")
		}
	}
	if len(args) != len(f.UserArgs) {
		return 0, FuncCall{}, fmt.Errorf("Unexpected number of parameters to %v: got %v want %v.", tokens[start], len(f.UserArgs), len(args))
	}
	return argStart - start - 1, f, nil
}

func consumeLetStmt(start int, tokens []token.Token, c *Context) (int, Node, error) {
	l := LetStmt{}

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
				c.Variables[t.String()] = l.Var
			} else if l.Var.Typ == "" {
				l.Var.Typ = Type(t.String())
				if !c.ValidType(l.Var.Typ) {
					return 0, nil, fmt.Errorf("Invalid type: %v", t.String())
				}

				c.Variables[string(l.Var.Name)] = l.Var
			} else {
				return 0, nil, fmt.Errorf("Invalid name for let statement")
			}
		case token.Type:
			if l.Var.Typ == "" {
				l.Var.Typ = Type(t.String())
				c.Variables[string(l.Var.Name)] = l.Var

			} else {
				return 0, nil, fmt.Errorf("Unexpected type in let statement")

			}
		case token.Operator:
			if t == token.Operator("=") {
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				if l.Var.Type() == "" {
					l.Var.Typ = v.Type()
					c.Variables[string(l.Var.Name)] = l.Var
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

	if tokens[start] != token.Keyword("mut") {
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
				c.Variables[t.String()] = l.Var
				c.Mutables[t.String()] = l.Var
			} else if l.Var.Typ == "" {
				l.Var.Typ = Type(t.String())
				if !c.ValidType(l.Var.Typ) {
					return 0, nil, fmt.Errorf("Invalid type: %v", t.String())
				}
				c.Variables[string(l.Var.Name)] = l.Var
				c.Mutables[string(l.Var.Name)] = l.Var
			} else {
				return 0, nil, fmt.Errorf("Invalid name for mut statement")
			}
		case token.Type:
			if l.Var.Typ == "" {
				l.Var.Typ = Type(t.String())
				c.Variables[string(l.Var.Name)] = l.Var
				c.Mutables[string(l.Var.Name)] = l.Var
			} else {
				return 0, nil, fmt.Errorf("Unexpected type in mut statement")

			}
		case token.Operator:
			if t == token.Operator("=") {
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				if l.Var.Type() == "" {
					l.Var.Typ = v.Type()
					c.Variables[string(l.Var.Name)] = l.Var
					c.Mutables[string(l.Var.Name)] = l.Var
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
			return 0, nil, fmt.Errorf("Invalid mut statement")
		case token.Whitespace:
		default:
			return 0, nil, fmt.Errorf("Invalid mut statement: %v", t)
		}

	}
	return 0, nil, fmt.Errorf("Invalid mut statement")
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
			}
		case token.Type:
			if parsingNames {
				return 0, nil, fmt.Errorf("Expected name, got type")
			}
			for _, n := range names {
				args = append(args, VarWithType{
					Name: n,
					Typ:  Type(t.String()),
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

func consumeTypeList(start int, tokens []token.Token) (int, []VarWithType, error) {
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
			args = append(args, VarWithType{
				Name: "",
				Typ:  Type(t.String()),
			})
			continue
		default:
			return 0, nil, fmt.Errorf("Invalid token in argument list: %v", t)
		}
	}
	return 0, nil, fmt.Errorf("Could not parse arguments")
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

func consumeSumTypeList(start int, tokens []token.Token, c *Context) (int, []string, error) {
	var vals []string
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			vals = append(vals, t.String())
			if tokens[i+1] == token.Operator("|") {
				i += 1
			}
		case token.Keyword:
			return i - start, vals, nil
		default:
			return 0, nil, fmt.Errorf("Invalid token in sumtype: %v", t)
		}
	}
	return 0, nil, fmt.Errorf("Could not consume type list")
}
