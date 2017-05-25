package irgen

import (
	"testing"

	"github.com/driusan/lang/compiler/ir"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/sampleprograms"
)

func compareOp(a, b ir.Opcode) bool {
	switch a1 := a.(type) {
	case ir.CALL:
		b1, ok := b.(ir.CALL)
		if !ok {
			return false
		}
		if b1.FName != a1.FName {
			return false
		}
		if len(b1.Args) != len(a1.Args) {
			return false
		}
		for i := range a1.Args {
			if a1.Args[i] != b1.Args[i] {
				return false
			}
		}
		return b1.TailCall == a1.TailCall
	default:
		return a == b
	}
}
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
		ir.CALL{FName: "printf", Args: []ir.Register{ir.StringLiteral(`Hello, world!\n`)}},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}
	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "printf", Args: []ir.Register{ir.StringLiteral(`%d\n`), ir.LocalValue(0)}},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`%s %s\n %s`),
			ir.StringLiteral(`Hello, world!\n`),
			ir.StringLiteral(`World??`),
			ir.StringLiteral(`Hello, world!\n`),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "foo"},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.LocalValue(0),
		},
		ir.CALL{FName: "printf",
			Args: []ir.Register{
				ir.StringLiteral(`%d`),
				ir.LocalValue(0),
			},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "foo", Args: []ir.Register{}},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.LocalValue(0),
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`%d`),
			ir.LocalValue(0),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`%d\n`),
			ir.LocalValue(0),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "foo", Args: []ir.Register{}},
		ir.MOV{Src: ir.FuncRetVal(0), Dst: ir.LocalValue(0)},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral("%d"),
			ir.LocalValue(0),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "sum", Args: []ir.Register{
			ir.IntLiteral(10),
		},
		},
		ir.MOV{
			Src: ir.FuncRetVal(0),
			Dst: ir.LocalValue(0),
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`%d\n`),
			ir.LocalValue(0),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{
			FName: "partial_sum",
			Args: []ir.Register{
				ir.IntLiteral(0),
				ir.FuncArg(0),
			},
			TailCall: true,
		},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
			Src: ir.FuncArg(1),
			Dst: ir.LocalValue(1),
		},
		ir.SUB{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(1),
		},
		ir.CALL{
			FName: "partial_sum",
			Args: []ir.Register{
				ir.LocalValue(0),
				ir.LocalValue(1),
			},
			TailCall: true,
		},
		ir.RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "sum", Args: []ir.Register{ir.IntLiteral(10)}},
		ir.MOV{Src: ir.FuncRetVal(0), Dst: ir.LocalValue(0)},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`%d\n`),
			ir.LocalValue(0),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
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
		ir.CALL{FName: "printf", Args: []ir.Register{ir.StringLiteral(`fizzbuzz\n`)}},
		ir.JMP{"if3elsedone"},
		ir.Label("if3else"),
		ir.MOD{Left: ir.LocalValue(1), Right: ir.IntLiteral(5), Dst: ir.LocalValue(3)},
		ir.JNE{ir.ConditionalJump{Label: "if4else", Src: ir.LocalValue(3), Dst: ir.IntLiteral(0)}},
		ir.CALL{FName: "printf", Args: []ir.Register{ir.StringLiteral(`buzz\n`)}},
		ir.JMP{"if4elsedone"},
		ir.Label("if4else"),
		ir.MOD{Left: ir.LocalValue(1), Right: ir.IntLiteral(3), Dst: ir.LocalValue(4)},
		ir.JNE{ir.ConditionalJump{Label: "if5else", Src: ir.LocalValue(4), Dst: ir.IntLiteral(0)}},
		ir.CALL{FName: "printf", Args: []ir.Register{ir.StringLiteral(`fizz\n`)}},
		ir.JMP{"if5elsedone"},
		ir.Label("if5else"),
		ir.CALL{FName: "printf", Args: []ir.Register{ir.StringLiteral(`%d\n`), ir.LocalValue(1)}},
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
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}

func TestIRGenSomeMathStatement(t *testing.T) {
	ast, err := ast.Parse(sampleprograms.SomeMath)
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
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(1),
		},
		ir.ADD{
			Src: ir.IntLiteral(2),
			Dst: ir.LocalValue(1),
		},
		ir.MOV{
			Src: ir.LocalValue(1),
			Dst: ir.LocalValue(0),
		},

		ir.MOV{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(3),
		},
		ir.SUB{
			Src: ir.IntLiteral(2),
			Dst: ir.LocalValue(3),
		},
		ir.MOV{
			Src: ir.LocalValue(3),
			Dst: ir.LocalValue(2),
		},

		ir.MUL{
			Left:  ir.IntLiteral(2),
			Right: ir.IntLiteral(3),
			Dst:   ir.LocalValue(5),
		},
		ir.MOV{
			Src: ir.LocalValue(5),
			Dst: ir.LocalValue(4),
		},
		ir.DIV{
			Left:  ir.IntLiteral(6),
			Right: ir.IntLiteral(2),
			Dst:   ir.LocalValue(7),
		},
		ir.MOV{
			Src: ir.LocalValue(7),
			Dst: ir.LocalValue(6),
		},
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue(9),
		},
		ir.MUL{
			Left:  ir.IntLiteral(2),
			Right: ir.IntLiteral(3),
			Dst:   ir.LocalValue(11),
		},
		ir.DIV{
			Left:  ir.IntLiteral(4),
			Right: ir.IntLiteral(2),
			Dst:   ir.LocalValue(12),
		},
		ir.SUB{
			Src: ir.LocalValue(12),
			Dst: ir.LocalValue(11),
		},
		ir.ADD{
			Src: ir.LocalValue(11),
			Dst: ir.LocalValue(9),
		},
		ir.MOV{
			Src: ir.LocalValue(9),
			Dst: ir.LocalValue(8),
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`Add: %d\n`),
			ir.LocalValue(0),
		},
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`Sub: %d\n`),
			ir.LocalValue(2),
		},
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`Mul: %d\n`),
			ir.LocalValue(4),
		},
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`Div: %d\n`),
			ir.LocalValue(6),
		},
		},
		ir.CALL{FName: "printf", Args: []ir.Register{
			ir.StringLiteral(`Complex: %d\n`),
			ir.LocalValue(8),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
}
