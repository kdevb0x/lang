package irgen

import (
	"testing"

	"github.com/driusan/lang/compiler/ir"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/sampleprograms"
)

func TestIRGenEmptyMain(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.EmptyMain)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	if len(i.Body) != 0 {
		t.Error("Unexpected body for empty main function.")
	}
}

func TestIRGenHelloWorld(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.HelloWorld)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}

	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.StringLiteral(`Hello, world!\n`),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{ir.Fname("printf")},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}
	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenLetStatement(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.LetStatement)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(5),
			Dst: ir.LocalValue(0),
		},
		ir.MOV{
			Src: ir.StringLiteral(`%d\n`),
			Dst: ir.FuncCallArg(0),
		},
		ir.MOV{
			Src: ir.LocalValue(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.CALL{ir.Fname("printf")},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenHelloWorld2(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.HelloWorld2)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.StringLiteral(`%s %s\n %s`),
			Dst: ir.FuncCallArg(0),
		},
		ir.MOV{
			Src: ir.StringLiteral(`Hello, world!\n`),
			Dst: ir.FuncCallArg(1),
		},
		ir.MOV{
			Src: ir.StringLiteral(`World??`),
			Dst: ir.FuncCallArg(2),
		},
		ir.MOV{
			Src: ir.StringLiteral(`Hello, world!\n`),
			Dst: ir.FuncCallArg(3),
		},
		ir.CALL{ir.Fname("printf")},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenTwoProcs(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.TwoProcs)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.FuncRetVal(0),
		},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, err = GenerateIR(ast[1])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.CALL{Fname: "foo"},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.MOV{
			Src: ir.StringLiteral("%d"),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{Fname: "printf"},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenOutOfOrder(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.OutOfOrder)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.CALL{Fname: "foo"},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.MOV{
			Src: ir.StringLiteral("%d"),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{Fname: "printf"},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, err = GenerateIR(ast[1])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected = []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.FuncRetVal(0),
		},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenMutAddition(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.MutAddition)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.LocalValue(0),
		},
		ir.ADD{
			Src: ir.LocalValue(0),
			Dst: ir.LocalValue(2),
		},
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(2),
		},
		ir.MOV{
			Src: ir.LocalValue(2),
			Dst: ir.LocalValue(1),
		},
		ir.ADD{
			Src: ir.LocalValue(0),
			Dst: ir.LocalValue(3),
		},
		ir.ADD{
			Src: ir.LocalValue(1),
			Dst: ir.LocalValue(4),
		},
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(4),
		},
		ir.ADD{
			Src: ir.LocalValue(4),
			Dst: ir.LocalValue(3),
		},
		ir.MOV{
			Src: ir.LocalValue(3),
			Dst: ir.LocalValue(0),
		},
		ir.MOV{
			Src: ir.StringLiteral(`%d\n`),
			Dst: ir.FuncCallArg(0),
		},
		ir.MOV{
			Src: ir.LocalValue(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.CALL{Fname: "printf"},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenSimpleFunc(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.SimpleFunc)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.FuncRetVal(0),
		},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, err = GenerateIR(ast[1])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.CALL{Fname: "foo"},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.MOV{
			Src: ir.StringLiteral("%d"),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{Fname: "printf"},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenSumToTen(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.SumToTen)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "sum")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.FuncArg(0),
			Dst: ir.LocalValue(0),
		},
		ir.MOV{
			Src: ir.IntLiteral(0),
			Dst: ir.LocalValue(1),
		},
		ir.Label("loop0cond"),
		ir.JLE{
			ir.ConditionalJump{Label: ir.Label("loop0end"), Src: ir.LocalValue(0), Dst: ir.IntLiteral(0)},
		},
		ir.ADD{
			Src: ir.LocalValue(1),
			Dst: ir.LocalValue(2),
		},
		ir.ADD{
			Src: ir.LocalValue(0),
			Dst: ir.LocalValue(2),
		},
		ir.MOV{
			Src: ir.LocalValue(2),
			Dst: ir.LocalValue(1),
		},
		ir.MOV{
			Src: ir.LocalValue(0),
			Dst: ir.LocalValue(3),
		},
		ir.SUB{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(3),
		},
		ir.MOV{
			Src: ir.LocalValue(3),
			Dst: ir.LocalValue(0),
		},
		ir.JMP{"loop0cond"},
		ir.Label("loop0end"),
		ir.MOV{
			Src: ir.LocalValue(1),
			Dst: ir.FuncRetVal(0),
		},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
	i, err = GenerateIR(ast[1])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(10),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{"sum"},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.MOV{
			Src: ir.StringLiteral(`%d\n`),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{"printf"},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenSumToTenRecursive(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.SumToTenRecursive)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "sum")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(0),
			Dst: ir.FuncCallArg(0),
		},
		ir.MOV{
			Src: ir.FuncArg(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.CALL{Fname: "partial_sum"},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, err = GenerateIR(ast[1])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "partial_sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "partial_sum")
	}
	expected = []ir.Opcode{
		ir.JNE{
			ir.ConditionalJump{Label: "if1else", Src: ir.FuncArg(1), Dst: ir.IntLiteral(0)},
		},
		ir.MOV{
			Src: ir.FuncArg(0),
			Dst: ir.FuncRetVal(0),
		},
		ir.RET{},
		ir.JMP{"if1elsedone"},
		ir.Label("if1else"),
		ir.Label("if1elsedone"),
		ir.ADD{
			Src: ir.FuncArg(0),
			Dst: ir.LocalValue(0),
		},
		ir.ADD{
			Src: ir.FuncArg(1),
			Dst: ir.LocalValue(0),
		},
		ir.MOV{
			Src: ir.LocalValue(0),
			Dst: ir.FuncCallArg(0),
		},
		ir.MOV{
			Src: ir.FuncArg(1),
			Dst: ir.LocalValue(1),
		},
		ir.SUB{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(1),
		},
		ir.MOV{
			Src: ir.LocalValue(1),
			Dst: ir.FuncCallArg(1),
		},
		ir.CALL{"partial_sum"},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, err = GenerateIR(ast[2])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(10),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{"sum"},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.FuncCallArg(1),
		},
		ir.MOV{
			Src: ir.StringLiteral(`%d\n`),
			Dst: ir.FuncCallArg(0),
		},
		ir.CALL{"printf"},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

}

