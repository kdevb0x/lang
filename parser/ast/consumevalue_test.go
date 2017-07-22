package ast

import (
	"strings"
	"testing"

	"github.com/driusan/lang/parser/token"
)

func TestConsumeValue(t *testing.T) {
	cases := []struct {
		Code      string
		Expected  Node
		ExpectedN int
	}{

		{"1", IntLiteral(1), 1},
		{"1 + 2", AdditionOperator{Left: IntLiteral(1), Right: IntLiteral(2)}, 3},
		{"1 - 2", SubtractionOperator{Left: IntLiteral(1), Right: IntLiteral(2)}, 3},
		{"1 * 2", MulOperator{Left: IntLiteral(1), Right: IntLiteral(2)}, 3},
		{"1 / 2", DivOperator{Left: IntLiteral(1), Right: IntLiteral(2)}, 3},
		{
			"1 - 2 * 4",
			SubtractionOperator{
				Left: IntLiteral(1),
				Right: MulOperator{
					Left:  IntLiteral(2),
					Right: IntLiteral(4),
				},
			},
			5,
		},
		{
			"1 * 2 - 4",
			SubtractionOperator{
				Left: MulOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: IntLiteral(4),
			},
			5,
		},
		{
			"4 - 4 / 2",
			SubtractionOperator{
				Left: IntLiteral(4),
				Right: DivOperator{
					Left:  IntLiteral(4),
					Right: IntLiteral(2),
				},
			},
			5,
		},
		{
			"2 * 4 - 6 / 3",
			SubtractionOperator{
				Left: MulOperator{
					Left:  IntLiteral(2),
					Right: IntLiteral(4),
				},
				Right: DivOperator{
					Left:  IntLiteral(6),
					Right: IntLiteral(3),
				},
			},
			7,
		},
		{
			"1 - 2 * 4 - 4 / 2",
			SubtractionOperator{
				Left: IntLiteral(1),
				Right: SubtractionOperator{
					Left:  MulOperator{IntLiteral(2), IntLiteral(4)},
					Right: DivOperator{IntLiteral(4), IntLiteral(2)},
				},
			},
			9,
		},
		{
			"1 + 2 * 4 - 4 / 2",
			AdditionOperator{
				Left: IntLiteral(1),
				Right: SubtractionOperator{
					Left:  MulOperator{IntLiteral(2), IntLiteral(4)},
					Right: DivOperator{IntLiteral(4), IntLiteral(2)},
				},
			},
			9,
		},
		{
			"1 - 2 + 3 / 4",
			SubtractionOperator{
				Left: IntLiteral(1),
				Right: AdditionOperator{
					Left: IntLiteral(2),
					Right: DivOperator{
						Left:  IntLiteral(3),
						Right: IntLiteral(4),
					},
				},
			},
			7,
		},
		{
			"1 / 2 + 3 - 4",
			SubtractionOperator{
				Left: AdditionOperator{
					Left: DivOperator{
						Left:  IntLiteral(1),
						Right: IntLiteral(2),
					},
					Right: IntLiteral(3),
				},
				Right: IntLiteral(4),
			},
			7,
		},
		{
			"1 - 2 + 3 - 4 / 5 + 6 - 7 + 8",
			SubtractionOperator{
				Left: IntLiteral(1),
				Right: AdditionOperator{
					Left: IntLiteral(2),
					Right: SubtractionOperator{
						Left: IntLiteral(3),
						Right: SubtractionOperator{
							Left: AdditionOperator{
								Left: DivOperator{
									Left:  IntLiteral(4),
									Right: IntLiteral(5),
								},
								Right: IntLiteral(6),
							},
							Right: AdditionOperator{
								IntLiteral(7),
								IntLiteral(8),
							},
						},
					},
				},
			},
			15,
		},
		{
			"1 + 2 != 3 + 4",
			NotEqualsComparison{
				Left: AdditionOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: AdditionOperator{
					Left:  IntLiteral(3),
					Right: IntLiteral(4),
				},
			},
			7,
		},
		{
			"1 + 2 == 3 + 4",
			EqualityComparison{
				Left: AdditionOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: AdditionOperator{
					Left:  IntLiteral(3),
					Right: IntLiteral(4),
				},
			},
			7,
		},
		{
			"1 + 2 >= 3 + 4",
			GreaterOrEqualComparison{
				Left: AdditionOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: AdditionOperator{
					Left:  IntLiteral(3),
					Right: IntLiteral(4),
				},
			},
			7,
		},
		{
			"1 + 2 > 3 + 4",
			GreaterComparison{
				Left: AdditionOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: AdditionOperator{
					Left:  IntLiteral(3),
					Right: IntLiteral(4),
				},
			},
			7,
		},
		{
			"1 + 2 < 3 + 4",
			LessThanComparison{
				Left: AdditionOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: AdditionOperator{
					Left:  IntLiteral(3),
					Right: IntLiteral(4),
				},
			},
			7,
		},
		{
			"1 + 2 <= 3 + 4",
			LessThanOrEqualComparison{
				Left: AdditionOperator{
					Left:  IntLiteral(1),
					Right: IntLiteral(2),
				},
				Right: AdditionOperator{
					Left:  IntLiteral(3),
					Right: IntLiteral(4),
				},
			},
			7,
		},
		{
			"1 + 2*3",
			AdditionOperator{
				Left: IntLiteral(1),
				Right: MulOperator{
					Left:  IntLiteral(2),
					Right: IntLiteral(3),
				},
			},
			5,
		},
	}

	for i, tc := range cases {
		tokens, err := token.Tokenize(strings.NewReader(tc.Code))
		tokens = stripWhitespace(tokens)

		n, value, err := consumeValue(0, tokens, &Context{})
		if err != nil {
			t.Fatal(err)
		}

		if !compare(value, tc.Expected) {
			t.Errorf("Case %v: got %v want %v", i, value, tc.Expected)
		}

		if n != tc.ExpectedN {
			t.Errorf("Case %v: Consumed %v tokens, expected %v", i, n, tc.ExpectedN)
		}
	}
}
