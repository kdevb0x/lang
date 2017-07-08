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
	ast, ti, err := ast.Parse(sampleprograms.EmptyMain)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(ast[0], ti, nil)
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
	ast, ti, err := ast.Parse(sampleprograms.HelloWorld)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(ast[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}

	expected := []ir.Opcode{
		ir.CALL{FName: "PrintString", Args: []ir.Register{ir.StringLiteral(`Hello, world!\n`)}},
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
	as, ti, err := ast.Parse(sampleprograms.LetStatement)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(5),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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

func TestIRGenLetStatementShadow(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.LetStatementShadow)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(5),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},
		ir.MOV{
			Src: ir.StringLiteral("hello"),
			Dst: ir.LocalValue{1, ast.TypeInfo{0, false}},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.LocalValue{1, ast.TypeInfo{0, false}},
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

func TestIRGenHelloWorld2(t *testing.T) {
	ast, ti, err := ast.Parse(sampleprograms.HelloWorld2)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(ast[0], ti, nil)
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
	as, ti, err := ast.Parse(sampleprograms.TwoProcs)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
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

	i, _, err = GenerateIR(as[1], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.CALL{FName: "foo"},
		ir.MOV{
			Src: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt",
			Args: []ir.Register{
				ir.LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, err := ast.Parse(sampleprograms.OutOfOrder)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.CALL{FName: "foo", Args: []ir.Register{}},
		ir.MOV{
			Src: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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

	i, _, err = GenerateIR(as[1], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected = []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
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
	as, ti, err := ast.Parse(sampleprograms.MutAddition)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.LocalValue{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{2, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{2, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{2, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.LocalValue{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.LocalValue{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{4, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{4, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.LocalValue{4, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{3, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, err := ast.Parse(sampleprograms.SimpleFunc)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(3),
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
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

	i, _, err = GenerateIR(as[1], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.CALL{FName: "foo", Args: []ir.Register{}},
		ir.MOV{
			Src: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, err := ast.Parse(sampleprograms.SumToTen)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "sum")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.FuncArg{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.IntLiteral(0),
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.Label("loop0cond"),
		ir.JLE{
			ir.ConditionalJump{Label: ir.Label("loop0end"), Src: ir.LocalValue{0, ast.TypeInfo{8, true}}, Dst: ir.IntLiteral(0)},
		},
		ir.ADD{
			Src: ir.LocalValue{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{2, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.LocalValue{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{2, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{2, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.SUB{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{3, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.JMP{"loop0cond"},
		ir.Label("loop0end"),
		ir.MOV{
			Src: ir.LocalValue{1, ast.TypeInfo{8, true}},
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
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
	i, _, err = GenerateIR(as[1], ti, nil)
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
			Src: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, err := ast.Parse(sampleprograms.SumToTenRecursive)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
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
				ir.FuncArg{0, ast.TypeInfo{8, true}},
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

	i, _, err = GenerateIR(as[1], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "partial_sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "partial_sum")
	}
	expected = []ir.Opcode{
		ir.JNE{
			ir.ConditionalJump{
				Label: "if1else",
				Src:   ir.FuncArg{1, ast.TypeInfo{8, true}},
				Dst:   ir.IntLiteral(0),
			},
		},
		ir.MOV{
			Src: ir.FuncArg{0, ast.TypeInfo{8, true}},
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		ir.RET{},
		ir.JMP{"if1elsedone"},
		ir.Label("if1else"),
		ir.Label("if1elsedone"),
		ir.ADD{
			Src: ir.FuncArg{0, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.FuncArg{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.FuncArg{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.SUB{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.CALL{
			FName: "partial_sum",
			Args: []ir.Register{
				ir.LocalValue{0, ast.TypeInfo{8, true}},
				ir.LocalValue{1, ast.TypeInfo{8, true}},
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

	i, _, err = GenerateIR(as[2], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.CALL{FName: "sum", Args: []ir.Register{ir.IntLiteral(10)}},
		ir.MOV{Src: ir.FuncRetVal{0, ast.TypeInfo{8, true}}, Dst: ir.LocalValue{0, ast.TypeInfo{8, true}}},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, err := ast.Parse(sampleprograms.Fizzbuzz)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{Src: ir.IntLiteral(0), Dst: ir.LocalValue{0, ast.TypeInfo{1, false}}},
		ir.MOV{Src: ir.IntLiteral(1), Dst: ir.LocalValue{1, ast.TypeInfo{8, true}}},
		ir.Label("loop2cond"),
		ir.JE{ir.ConditionalJump{Label: "loop2end", Src: ir.LocalValue{0, ast.TypeInfo{1, false}}, Dst: ir.IntLiteral(1)}},
		ir.MOD{Left: ir.LocalValue{1, ast.TypeInfo{8, true}}, Right: ir.IntLiteral(15), Dst: ir.LocalValue{2, ast.TypeInfo{8, true}}},
		ir.JNE{ir.ConditionalJump{Label: "if3else", Src: ir.LocalValue{2, ast.TypeInfo{8, true}}, Dst: ir.IntLiteral(0)}},
		ir.CALL{FName: "PrintString", Args: []ir.Register{ir.StringLiteral(`fizzbuzz`)}},
		ir.JMP{"if3elsedone"},
		ir.Label("if3else"),
		ir.MOD{Left: ir.LocalValue{1, ast.TypeInfo{8, true}}, Right: ir.IntLiteral(5), Dst: ir.LocalValue{3, ast.TypeInfo{8, true}}},
		ir.JNE{ir.ConditionalJump{Label: "if4else", Src: ir.LocalValue{3, ast.TypeInfo{8, true}}, Dst: ir.IntLiteral(0)}},
		ir.CALL{FName: "PrintString", Args: []ir.Register{ir.StringLiteral(`buzz`)}},
		ir.JMP{"if4elsedone"},
		ir.Label("if4else"),
		ir.MOD{Left: ir.LocalValue{1, ast.TypeInfo{8, true}}, Right: ir.IntLiteral(3), Dst: ir.LocalValue{4, ast.TypeInfo{8, true}}},
		ir.JNE{ir.ConditionalJump{Label: "if5else", Src: ir.LocalValue{4, ast.TypeInfo{8, true}}, Dst: ir.IntLiteral(0)}},
		ir.CALL{FName: "PrintString", Args: []ir.Register{ir.StringLiteral(`fizz`)}},
		ir.JMP{"if5elsedone"},
		ir.Label("if5else"),
		ir.CALL{FName: "PrintInt", Args: []ir.Register{ir.LocalValue{1, ast.TypeInfo{8, true}}}},
		ir.Label("if5elsedone"),
		ir.Label("if4elsedone"),
		ir.Label("if3elsedone"),
		ir.CALL{FName: "PrintString", Args: []ir.Register{ir.StringLiteral(`\n`)}},
		ir.ADD{Src: ir.LocalValue{1, ast.TypeInfo{8, true}}, Dst: ir.LocalValue{5, ast.TypeInfo{8, true}}},
		ir.ADD{Src: ir.IntLiteral(1), Dst: ir.LocalValue{5, ast.TypeInfo{8, true}}},
		ir.MOV{Src: ir.LocalValue{5, ast.TypeInfo{8, true}}, Dst: ir.LocalValue{1, ast.TypeInfo{8, true}}},
		ir.JL{ir.ConditionalJump{Label: ir.Label("if6else"), Src: ir.LocalValue{1, ast.TypeInfo{8, true}}, Dst: ir.IntLiteral(100)}},
		ir.MOV{Src: ir.IntLiteral(1), Dst: ir.LocalValue{0, ast.TypeInfo{1, false}}},
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
	as, ti, err := ast.Parse(sampleprograms.SomeMath)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.IntLiteral(2),
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},

		ir.MOV{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.SUB{
			Src: ir.IntLiteral(2),
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{3, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{2, ast.TypeInfo{8, true}},
		},

		ir.MUL{
			Left:  ir.IntLiteral(2),
			Right: ir.IntLiteral(3),
			Dst:   ir.LocalValue{5, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{5, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{4, ast.TypeInfo{8, true}},
		},
		ir.DIV{
			Left:  ir.IntLiteral(6),
			Right: ir.IntLiteral(2),
			Dst:   ir.LocalValue{7, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{7, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{6, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.IntLiteral(1),
			Dst: ir.LocalValue{9, ast.TypeInfo{8, true}},
		},
		ir.MUL{
			Left:  ir.IntLiteral(2),
			Right: ir.IntLiteral(3),
			Dst:   ir.LocalValue{11, ast.TypeInfo{8, true}},
		},
		ir.DIV{
			Left:  ir.IntLiteral(4),
			Right: ir.IntLiteral(2),
			Dst:   ir.LocalValue{12, ast.TypeInfo{8, true}},
		},
		ir.SUB{
			Src: ir.LocalValue{12, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{11, ast.TypeInfo{8, true}},
		},
		ir.ADD{
			Src: ir.LocalValue{11, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{9, ast.TypeInfo{8, true}},
		},
		ir.MOV{
			Src: ir.LocalValue{9, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{8, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`Add: `),
		},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},

		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`Sub: `),
		},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{2, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},

		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`Mul: `),
		},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{4, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},

		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`Div: `),
		},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{6, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},

		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`Complex: `),
		},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{8, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
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

func TestIRGenUserType(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.UserDefinedType)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[1], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(4),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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

func TestIRGenConcreteTypeUint8(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeUint8)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(4),
			Dst: ir.LocalValue{0, ast.TypeInfo{1, false}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{1, false}},
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

func TestIRGenConcreteTypeInt8(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeInt8)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(-4),
			Dst: ir.LocalValue{0, ast.TypeInfo{1, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{1, true}},
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

func TestIRGenConcreteTypeUint16(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeUint16)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(4),
			Dst: ir.LocalValue{0, ast.TypeInfo{2, false}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{2, false}},
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

func TestIRGenConcreteTypeInt16(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeInt16)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(-4),
			Dst: ir.LocalValue{0, ast.TypeInfo{2, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{2, true}},
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

func TestIRGenConcreteTypeUint32(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeUint32)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(4),
			Dst: ir.LocalValue{0, ast.TypeInfo{4, false}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{4, false}},
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

func TestIRGenConcreteTypeInt32(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeInt32)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(-4),
			Dst: ir.LocalValue{0, ast.TypeInfo{4, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{4, true}},
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

func TestIRGenConcreteTypeUint64(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeUint64)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(4),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, false}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, false}},
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

func TestIRGenConcreteTypeInt64(t *testing.T) {
	as, ti, err := ast.Parse(sampleprograms.ConcreteTypeInt64)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(-4),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, true}},
		},
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{0, ast.TypeInfo{8, true}},
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

func TestIRGenFibonacci(t *testing.T) {
	loopNum = 0
	as, ti, err := ast.Parse(sampleprograms.Fibonacci)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "fib_rec" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.ADD{
			Src: ir.FuncArg{0, ast.TypeInfo{8, false}},
			Dst: ir.LocalValue{1, ast.TypeInfo{8, false}},
		},
		ir.ADD{
			Src: ir.FuncArg{1, ast.TypeInfo{8, false}},
			Dst: ir.LocalValue{1, ast.TypeInfo{8, false}},
		},
		ir.MOV{
			Src: ir.LocalValue{1, ast.TypeInfo{8, false}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, false}},
		},
		ir.JL{ir.ConditionalJump{Label: "if0else", Src: ir.LocalValue{0, ast.TypeInfo{8, false}}, Dst: ir.IntLiteral(200)}},
		ir.MOV{
			Src: ir.FuncArg{1, ast.TypeInfo{8, false}},
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		ir.RET{},
		ir.JMP{"if0elsedone"},
		ir.Label("if0else"),
		ir.Label("if0elsedone"),
		ir.CALL{
			FName: "PrintInt",
			Args: []ir.Register{
				ir.LocalValue{0, ast.TypeInfo{8, false}},
			},
		},
		ir.CALL{
			FName: "PrintString",
			Args: []ir.Register{
				ir.StringLiteral(`\n`),
			},
		},
		ir.CALL{
			FName: "fib_rec",
			Args: []ir.Register{
				ir.FuncArg{1, ast.TypeInfo{8, false}},
				ir.LocalValue{0, ast.TypeInfo{8, false}},
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

	i, _, err = GenerateIR(as[1], ti, nil)
	if err != nil {
		t.Fatal(err)
	}

	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []ir.Opcode{
		ir.CALL{FName: "fib_rec", Args: []ir.Register{
			ir.IntLiteral(1),
			ir.IntLiteral(1),
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

func TestIREnumType(t *testing.T) {
	loopNum = 0
	as, ti, err := ast.Parse(sampleprograms.EnumType)
	if err != nil {
		t.Fatal(err)
	}

	_, enums, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := GenerateIR(as[1], ti, enums)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []ir.Opcode{
		ir.MOV{
			Src: ir.IntLiteral(0),
			Dst: ir.LocalValue{0, ast.TypeInfo{8, false}},
		},
		ir.JE{ir.ConditionalJump{Label: "match0v0", Src: ir.LocalValue{0, ast.TypeInfo{8, false}}, Dst: ir.IntLiteral(0)}},
		ir.JE{ir.ConditionalJump{Label: "match0v1", Src: ir.LocalValue{0, ast.TypeInfo{8, false}}, Dst: ir.IntLiteral(1)}},
		ir.JMP{"match0done"},
		ir.Label("match0v0"),
		ir.CALL{
			FName: "PrintString",
			Args: []ir.Register{
				ir.StringLiteral(`I am A!\n`),
			},
		},
		ir.JMP{"match0done"},

		ir.Label("match0v1"),
		ir.CALL{
			FName: "PrintString",
			Args: []ir.Register{
				ir.StringLiteral(`I am B!\n`),
			},
		},
		ir.JMP{"match0done"},
		ir.Label("match0done"),
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

func TestIRGenericEnumType(t *testing.T) {
	loopNum = 0
	as, ti, err := ast.Parse(sampleprograms.GenericEnumType)
	if err != nil {
		t.Fatal(err)
	}
	_, enums, err := GenerateIR(as[0], ti, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := GenerateIR(as[1], ti, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected := []ir.Opcode{
		ir.JLE{ir.ConditionalJump{
			Label: "if0else",
			Src:   ir.FuncArg{0, ast.TypeInfo{8, true}},
			Dst:   ir.IntLiteral(3),
		},
		},
		ir.MOV{
			Src: ir.IntLiteral(0),
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		ir.RET{},
		ir.JMP{"if0elsedone"},
		ir.Label("if0else"),
		ir.Label("if0elsedone"),
		// Enum type goes into the first word
		ir.MOV{
			Src: ir.IntLiteral(1),
			Dst: ir.FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		// The concrete parameter is an int, which goes into the
		// next word.
		ir.MOV{
			Src: ir.IntLiteral(5),
			Dst: ir.FuncRetVal{1, ast.TypeInfo{8, true}},
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
	i, _, err = GenerateIR(as[2], ti, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected = []ir.Opcode{
		ir.CALL{FName: "DoSomething", Args: []ir.Register{
			ir.IntLiteral(3),
		},
		},
		ir.MOV{
			Src: ir.FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: ir.LocalValue{0, ast.TypeInfo{8, false}},
		},
		ir.MOV{
			Src: ir.FuncRetVal{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		ir.JE{ir.ConditionalJump{
			Label: "match1v0",
			Src:   ir.LocalValue{0, ast.TypeInfo{8, false}},
			Dst:   ir.IntLiteral(0),
		},
		},
		ir.JE{ir.ConditionalJump{
			Label: "match1v1",
			Src:   ir.LocalValue{0, ast.TypeInfo{8, false}},
			Dst:   ir.IntLiteral(1),
		},
		},
		ir.JMP{"match1done"},
		ir.Label("match1v0"),
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`I am nothing!\n`),
		},
		},
		ir.JMP{"match1done"},
		ir.Label("match1v1"),
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{1, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},

		ir.JMP{"match1done"},
		ir.Label("match1done"),
		ir.CALL{FName: "DoSomething", Args: []ir.Register{
			ir.IntLiteral(4),
		},
		},
		ir.MOV{
			Src: ir.FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: ir.LocalValue{2, ast.TypeInfo{8, false}},
		},
		ir.MOV{
			Src: ir.FuncRetVal{1, ast.TypeInfo{8, true}},
			Dst: ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		ir.JE{ir.ConditionalJump{
			Label: "match2v0",
			Src:   ir.LocalValue{2, ast.TypeInfo{8, false}},
			Dst:   ir.IntLiteral(0),
		},
		},
		ir.JE{ir.ConditionalJump{
			Label: "match2v1",
			Src:   ir.LocalValue{2, ast.TypeInfo{8, false}},
			Dst:   ir.IntLiteral(1),
		},
		},
		ir.JMP{"match2done"},
		ir.Label("match2v0"),
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`I am nothing!\n`),
		},
		},
		ir.JMP{"match2done"},
		ir.Label("match2v1"),
		ir.CALL{FName: "PrintInt", Args: []ir.Register{
			ir.LocalValue{3, ast.TypeInfo{8, true}},
		},
		},
		ir.CALL{FName: "PrintString", Args: []ir.Register{
			ir.StringLiteral(`\n`),
		},
		},
		ir.JMP{"match2done"},
		ir.Label("match2done"),
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
