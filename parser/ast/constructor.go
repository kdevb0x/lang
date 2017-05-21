package ast

import (
	"fmt"
	"github.com/driusan/lang/parser/token"
	// "os"
	"strings"
)

type Context struct {
	Variables map[string]bool
	Mutables  map[string]bool
	Functions map[string]Callable
}

func NewContext() Context {
	return Context{
		Variables: make(map[string]bool),
		Functions: map[string]Callable{
			"print": FuncDecl{Name: "print"},
		},
		Mutables: make(map[string]bool),
	}
}
func (c Context) Clone() Context {
	var c2 Context
	c.Variables = make(map[string]bool)
	c.Functions = make(map[string]Callable)
	c.Mutables = make(map[string]bool)
	for k, v := range c.Variables {
		c2.Variables[k] = v
	}
	for k, v := range c.Mutables {
		c2.Mutables[k] = v
	}

	for k, v := range c.Functions {
		c2.Functions[k] = v
	}
	return c2
}

func (c Context) IsVariable(s string) bool {
	for k := range c.Variables {
		if k == s {
			return true
		}
	}
	return false
}

func (c Context) IsMutable(s string) bool {
	for k := range c.Mutables {
		if k == s {
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

func topLevelNode(T token.Token) (Node, error) {
	switch T.(type) {
	case token.Whitespace:
		return nil, nil
	default:
		if T.String() == "proc" {
			return ProcDecl{}, nil
		} else if T.String() == "func" {
			return FuncDecl{}, nil
		}
	}
	return nil, fmt.Errorf("Invalid top level token %v", T)
}

func Parse(val string) ([]Node, error) {
	tokens, err := token.Tokenize(strings.NewReader(val))
	if err != nil {
		return nil, err
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
func Construct(tokens []token.Token) ([]Node, error) {
	var nodes []Node

	c := NewContext()

	tokens = stripWhitespace(tokens)
	/*
		// For debugging only.
		for i := 0; i < len(tokens); i++ {
			fmt.Fprintf(os.Stderr, "%d: '%v'\n", i, tokens[i].String())
		}
	*/

	err := extractPrototypes(tokens, c)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tokens); i++ {
		// Parse the top level "func" or "proc" keyword
		cn, err := topLevelNode(tokens[i])
		if err != nil {
			return nil, err
		}

		switch cur := cn.(type) {
		case ProcDecl:
			// move past the "proc" keyword and reset the local
			// variables and mutables, since we're in a new function.
			c.Variables = make(map[string]bool)
			c.Mutables = make(map[string]bool)
			i++

			// FIXME: This should check that the name is valid and
			// i isn't out of bounds.
			cur.Name = tokens[i].String()
			i++

			n, a, r, err := consumePrototype(i, tokens, c)
			if err != nil {
				return nil, err
			}
			cur.Args = a
			cur.Return = r
			i += n

			n, block, err := consumeBlock(i, tokens, c)
			if err != nil {
				return nil, err
			}
			cur.Body = block

			i += n

			nodes = append(nodes, cur)
		case FuncDecl:
			// move past the "proc" keyword and reset the local
			// variables and mutables, since we're in a new function.
			c.Variables = make(map[string]bool)
			c.Mutables = make(map[string]bool)
			i++

			// FIXME: This should check that the name is valid and
			// i isn't out of bounds.
			cur.Name = tokens[i].String()
			i++

			n, a, r, err := consumePrototype(i, tokens, c)
			if err != nil {
				return nil, err
			}
			cur.Args = a
			cur.Return = r
			i += n

			n, block, err := consumeBlock(i, tokens, c)
			if err != nil {
				return nil, err
			}
			cur.Body = block

			i += n

			nodes = append(nodes, cur)
		}
	}
	return nodes, nil
}

func consumePrototype(start int, tokens []token.Token, c Context) (n int, args []VarWithType, retn []VarWithType, err error) {
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

func extractPrototypes(tokens []token.Token, c Context) error {
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

		}
	}
	return nil
}

// consumeBlock consumes a balanced number of tokens delimited by a balanced
// number of "{" and "}" characters. It returns the ASTNode for the block, and
// the number of tokens that were consumed.
func consumeBlock(start int, tokens []token.Token, c Context) (int, BlockStmt, error) {
	if tokens[start] != token.Char("{") {
		return 0, BlockStmt{}, fmt.Errorf("Invalid block. (%v)", tokens[start])
	}
	var blockStmt BlockStmt

	for i := start + 1; i < len(tokens); i++ {
		// First handle open or close brackets
		if tokens[i] == token.Char("{") {
			c2 := c.Clone()
			n, subblock, err := consumeBlock(i, tokens, c2)
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
					return 0, BlockStmt{}, fmt.Errorf("Can not mutate immutable variable %v", tokens[i])
				}
				n, val, err := consumeValue(i+2, tokens, c)
				if err != nil {
					return 0, BlockStmt{}, err
				}
				blockStmt.Stmts = append(blockStmt.Stmts, AssignmentOperator{
					Variable: Variable(tokens[i].String()),
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
			default:
				panic(fmt.Sprintf("Unimplemented keyword: %v at %v", tokens[i], i))
			}
		default:
			panic(fmt.Sprintf("Unhandled token type in block %v for token %v", tokens[i].String(), i))
		}

	}
	return 0, BlockStmt{}, fmt.Errorf("Unterminated block statement")
}

func consumeFuncCall(start int, tokens []token.Token, c Context) (int, FuncCall, error) {
	name := tokens[start].String()
	f := FuncCall{
		Name: name,
	}

	decl, ok := c.Functions[name]
	if !ok {
		return 0, FuncCall{}, fmt.Errorf("Undefined function: %v", name)
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
						f.Args = append(f.Args, strLit)
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
					f.Args = append(f.Args, Variable(t.String()))
				} else if c.IsFunction(t.String()) {
					n, subcall, err := consumeFuncCall(i, tokens, c)
					if err != nil {
						return 0, FuncCall{}, err
					}
					f.Args = append(f.Args, subcall)
					i += n
				} else {
					panic("Unhandled token" + t.String())
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
	if len(args) == 0 {
		if tokens[start+1] == token.Char("(") && tokens[start+2] == token.Char(")") {
			return 2, f, nil
		}
		return 0, FuncCall{}, fmt.Errorf("Unexpected parameter to %v", name)
	}
	argStart := start + 2 // start = name, +1 = "(", +2 = the first param..
	for i := 0; i < len(args); i++ {
		n, val, err := consumeValue(argStart, tokens, c)
		if err != nil {
			return 0, FuncCall{}, err
		}
		argStart += n
		if i != len(args)-1 {
			if tokens[argStart] != token.Char(",") {
				return 0, FuncCall{}, fmt.Errorf("Invalid call to %v: missing comma between parameters at token %v (arg %d of %d)", name, argStart, i, len(args))
			}
		} else if tokens[argStart] != token.Char(")") {
			return 0, FuncCall{}, fmt.Errorf("Invalid call to %v: missing closing bracket at token %v", name, argStart)
		}
		argStart++
		f.Args = append(f.Args, val)
	}
	return argStart - start - 1, f, nil
	//return 0, FuncCall{}, fmt.Errorf("Invalid function call to %v. Missing closing bracket", tokens[start])
}

func consumeLetStmt(start int, tokens []token.Token, c Context) (int, Node, error) {
	l := LetStmt{}

	if tokens[start] != token.Keyword("let") {
		return 0, nil, fmt.Errorf("Invalid let statement")
	}
	for i := start + 1; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			if l.Var.Name == "" {
				l.Var.Name = Variable(t.String())
				c.Variables[t.String()] = true
			} else if l.Var.Type == "" {
				l.Var.Type = t.String()
			} else {
				return 0, nil, fmt.Errorf("Invalid name for let statement")
			}
		case token.Type:
			if l.Var.Type == "" {
				l.Var.Type = t.String()
			} else {
				return 0, nil, fmt.Errorf("Unexpected type in let statement")

			}
		case token.Operator:
			if t == token.Operator("=") {
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
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

func consumeMutStmt(start int, tokens []token.Token, c Context) (int, Node, error) {
	l := MutStmt{}

	if tokens[start] != token.Keyword("mut") {
		return 0, nil, fmt.Errorf("Invalid mutable variable statement")
	}
	for i := start + 1; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case token.Unknown:
			if l.Var.Name == "" {
				l.Var.Name = Variable(t.String())
				c.Variables[t.String()] = true
				c.Mutables[t.String()] = true
			} else if l.Var.Type == "" {
				l.Var.Type = t.String()
			} else {
				return 0, nil, fmt.Errorf("Invalid name for mut statement")
			}
		case token.Type:
			if l.Var.Type == "" {
				l.Var.Type = t.String()
			} else {
				return 0, nil, fmt.Errorf("Unexpected type in mut statement")

			}
		case token.Operator:
			if t == token.Operator("=") {
				if l.Var.Type == "" {
					return 0, nil, fmt.Errorf("Type inference not yet implemented")
				}
				n, v, err := consumeValue(i+1, tokens, c)
				if err != nil {
					return 0, nil, err
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

func consumeArgs(start int, tokens []token.Token, c Context) (int, []VarWithType, error) {
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
					Type: t.String(),
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
				Type: t.String(),
			})
			continue
		default:
			return 0, nil, fmt.Errorf("Invalid token in argument list: %v", t)
		}
	}
	return 0, nil, fmt.Errorf("Could not parse arguments")
}
