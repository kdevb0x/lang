package ast

import (
	"fmt"
	"strings"

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
		cn, cv, err = consumeValue(start+1, tokens, c)
		if err != nil {
			return 0, MatchStmt{}, err
		}
	}
	l.Condition = cv

	if tokens[start+cn+1] != token.Char("{") {
		return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
	}

	concreteMap := make(map[Type]Type)
	for i := start + cn + 2; i < len(tokens); {
		c2 := c.Clone()

		if concretes := strings.Fields(string(l.Condition.Type())); len(concretes) > 1 {
			// Convert the generic type's parameters to the concrete typename,
			// not the generic type, so that the case can look them up properly.
			// FIXME: There should be a better way than splitting the type name
			// on whitespace
			td, ok := c.Types[concretes[0]]
			if !ok {
				panic("Expected parameterized enumerated option")
			}
			for i, v := range concretes[1:] {
				concreteMap[td.Parameters[i]] = TypeLiteral(v)
			}
		}

		n, cs, err := consumeCase(i, tokens, &c2, concreteMap)
		if err != nil {
			return 0, MatchStmt{}, err
		}
		l.Cases = append(l.Cases, cs)
		i += n
		if tokens[i] == token.Char("}") {
			ct := c.Types[l.Condition.Type()].ConcreteType
			if ct != nil && c.Types[l.Condition.Type()].ConcreteType.Type() == "sumtype" {
				if err := checkExhaustiveness(l.Condition, l.Cases, c); err != nil {
					return 0, MatchStmt{}, err
				}
			}
			return i - start, l, nil
		}
	}
	return 0, MatchStmt{}, fmt.Errorf("Invalid match statement")
}

func consumeCase(start int, tokens []token.Token, c *Context, genericMap map[Type]Type) (int, MatchCase, error) {
	l := MatchCase{}
	var n int
	if tokens[start] != token.Keyword("case") {
		return 0, MatchCase{}, fmt.Errorf("Invalid case statement. Unexpected '%v' at %d", tokens[start], start)
	}
	if eo := c.EnumeratedOption(tokens[start+1].String()); eo != nil {
		n = 1
		for _, t := range eo.Parameters {
			varname := tokens[start+1+n].String()
			c.Variables[varname] = VarWithType{
				Name: Variable(varname),
				Typ:  genericMap[t],
			}
			l.LocalVariables = append(l.LocalVariables, c.Variables[varname])
			n += 1
		}
		l.Variable = *eo
	} else {
		n2, v, err := consumeValue(start+1, tokens, c)
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

func checkExhaustiveness(t Type, mc []MatchCase, c *Context) error {
	allcases := make(map[string]bool)
	for _, eo := range c.EnumOptions {
		if eo.Type() == t.Type() {
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
			return fmt.Errorf(`Inexhaustive match for enum type "%v": Missing case "%v".`, t.Type(), c)
		}
	}
	return nil
}