func TestIRGenFizzBuzz(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.Fizzbuzz)
	if err != nil {
		t.Fatal(err)
	}

	i, err := GenerateIR(ast[0])
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{Src: ir.IntLiteral(0), Dst: ir.LocalValue(0)},
		ir.MOV{Src: ir.IntLiteral(1), Dst: ir.LocalValue(1)},
		ir.Label("loop2cond"),
		ir.JE{ir.ConditionalJump{Label: "loop2end", Src: ir.LocalValue(0), Dst: ir.IntLiteral(1)}},
		ir.MOD{Left: ir.LocalValue(1), Right: ir.IntLiteral(15), Dst: ir.LocalValue(2)},
		ir.JNE{ir.ConditionalJump{Label: "if3else", Src: ir.LocalValue(2), Dst: ir.IntLiteral(0)}},
		ir.MOV{Src: ir.StringLiteral(`fizzbuzz\n`), Dst: ir.FuncCallArg(0)},
		ir.CALL{"printf"},
		ir.JMP{"if3elsedone"},
		ir.Label("if3else"),
		ir.MOD{Left: ir.LocalValue(1), Right: ir.IntLiteral(5), Dst: ir.LocalValue(3)},
		ir.JNE{ir.ConditionalJump{Label: "if4else", Src: ir.LocalValue(3), Dst: ir.IntLiteral(0)}},
		ir.MOV{Src: ir.StringLiteral(`buzz\n`), Dst: ir.FuncCallArg(0)},
		ir.CALL{"printf"},
		ir.JMP{"if4elsedone"},
		ir.Label("if4else"),
		ir.MOD{Left: ir.LocalValue(1), Right: ir.IntLiteral(3), Dst: ir.LocalValue(4)},
		ir.JNE{ir.ConditionalJump{Label: "if5else", Src: ir.LocalValue(4), Dst: ir.IntLiteral(0)}},
		ir.MOV{Src: ir.StringLiteral(`fizz\n`), Dst: ir.FuncCallArg(0)},
		ir.CALL{"printf"},
		ir.JMP{"if5elsedone"},
		ir.Label("if5else"),
		ir.MOV{Src: ir.StringLiteral(`%d\n`), Dst: ir.FuncCallArg(0)},
		ir.MOV{Src: ir.LocalValue(1), Dst: ir.FuncCallArg(1)},
		ir.CALL{"printf"},
		ir.Label("if5elsedone"),
		ir.Label("if4elsedone"),
		ir.Label("if3elsedone"),
		ir.ADD{Src: ir.LocalValue(1), Dst: ir.LocalValue(5)},
		ir.ADD{Src: ir.IntLiteral(1), Dst: ir.LocalValue(5)},
		ir.MOV{Src: ir.LocalValue(5), Dst: ir.LocalValue(1)},
		ir.JL{ir.ConditionalJump{Label: ir.Label("if6else"), Src: ir.LocalValue(1), Dst: ir.IntLiteral(100)}},
		ir.MOV{Src: ir.IntLiteral(1), Dst: ir.LocalValue(0)},
		ir.JMP{"if6elsedone"},
		ir.Label("if6else"),
		ir.Label("if6elsedone"),
		ir.JMP{"loop2cond"},
		ir.Label("loop2end"),
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}
