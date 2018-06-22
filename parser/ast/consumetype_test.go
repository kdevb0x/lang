package ast

import (
	"strings"
	"testing"

	"github.com/driusan/lang/parser/token"
)

func TestConsumeType(t *testing.T) {
	cases := []struct {
		Code      string
		Expected  Type
		ExpectedN int
	}{

		{"int", TypeLiteral("int"), 1},
		// Needs context:{"Maybe int", TypeLiteral("Maybe int"), 2},
		{"[5]int", ArrayType{Base: TypeLiteral("int"), Size: IntLiteral(5)}, 4},
		{"int | string", SumType{TypeLiteral("int"), TypeLiteral("string")}, 3},
	}

	for i, tc := range cases {
		tokens, err := token.Tokenize(strings.NewReader(tc.Code))
		tokens = stripWhitespaceAndComments(tokens)

		c := NewContext()
		n, value, err := consumeType(0, tokens, &c)
		if err != nil {
			t.Fatal(err)
		}

		if value == nil {
			t.Errorf("Case %v: got nil want %v", i, tc.Expected)
			continue
		}

		if !compare(value, tc.Expected) {
			t.Errorf("Case %v: got %v want %v", i, value, tc.Expected)
		}

		if n != tc.ExpectedN {
			t.Errorf("Case %v: Consumed %v tokens, expected %v", i, n, tc.ExpectedN)
		}
	}
}
