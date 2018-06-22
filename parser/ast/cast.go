package ast

import (
	"fmt"

	"github.com/driusan/lang/parser/token"
)

func consumeCastStmt(start int, tokens []token.Token, c *Context) (int, Cast, error) {
	if tokens[start+1] != token.Char("(") {
		return 0, Cast{}, fmt.Errorf("Invalid cast (missing value)")

	}
	// +2 to skip past the bracket
	vn, v, err := consumeValue(start+2, tokens, c, false)
	if err != nil {
		return 0, Cast{}, err
	}

	if tokens[start+vn+2] != token.Char(")") {
		return 0, Cast{}, fmt.Errorf("Missing closing bracket for cast value")
	}
	// +3 for the extra ")" that wasn't included in the consumeValue, since
	// it skipped past the bracket
	if tokens[start+vn+3] != token.Keyword("as") {
		return 0, Cast{}, fmt.Errorf("Invalid cast (missing \"as\")")
	}

	// 4 includes another 1 for "as"
	tn, t, err := consumeType(start+vn+4, tokens, c)
	if err != nil {
		return 0, Cast{}, err
	}

	return tn + vn + 4, Cast{Val: v, Typ: t}, nil
}
