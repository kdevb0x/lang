package ast

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/driusan/lang/parser/token"
)

var debug = false

type Context struct {
	Variables   map[string]VarWithType
	Mutables    map[string]VarWithType
	Functions   map[string]Callable
	Types       map[string]TypeDefn
	PureContext bool // true if inside a pure function.
	EnumOptions map[string]EnumOption
	CurFunc     Callable
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
					{"str", TypeLiteral("string"), false},
				},
			},
			"PrintInt": ProcDecl{
				Name: "PrintInt",
				Args: []VarWithType{
					{"x", TypeLiteral("int"), false},
				},
			},
			"PrintByteSlice": ProcDecl{
				Name: "PrintByteSlice",
				Args: []VarWithType{
					{"slice", SliceType{TypeLiteral("byte")}, false},
				},
			},
			"len": FuncDecl{
				Name: "len",
				Args: []VarWithType{
					// FIXME: This should be any slice, not just string slices, but
					// string slices take priority until there's some sort of generic func
					// decl, because they're passed to main..
					{"slice", SliceType{TypeLiteral("string")}, false},
				},
				Return: []VarWithType{
					{"", TypeLiteral("uint64"), false},
				},
			},
			// FIXME: These should be moved out of the compiler
			// and into a standard library, once enough of the
			// compiler is implemented to have a standard
			// library.
			"Write": ProcDecl{
				Name: "Write",
				Args: []VarWithType{
					{"fd", TypeLiteral("uint64"), false},
					{"val", SliceType{TypeLiteral("byte")}, false},
				},
			},
			"Read": ProcDecl{
				Args: []VarWithType{
					{"fd", TypeLiteral("uint64"), false},
					// NB. this should be []byte, once arrays are implemented.
					// NB2. This will read exactly the length of string from fd int
					//     dst and overwrite what's there.
					//     It needs a way to mark parameters mutable.
					{"dst", SliceType{TypeLiteral("byte")}, true},
				},
				Return: []VarWithType{{"", TypeLiteral("uint64"), false}},
			},
			"Open": ProcDecl{
				Name: "Open",
				Args: []VarWithType{
					{"val", TypeLiteral("string"), false},
					/*
						Follow Go/Plan9 conventions. Open just opens, Create
						just creates with a default umask, rather than having
						options on a generic Open like POSIX does.
						This makes it easier to port to Plan 9.
						{"mode", TypeLiteral("int"), false},
						{"createperms", TypeLiteral("int"), false},
					*/
				},
				Return: []VarWithType{{"", TypeLiteral("uint64"), false}},
			},
			"Create": ProcDecl{
				Name: "Create",
				Args: []VarWithType{
					{"val", TypeLiteral("string"), false},
					/*
						Follow Go/Plan9 conventions. Open just opens, Create
						just creates with a default umask, rather than having
						options on a generic Open like POSIX does.
						This makes it easier to port to Plan 9.
						{"mode", TypeLiteral("int"), false},
						{"createperms", TypeLiteral("int"), false},
					*/
				},
				Return: []VarWithType{{"", TypeLiteral("uint64"), false}},
			},
			"Close": ProcDecl{
				Name: "Close",
				Args: []VarWithType{
					{"val", TypeLiteral("uint64"), false},
				},
			},
		},
		Mutables: make(map[string]VarWithType),
		Types: map[string]TypeDefn{
			"int":    TypeDefn{TypeLiteral("int"), TypeLiteral("int"), nil},
			"uint":   TypeDefn{TypeLiteral("uint"), TypeLiteral("uint"), nil},
			"uint8":  TypeDefn{TypeLiteral("uint8"), TypeLiteral("uint8"), nil},
			"byte":   TypeDefn{TypeLiteral("byte"), TypeLiteral("byte"), nil},
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
	c2.CurFunc = c.CurFunc
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
type Callables map[string][]Callable

func Parse(val string) ([]Node, TypeInformation, Callables, error) {
	tokens, err := token.Tokenize(strings.NewReader(val))
	if err != nil {
		return nil, nil, nil, err
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
func Construct(tokens []token.Token) ([]Node, TypeInformation, Callables, error) {
	var nodes []Node
	ti := TypeInformation{
		("int"):     TypeInfo{0, true},
		("uint"):    TypeInfo{0, false},
		("int8"):    TypeInfo{1, true},
		("uint8"):   TypeInfo{1, false},
		("byte"):    TypeInfo{1, false},
		("int16"):   TypeInfo{2, true},
		("uint16"):  TypeInfo{2, false},
		("int32"):   TypeInfo{4, true},
		("uint32"):  TypeInfo{4, false},
		("int64"):   TypeInfo{8, true},
		("uint64"):  TypeInfo{8, false},
		("bool"):    TypeInfo{1, false},
		("string"):  TypeInfo{0, false},
		("sumtype"): TypeInfo{4, false},
	}

	c := NewContext()

	tokens = stripWhitespace(tokens)
	if debug {
		for i := 0; i < len(tokens); i++ {
			fmt.Fprintf(os.Stderr, "%d: '%v'\n", i, tokens[i].String())
		}
	}

	callables := make(Callables)
	for k, v := range c.Functions {
		callables[k] = append(callables[k], v)
	}
	err := extractPrototypes(tokens, &c)
	if err != nil {
		return nil, nil, nil, err
	}

	for i := 0; i < len(tokens); i++ {
		// Parse the top level "func" or "proc" keyword
		cn, err := topLevelNode(tokens[i])
		if err != nil {
			return nil, nil, nil, err
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
				return nil, nil, nil, err
			}
			cur.Args = a
			cur.Return = r
			i += n

			for _, v := range cur.Args {
				c.Variables[string(v.Name)] = v
				if v.Reference {
					c.Mutables[string(v.Name)] = v
				}
			}
			c.CurFunc = cur
			n, block, err := consumeBlock(i, tokens, &c)
			if err != nil {
				return nil, nil, nil, err
			}
			cur.Body = block

			i += n - 1

			nodes = append(nodes, cur)
			callables[cur.Name] = append(callables[cur.Name], cur)
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
				return nil, nil, nil, err
			}
			cur.Args = a
			cur.Return = r
			c.CurFunc = cur
			i += n
			for _, v := range cur.Args {
				c.Variables[string(v.Name)] = v
			}

			n, block, err := consumeBlock(i, tokens, &c)
			if err != nil {
				return nil, nil, nil, err
			}
			cur.Body = block

			i += n - 1
			nodes = append(nodes, cur)
			callables[cur.Name] = append(callables[cur.Name], cur)
		case TypeDefn:
			typeName := tokens[i+1].String()
			nodes = append(nodes, c.Types[typeName])
			switch concrete := c.Types[typeName].ConcreteType.Type(); concrete {
			case "int", "uint", "int8", "uint8", "int16", "uint16",
				"int32", "uint32", "int64", "uint64", "bool", "string",
				"sumtype", "byte":
				ti[typeName] = ti[c.Types[typeName].ConcreteType.Type()]
			default:
				panic("Unhandled concrete type: " + string(c.Types[typeName].ConcreteType.Type()))
			}
			i += 2
		case SumTypeDefn:
			n, typeNames, err := consumeIdentifiersUntilEquals(i+1, tokens, &c)
			if err != nil {
				return nil, nil, nil, err
			}
			i += n + 1

			cur.Name = TypeLiteral(typeNames[0].String())
			n, options, err := consumeSumTypeList(i+1, tokens, &c)
			if err != nil {
				return nil, nil, nil, err
			}
			for _, constructor := range options {
				constructor.ParentType = cur.Name
				cur.Options = append(cur.Options, constructor)
			}

			ti[cur.Name.Type()] = TypeInfo{0, false}

			i += n
			nodes = append(nodes, cur)
		}
	}
	return nodes, ti, callables, nil
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
			return i + 1 - start, blockStmt, nil
		}
		n, stmt, err := consumeStmt(i, tokens, c)
		if err != nil {
			return 0, BlockStmt{}, err
		}
		blockStmt.Stmts = append(blockStmt.Stmts, stmt)
		i += n - 1
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
				n, fc, err := consumeFuncCall(start, tokens, c, nil)
				if err == nil {
					return n + 1, fc, nil
				}
				return 0, nil, err
			} else {
				return 0, nil, fmt.Errorf("Call to undefined function: %v", tokens[start])
			}
		case token.Char("."):
			n, fc, err := consumeFuncCall(start+2, tokens, c, []Value{c.Variables[tokens[start].String()]})
			if err != nil {
				return 0, nil, err
			}
			return n + 3, fc, nil
		case token.Char("["):
			// We're indexing into an array (probably)
			// use consumeValue to get the ArrayValue for the index
			// that we're checking.
			n, v, err := consumeValue(start, tokens, c)
			if err != nil {
				return 0, BlockStmt{}, err
			}

			// After indexing into an array, the only operation that makes
			// sense is assignment. consumeValue would have taken care of infix
			// operations above.
			if tokens[start+n] != token.Operator("=") {
				return 0, BlockStmt{}, fmt.Errorf("Invalid variable assignment.")
			}

			av, ok := v.(ArrayValue)
			if !ok {
				return 0, BlockStmt{}, fmt.Errorf("Can not index non-array value")
			}

			valn, val, err := consumeValue(start+n+1, tokens, c)
			if err != nil {
				return 0, BlockStmt{}, err
			}
			return n + valn + 1, AssignmentOperator{
				Variable: av,
				Value:    val,
			}, nil
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

			// n for the value, one for the token, one for the = sign.
			return n + 2, AssignmentOperator{
				Variable: c.Variables[tokens[start].String()],
				Value:    val,
			}, nil
		default:
			return 0, nil, fmt.Errorf("Don't know how to handle token: %v(%v) at token %d [%v]", reflect.TypeOf(tokens[start+1]), tokens[start+1], start+1, tokens[start+1:])
		}
	case token.Keyword:
		switch t {
		case "let":
			return consumeLetStmt(start, tokens, c)
		case "mutable":
			return consumeMutStmt(start, tokens, c)
		case "return":
			if len(c.CurFunc.ReturnTuple()) == 0 {
				return 1, ReturnStmt{}, nil
			}

			n, nd, err := consumeValue(start+1, tokens, c)
			if err != nil {
				return 0, ReturnStmt{}, err
			}
			return n + 1, ReturnStmt{Val: nd}, nil
		case "while":
			return consumeWhileLoop(start, tokens, c)
		case "if":
			return consumeIfStmt(start, tokens, c)
		case "match":
			return consumeMatchStmt(start, tokens, c)
		default:
			panic(fmt.Sprintf("Unimplemented keyword: %v at %v", tokens[start], start))
		}
	default:
		panic(fmt.Sprintf("Unhandled token type in block %v for token %v [%v]", tokens[start].String(), start, tokens[start:]))
	}

}
func consumeFuncCall(start int, tokens []token.Token, c *Context, mvals []Value) (int, FuncCall, error) {
	name := tokens[start].String()
	f := FuncCall{
		Name:     name,
		UserArgs: mvals,
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

	// FIXME: This should support variadic functions, too.
	args := decl.GetArgs()

	f.Returns = decl.ReturnTuple()

	if tokens[start+1] == token.Char("(") && tokens[start+2] == token.Char(")") {
		if len(args) == len(mvals) {
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
	// Check that the arguments we got were compatible.
	// As a temporary hack, we don't check PrintInt or len, because PrintInt
	// currently deals with all int types and len deals with both strings and
	// slices, and there's not yet any casting
	if name != "PrintInt" && name != "len" {
		for i, arg := range args {
			if IsLiteral(f.UserArgs[i]) {
				if err := IsCompatibleType(c.Types[arg.Type()], f.UserArgs[i]); err != nil {
					return 0, FuncCall{}, fmt.Errorf("Incompatible call to %v: argument %v must be of type %v (got %v)", name, arg.Name, arg.Type(), f.UserArgs[i].Type())
				}
			} else {
				if arg.Type() != f.UserArgs[i].Type() {
					return 0, FuncCall{}, fmt.Errorf("Incompatible call to %v: argument %v must be of type %v (got %v)", name, arg.Name, arg.Type(), f.UserArgs[i].Type())
				}
			}
		}
	}
	return argStart - start - 1, f, nil
}

func consumeLetStmt(start int, tokens []token.Token, c *Context) (int, Value, error) {
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
				n, ty, err := consumeType(i, tokens, c)
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
					case ArrayType, SliceType:
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
				l.Val = v
				return i + n - start + 1, l, nil
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
				n, ty, err := consumeType(i, tokens, c)
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
					case ArrayType, SliceType:
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
				return i + n - start + 1, l, nil
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
	mutable := false
	for i := start; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Char:
			if t == "," {
				parsingNames = true
				mutable = false
				continue
			} else if t == "(" {
				started = true
				continue

			} else if t == ")" {
				return i + 1 - start, args, nil
			} else if t != "[" {
				// If it's a [ char, treat it as if it was the
				// start of a type below
				return 0, nil, fmt.Errorf("Invalid argument")
			}

			n, tk, err := consumeType(i, tokens, c)
			if err != nil {
				return 0, nil, err
			}
			i += n - 1

			for _, n := range names {
				args = append(args, VarWithType{
					Name:      n,
					Typ:       tk,
					Reference: mutable,
				})
			}

			names = nil
			parsingNames = true
			mutable = false
		case token.Unknown:
			if parsingNames {
				names = append(names, Variable(t.String()))
				parsingNames = false
			} else {
				n, tk, err := consumeType(i, tokens, c)
				if err != nil {
					return 0, nil, err
				}
				i += n - 1

				for _, n := range names {
					args = append(args, VarWithType{
						Name:      n,
						Typ:       tk,
						Reference: mutable,
					})
				}

				names = nil
				parsingNames = true
				mutable = false
			}
		case token.Type:
			if parsingNames {
				return 0, nil, fmt.Errorf("Expected name, got type")
			}
			for _, n := range names {
				args = append(args, VarWithType{
					Name:      n,
					Typ:       TypeLiteral(t.String()),
					Reference: mutable,
				})
			}
			names = nil
			parsingNames = true
			mutable = false
			continue
		case token.Keyword:
			if t != "mutable" {
				return 0, nil, fmt.Errorf("Unexpected keyword in argument list: %v %v", t, started)
			}
			mutable = true
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
			n, tk, err := consumeType(i, tokens, &c)
			if err != nil {
				return 0, nil, err
			}
			i += n - 1
			args = append(args, VarWithType{"", tk, false})
			continue
		default:
			return 0, nil, fmt.Errorf("Invalid token in argument list: %v", t)
		}
	}
	return 0, nil, fmt.Errorf("Could not parse arguments")
}

func consumeType(start int, tokens []token.Token, c *Context) (int, Type, error) {
	nm := tokens[start].String()
	if nm == "[" {
		if tokens[start+1] == token.Char("]") {
			// It is a slice.
			n, base, err := consumeType(start+2, tokens, c)
			if err != nil {
				return 0, nil, err
			}
			tn := fmt.Sprintf("[]%v", base)
			t := SliceType{
				Base: base,
			}
			c.Types[tn] = TypeDefn{t, t, nil}
			return n + 2, t, nil
		}
		// It is an array
		in, size, err := consumeValue(start+1, tokens, c)
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
