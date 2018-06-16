package token

import (
	"github.com/driusan/lang/parser/sampleprograms"
	"strings"
	"testing"
)

func TestParseFizzbuzz(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.Fizzbuzz))
	expected := []Token{
		Keyword("func"), // 0
		Whitespace(" "),
		Unknown("main"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Operator("->"),
		Whitespace(" "),
		Keyword("affects"),
		Char("("),
		Unknown("IO"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"), // 10
		Keyword("mutable"),
		Whitespace(" "),
		Unknown("terminate"),
		Whitespace(" "),
		Type("bool"),
		Whitespace(" "),
		Operator("="),
		Whitespace(" "),
		Unknown("false"),
		Whitespace("\n\t"), // 20
		Keyword("mutable"),
		Whitespace(" "),
		Unknown("i"),
		Whitespace(" "),
		Type("int"),
		Whitespace(" "),
		Operator("="),
		Whitespace(" "),
		Unknown("1"),
		Whitespace("\n\t"), // 30
		Keyword("while"),
		Whitespace(" "),
		Unknown("terminate"),
		Whitespace(" "),
		Operator("!="),
		Whitespace(" "),
		Unknown("true"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t\t"),
		Keyword("if"),
		Whitespace(" "),
		Unknown("i"),
		Whitespace(" "), // 40
		Operator("%"),
		Whitespace(" "),
		Unknown("15"),
		Whitespace(" "),
		Operator("=="),
		Whitespace(" "),
		Unknown("0"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t\t\t"),
		Unknown("PrintString"),
		Char(`(`),
		Char(`"`),
		String(`fizzbuzz`), // 50
		Char(`"`),
		Char(`)`),
		Whitespace("\n\t\t"),
		Char("}"),
		Whitespace(" "),
		Keyword("else"),
		Whitespace(" "),
		Keyword("if"),
		Whitespace(" "),
		Unknown("i"),
		Whitespace(" "), // 60
		Operator("%"),
		Whitespace(" "),
		Unknown("5"),
		Whitespace(" "),
		Operator("=="),
		Whitespace(" "),
		Unknown("0"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t\t\t"),
		Unknown("PrintString"),
		Char(`(`),
		Char(`"`),
		String(`buzz`), // 70
		Char(`"`),
		Char(`)`),
		Whitespace("\n\t\t"),
		Char("}"),
		Whitespace(" "),
		Keyword("else"),
		Whitespace(" "),
		Keyword("if"),
		Whitespace(" "),
		Unknown("i"), // 80
		Whitespace(" "),
		Operator("%"),
		Whitespace(" "),
		Unknown("3"),
		Whitespace(" "),
		Operator("=="),
		Whitespace(" "),
		Unknown("0"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t\t\t"),
		Unknown("PrintString"),
		Char(`(`),
		Char(`"`), // 90
		String(`fizz`),
		Char(`"`),
		Char(`)`),
		Whitespace("\n\t\t"),
		Char("}"),
		Whitespace(" "),
		Keyword("else"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t\t\t"), // 100
		Unknown("PrintInt"),
		Char(`(`),
		Unknown("i"),
		Char(`)`),
		Whitespace("\n\t\t"),
		Char("}"),
		Whitespace("\n\t\t"),
		Unknown("PrintString"),
		Char(`(`),
		Char(`"`),
		String(`\n`), // 50
		Char(`"`),
		Char(`)`),
		Whitespace("\n\n\t\t"),

		Unknown("i"),
		Whitespace(" "),
		Operator("="), // 110
		Whitespace(" "),
		Unknown("i"),
		Whitespace(" "),
		Operator("+"),
		Whitespace(" "),
		Unknown("1"),
		Whitespace("\n\t\t"),
		Keyword("if"),
		Whitespace(" "),
		Unknown("i"), // 120
		Whitespace(" "),
		Operator(">="),
		Whitespace(" "),
		Unknown("100"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t\t\t"),
		Unknown("terminate"),
		Whitespace(" "),
		Operator("="), // 130
		Whitespace(" "),
		Unknown("true"),
		Whitespace("\n\t\t"),
		Char("}"),
		Whitespace("\n\t"),
		Char("}"),
		Whitespace("\n"),
		Char("}"),
	}
	if err != nil {
		t.Fatal(err)
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}

}

func TestHelloWorld(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.HelloWorld))
	expected := []Token{
		Keyword("func"), // 0
		Whitespace(" "),
		Unknown("main"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Operator("->"),
		Whitespace(" "),
		Keyword("affects"),
		Char("("),
		Unknown("IO"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"), // 10
		Unknown("PrintString"),
		Char("("),
		Char(`"`),
		String(`Hello, world!\n`),
		Char(`"`),
		Char(")"),
		Whitespace("\n"),
		Char("}"),
	}
	if err != nil {
		t.Fatal(err)
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}

}

func TestTwoProcs(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.TwoProcs))
	expected := []Token{
		Keyword("func"),
		Whitespace(" "),
		Unknown("foo"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Type("int"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Keyword("return"),
		Whitespace(" "),
		Unknown("3"),
		Whitespace("\n"),
		Char("}"),
		Whitespace("\n\n"),
		Keyword("func"),
		Whitespace(" "),
		Unknown("main"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Operator("->"),
		Whitespace(" "),
		Keyword("affects"),
		Char("("),
		Unknown("IO"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Unknown("PrintInt"),
		Char("("),
		Unknown("foo"),
		Char("("),
		Char(")"),
		Char(")"),
		Whitespace("\n"),
		Char("}"),
	}
	if err != nil {
		t.Fatal(err)
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}

}

func TestMatchParam2(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.MatchParam2))
	expected := []Token{
		Keyword("enum"),
		Whitespace(" "),
		Unknown("Maybe"),
		Whitespace(" "),
		Unknown("x"),
		Whitespace(" "),
		Operator("="),
		Whitespace(" "),
		Unknown("Nothing"),
		Whitespace(" "),
		Operator("|"),
		Whitespace(" "),
		Unknown("Just"),
		Whitespace(" "),
		Unknown("x"),
		Whitespace("\n\n"),
		Keyword("func"),
		Whitespace(" "),
		Unknown("foo"),
		Whitespace(" "),
		Char("("),
		Unknown("x"),
		Whitespace(" "),
		Unknown("Maybe"),
		Whitespace(" "),
		Type("int"),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Type("int"),
		Char(")"),
		Whitespace(" "),
		Operator("->"),
		Whitespace(" "),
		Keyword("affects"),
		Char("("),
		Unknown("IO"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Unknown("PrintString"),
		Char("("),
		Char(`"`),
		String("x"),
		Char(`"`),
		Char(")"),
		Whitespace("\n\t"),
		Keyword("match"),
		Whitespace(" "),
		Unknown("x"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Keyword("case"),
		Whitespace(" "),
		Unknown("Just"),
		Whitespace(" "),
		Unknown("n"),
		Char(":"),
		Whitespace("\n\t\t"),
		Keyword("return"),
		Whitespace(" "),
		Unknown("n"),
		Whitespace("\n\t"),
		Keyword("case"),
		Whitespace(" "),
		Unknown("Nothing"),
		Char(":"),
		Whitespace("\n\t\t"),
		Keyword("return"),
		Whitespace(" "),
		Unknown("0"),
		Whitespace("\n\t"),
		Char("}"),
		Whitespace("\n"),
		Char("}"),
		Whitespace("\n\n"),
		Keyword("func"),
		Whitespace(" "),
		Unknown("main"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Operator("->"),
		Whitespace(" "),
		Keyword("affects"),
		Char("("),
		Unknown("IO"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Unknown("PrintInt"),
		Char("("),
		Unknown("foo"),
		Char("("),
		Unknown("Just"),
		Whitespace(" "),
		Unknown("5"),
		Char(")"),
		Char(")"),
		Whitespace("\n"),
		Char("}"),
	}
	if err != nil {
		t.Fatal(err)
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}

}
func TestNoWhitespace(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(`1+2*3`))
	if err != nil {
		t.Fatal(err)
	}
	expected := []Token{
		Unknown("1"),
		Operator("+"),
		Unknown("2"),
		Operator("*"),
		Unknown("3"),
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
}

func TestSingleCharInput(t *testing.T) {
	tk, err := Tokenize(strings.NewReader("%"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []Token{
		Operator("%"),
	}

	if len(tk) != 1 {
		t.Fatalf("Unexpected number of tokens. Got: %v", tk)
	}
	for i, tok := range expected {
		if tok != tk[i] {
			t.Errorf("Unexpected token: got %v want %v", tk[i], expected[i])
		}
	}
}

func TestSimpleArray(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.SimpleArray))
	expected := []Token{
		Keyword("func"),
		Whitespace(" "),
		Unknown("main"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Operator("->"),
		Whitespace(" "),
		Keyword("affects"),
		Char("("),
		Unknown("IO"),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Keyword("let"),
		Whitespace(" "),
		Unknown("n"),
		Whitespace(" "),
		Char("["),
		Unknown("5"),
		Char("]"),
		Type("int"),
		Whitespace(" "),
		Operator("="),
		Whitespace(" "),
		Char("{"),
		Whitespace(" "),
		Unknown("1"),
		Char(","),
		Whitespace(" "),
		Unknown("2"),
		Char(","),
		Whitespace(" "),
		Unknown("3"),
		Char(","),
		Whitespace(" "),
		Unknown("4"),
		Char(","),
		Whitespace(" "),
		Unknown("5"),
		Whitespace(" "),
		Char("}"),
		Whitespace("\n\t"),
		Unknown("PrintInt"),
		Char("("),
		Unknown("n"),
		Char("["),
		Unknown("3"),
		Char("]"),
		Char(")"),
		Whitespace("\n"),
		Char("}"),
	}
	if err != nil {
		t.Fatal(err)
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}

}

func TestAssertion(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.AssertionFail))
	expected := []Token{
		Keyword("func"),
		Whitespace(" "),
		Unknown("main"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"),
		Keyword("assert"),
		Char("("),
		Unknown("false"),
		Char(")"),
		Whitespace("\n"),
		Char("}"),
	}
	if err != nil {
		t.Fatal(err)
	}

	// We ensure that both there are no tokens that
	// were parsed that shouldn't have been, and no tokens
	// that weren't parsed by ranging through both what
	// we got and what we expected separately.
	for i, tok := range tokens {
		if i >= len(expected) {
			t.Errorf("Not enough tokens parsed")
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}
	for i, tok := range expected {
		if i >= len(tokens) {
			t.Errorf("Missing %d tokens", (i - len(tokens) + 1))
			break
		}
		if !tok.IsValid() {
			t.Errorf("Invalid token %v", tok)
		}
		if tokens[i] != expected[i] {
			t.Errorf(`Unexpected token %d: got "%v" want "%v"`, i, tokens[i], expected[i])
		}
	}

}
