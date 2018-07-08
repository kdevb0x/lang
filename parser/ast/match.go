package ast

import (
	"fmt"
	"reflect"

	"github.com/driusan/lang/parser/token"
)

type MatchCase struct {
	LocalVariables []VarWithType
	Variable       Value
	Body           BlockStmt
}

func (i MatchCase) Node() Node {
	return i
}

func (i MatchCase) String() string {
	return fmt.Sprintf("MatchCase{Locals: %v\n\tVariable: %v,\n\tBody: %v}", i.LocalVariables, i.Variable, i.Body)

}

func (i MatchCase) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

type MatchStmt struct {
	Condition Value
	Cases     []MatchCase
}

func (i MatchStmt) String() string {
	return fmt.Sprintf("MatchStmt{Condition: %v,\n\tBody: %v}", i.Condition, i.Cases)
}

func (i MatchStmt) Node() Node {
	return i
}

func (i MatchStmt) PrettyPrint(lvl int) string {
	panic("Not implemented")
}

func consumeMatchStmt(start int, tokens []token.Token, c *Context) (int, MatchStmt, error) {
	l := MatchStmt{}

	if tokens[start] != token.Keyword("match") {
		return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
	}
	var cn int
	var cv Value
	var err error
	if tokens[start+1] == token.Char("{") {
		cn = 0
		cv = BoolLiteral(true)
	} else {
		cn, cv, err = consumeValue(start+1, tokens, c, true)
		if err != nil {
			return 0, MatchStmt{}, err
		}
	}
	l.Condition = cv

	if tokens[start+cn+1] != token.Char("{") {
		return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
	}

	for i := start + cn + 2; i < len(tokens); {
		c2 := c.Clone()
		switch t := l.Condition.Type().(type) {
		case SumType:
			n, cs, err := consumeTypeCase(i, tokens, &c2, l.Condition)
			if err != nil {
				return 0, MatchStmt{}, err
			}
			l.Cases = append(l.Cases, cs)
			i += n
		case EnumTypeDefn:
			n, cs, err := consumeCase(i, tokens, &c2, t.Parameters)
			if err != nil {
				return 0, MatchStmt{}, err
			}
			l.Cases = append(l.Cases, cs)
			i += n
		default:
			n, cs, err := consumeCase(i, tokens, &c2, nil)
			if err != nil {
				return 0, MatchStmt{}, err
			}
			l.Cases = append(l.Cases, cs)
			i += n
		}
		if tokens[i] == token.Char("}") {
			ct := c.Types[l.Condition.Type().TypeName()].ConcreteType
			if _, ok := ct.(EnumTypeDefn); ok {
				if err := checkExhaustiveness(l.Condition.Type(), l.Cases, c); err != nil {
					return 0, MatchStmt{}, err
				}
			}
			return i + 1 - start, l, nil
		}
	}
	return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
}

func consumeCase(start int, tokens []token.Token, c *Context, subtypes []Type) (int, MatchCase, error) {
	l := MatchCase{}
	var n int
	if tokens[start] != token.Keyword("case") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement. Unexpected '%v' at %d", tokens[start], start)
	}
	if eo := c.EnumeratedOption(tokens[start+1].String()); eo != nil {
		n = 1
		for i, p := range eo.Parameters {
			if p == "unused" {
				break
			}
			varname := tokens[start+1+n].String()
			c.Variables[varname] = VarWithType{
				Name: Variable(varname),
				Typ:  subtypes[i],
			}
			l.LocalVariables = append(l.LocalVariables, c.Variables[varname])
			n += 1
		}
		l.Variable = *eo
	} else {
		n2, v, err := consumeValue(start+1, tokens, c, false)
		if err != nil {
			return 0, MatchCase{}, err
		}
		l.Variable = v
		n = n2

	}
	if tokens[start+n+1] != token.Char(":") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement at token %v. Expected ':', not '%v'", start, tokens[start+n+1])
	}
	for i := start + n + 2; i < len(tokens); {
		if tokens[i] == token.Keyword("case") || tokens[i] == token.Char("}") {
			return i - start, l, nil
		}
		n, stmt, err := consumeStmt(i, tokens, c)
		if err != nil {
			return 0, MatchCase{}, err
		}
		l.Body.Stmts = append(l.Body.Stmts, stmt)
		i += n
	}
	return 0, MatchCase{}, fmt.Errorf("Unterminated case statement")
}

func consumeTypeCase(start int, tokens []token.Token, c *Context, condition Value) (int, MatchCase, error) {
	l := MatchCase{}
	var n int
	if tokens[start] != token.Keyword("case") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement. Unexpected '%v' at %d", tokens[start], start)
	}
	n2, t, err := consumeType(start+1, tokens, c)
	if err != nil {
		return 0, MatchCase{}, err
	}
	n = n2

	switch v := condition.(type) {
	case VarWithType:
		v.Typ = t
		c.Variables[string(v.Name)] = v
		l.Variable = v
	default:
		panic(fmt.Sprintf("Can only destructure single variable sum types, got %v", reflect.TypeOf(v)))
	}
	if tokens[start+n+1] != token.Char(":") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement at token %v. Expected ':', not '%v'", start, tokens[start+n+1])
	}
	for i := start + n + 2; i < len(tokens); {
		if tokens[i] == token.Keyword("case") || tokens[i] == token.Char("}") {
			return i - start, l, nil
		}
		n, stmt, err := consumeStmt(i, tokens, c)
		if err != nil {
			return 0, MatchCase{}, err
		}
		l.Body.Stmts = append(l.Body.Stmts, stmt)
		i += n
	}
	return 0, MatchCase{}, fmt.Errorf("Unterminated case statement")
}

func checkExhaustiveness(t Type, mc []MatchCase, c *Context) error {
	allcases := make(map[string]bool)
	for _, eo := range c.EnumOptions {
		if eo.Type().TypeName() == t.TypeName() {
			allcases[eo.Constructor] = false
		}
	}

	for _, m := range mc {
		if eo, ok := m.Variable.(EnumOption); ok {
			allcases[eo.Constructor] = true
		}
	}
	for c, v := range allcases {
		if v == false {
			return fmt.Errorf(`Inexhaustive match for enum type "%v": Missing case "%v".`, t.TypeName(), c)
		}
	}
	return nil
}
