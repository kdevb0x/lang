package ast

import (
	"strings"
	"testing"

	"github.com/driusan/lang/parser/token"
)

func TestConsumeSumTypeList(t *testing.T) {
	cases := []struct {
		Code      string
		Expected  []EnumOption
		ExpectedN int
	}{

		{"A", []EnumOption{{Constructor: "A"}}, 1},
		{
			"A | B",
			[]EnumOption{
				{Constructor: "A"},
				{Constructor: "B"},
			},
			3,
		},
		{"Just a | B",
			[]EnumOption{
				{Constructor: "Just", Parameters: []Type{TypeLiteral("a")}},
				{Constructor: "B"},
			},
			4,
		},
	}

	for i, tc := range cases {
		tokens, err := token.Tokenize(strings.NewReader(tc.Code))
		tokens = stripWhitespace(tokens)

		n, value, err := consumeSumTypeList(0, tokens, &Context{})
		if err != nil {
			t.Fatal(err)
		}
		if n != tc.ExpectedN {
			t.Errorf("Case %v: Consumed %v tokens, expected %v", i, n, tc.ExpectedN)
			continue
		}

		for j := range tc.Expected {
			if !compare(tc.Expected[j], value[j]) {
				t.Errorf("Case %v[%v]: got %v want %v", i, j, value[j], tc.Expected[j])
			}
		}

	}
}
