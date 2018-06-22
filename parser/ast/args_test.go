package ast

import (
	"github.com/driusan/lang/parser/token"
	"strings"
	"testing"
)

func TestConsumeTypeTuple(t *testing.T) {
	tests := []struct {
		Val       string
		Expected  []VarWithType
		ExpectedN int
	}{
		{"()", nil, 2},
		{"(n int)", []VarWithType{{Name: "n", Typ: TypeLiteral("int")}}, 4},
		{
			"(partial int, x int)",
			[]VarWithType{
				{Name: "partial", Typ: TypeLiteral("int")},
				{Name: "x", Typ: TypeLiteral("int")},
			},
			7,
		},
	}

	c := NewContext()
	for _, tc := range tests {
		tokens, err := token.Tokenize(strings.NewReader(tc.Val))
		if err != nil {
			t.Fatal(err)
		}

		tokens = stripWhitespaceAndComments(tokens)
		n, vt, err := consumeTupleType(0, tokens, &c)
		if n != tc.ExpectedN {
			t.Errorf("Unexpected number of values returned: got %v want %v", n, tc.ExpectedN)
		}
		if len(vt) != len(tc.Expected) {
			t.Fatalf("Unexpected number of args: got %v want %v", len(vt), len(tc.Expected))
		}
		for i := range vt {
			if vt[i] != tc.Expected[i] {
				t.Errorf("Unexpected result: got %v want %v", vt[i], tc.Expected[i])
			}
		}

	}
}
