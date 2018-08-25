package gobackend

import (
	"testing"
	//	"fmt"
	"bufio"
	"os"
	"strings"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/token"
)

func parseTestCase(t *testing.T, filename string) []ast.Node {
	t.Helper()

	f, err := os.Open("../../testsuite/" + filename + ".l")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Add a RuneReader on top of f.
	tokens, err := token.Tokenize(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := ast.Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	return ast
}
func TestAssertions(t *testing.T) {
	tests := []struct {
		filename, expected string
	}{
		{
			"AssertionFail",
			`func main() {
	if !(false) {
		panic("Assertion failed")
	}
}
`,
		},
		{
			"AssertionFailWithMessage",
			`func main() {
	if !(false) {
		panic("This always fails")
	}
}
`,
		},
		{
			"AssertionPass",
			`func main() {
	if !(true) {
		panic("Assertion failed")
	}
}
`,
		},
		{
			"AssertionPassWithMessage",
			`func main() {
	if !(true) {
		panic("You should never see this")
	}
}
`,
		},
		{
			"AssertionFailWithContext",
			`func main() {
	fmt.Printf("%d",
		0,
	)
	if !(false) {
		panic("Assertion failed")
	}
	fmt.Printf("%d",
		1,
	)
}
`,
		},
		{
			"AssertionFailWithVariable",
			`func main() {
	x := 3
	if !(x > 3) {
		panic("Assertion failed")
	}
}
`,
		},
	}

	for _, tc := range tests {
		as := parseTestCase(t, tc.filename)
		var result strings.Builder
		if _, err := Convert(&result, as[0], nil); err != nil {
			t.Errorf("%v: %v", tc.filename, err)
			continue
		}
		if result.String() != tc.expected {
			t.Errorf("%v: Unexpected convertion: got %v want %v", tc.filename, result.String(), tc.expected)
			continue
		}
	}
}
