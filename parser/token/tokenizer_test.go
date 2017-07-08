package token

import (
	"github.com/driusan/lang/parser/sampleprograms"
	"strings"
	"testing"
)

func TestParseFizzbuzz(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.Fizzbuzz))
	expected := []Token{
		Keyword("proc"), // 0
		Whitespace(" "),
		Unknown("main"),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"), // 10
		Keyword("mut"),
		Whitespace(" "),
		Unknown("terminate"),
		Whitespace(" "),
		Type("bool"),
		Whitespace(" "),
		Operator("="),
		Whitespace(" "),
		Unknown("false"),
		Whitespace("\n\t"), // 20
		Keyword("mut"),
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
		Keyword("proc"), // 0
		Whitespace(" "),
		Unknown("main"),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
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

/*
func TestHelloWorld2(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.HelloWorld2))
	expected := []Token{
		Keyword("proc"), // 0
		Whitespace(" "),
		Unknown("main"),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("{"),
		Whitespace("\n\t"), // 10
		Unknown("PrintString"),
		Char("("),
		Char(`"`),
		String(`%s %s\n %s`),
		Char(`"`),
		Char(","),
		Whitespace(" "),
		Char(`"`),
		String(`Hello, world!\n`),
		Char(`"`), // 20
		Char(","),
		Whitespace(" "),
		Char(`"`),
		String(`World??`),
		Char(`"`),
		Char(","),
		Whitespace(" "),
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
*/
func TestTwoProcs(t *testing.T) {
	tokens, err := Tokenize(strings.NewReader(sampleprograms.TwoProcs))
	expected := []Token{
		Keyword("proc"),
		Whitespace(" "),
		Unknown("foo"),
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
		Keyword("proc"),
		Whitespace(" "),
		Unknown("main"),
		Char("("),
		Char(")"),
		Whitespace(" "),
		Char("("),
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
