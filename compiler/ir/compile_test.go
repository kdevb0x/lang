package ir

import (
	"fmt"
	"testing"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/sampleprograms"
)

func compareOp(a, b Opcode) bool {
	switch a1 := a.(type) {
	case CALL:
		b1, ok := b.(CALL)
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

func compareIR(i, expected []Opcode) error {
	if len(i) != len(expected) {
		return fmt.Errorf("Unexpected body: got %v want %v\n", i, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i[j]) {
			return fmt.Errorf("Unexpected value for opcode %d: got %v want %v", j, i[j], expected[j])
		}
	}
	return nil
}

func TestIRGenEmptyMain(t *testing.T) {
	ast, ti, c, err := ast.Parse(sampleprograms.EmptyMain)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(ast[0], ti, c, nil)
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
	ast, ti, c, err := ast.Parse(sampleprograms.HelloWorld)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(ast[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}

	expected := []Opcode{
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`Hello, world!\n`)}},
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
	as, ti, c, err := ast.Parse(sampleprograms.LetStatement)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.LetStatementShadow)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},
		MOV{
			Src: StringLiteral("hello"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		CALL{FName: "PrintString", Args: []Register{
			LocalValue{1, ast.TypeInfo{8, false}},
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
	as, ti, c, err := ast.Parse(sampleprograms.TwoProcs)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "foo"},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.OutOfOrder)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		CALL{FName: "foo", Args: []Register{}},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected = []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
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
	as, ti, c, err := ast.Parse(sampleprograms.MutAddition)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		MOV{
			Src: LocalValue{1, ast.TypeInfo{8, true}},
			Dst: TempValue(2),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(2),
		},
		ADD{
			Src: TempValue(2),
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.SimpleFunc)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "foo", Args: []Register{}},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.SumToTen)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "sum")
	}
	expected := []Opcode{
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, true}, false},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		Label("loop0cond"),
		JLE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: IntLiteral(0)},
		},
		MOV{
			Src: LocalValue{1, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		ADD{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		SUB{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		JMP{"loop0cond"},
		Label("loop0end"),
		MOV{
			Src: LocalValue{1, ast.TypeInfo{8, true}},
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if expected[j] != i.Body[j] {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "sum", Args: []Register{
			IntLiteral(10),
		},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.SumToTenRecursive)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "sum")
	}
	expected := []Opcode{
		CALL{
			FName: "partial_sum",
			Args: []Register{
				IntLiteral(0),
				FuncArg{0, ast.TypeInfo{8, true}, false},
			},
			TailCall: true,
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "partial_sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "partial_sum")
	}
	expected = []Opcode{
		JNE{
			ConditionalJump{
				Label: "if1else",
				Src:   FuncArg{1, ast.TypeInfo{8, true}, false},
				Dst:   IntLiteral(0),
			},
		},
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, true}, false},
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, true}, false},
			Dst: TempValue(0),
		},
		ADD{
			Src: FuncArg{1, ast.TypeInfo{8, true}, false},
			Dst: TempValue(0),
		},
		MOV{
			Src: FuncArg{1, ast.TypeInfo{8, true}, false},
			Dst: TempValue(1),
		},
		SUB{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		CALL{
			FName: "partial_sum",
			Args: []Register{
				TempValue(0),
				TempValue(1),
			},
			TailCall: true,
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, _, err = Generate(as[2], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "sum", Args: []Register{IntLiteral(10)}},
		MOV{Src: FuncRetVal{0, ast.TypeInfo{8, true}}, Dst: LocalValue{0, ast.TypeInfo{8, true}}},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.Fizzbuzz)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{Src: IntLiteral(0), Dst: LocalValue{0, ast.TypeInfo{1, false}}},
		MOV{Src: IntLiteral(1), Dst: LocalValue{1, ast.TypeInfo{8, true}}},
		Label("loop2cond"),
		JE{ConditionalJump{Label: "loop2end", Src: LocalValue{0, ast.TypeInfo{1, false}}, Dst: IntLiteral(1)}},
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(15), Dst: TempValue(0)},
		JNE{ConditionalJump{Label: "if3else", Src: TempValue(0), Dst: IntLiteral(0)}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`fizzbuzz`)}},
		JMP{"if3elsedone"},
		Label("if3else"),
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(5), Dst: TempValue(1)},
		JNE{ConditionalJump{Label: "if4else", Src: TempValue(1), Dst: IntLiteral(0)}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`buzz`)}},
		JMP{"if4elsedone"},
		Label("if4else"),
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(3), Dst: TempValue(2)},
		JNE{ConditionalJump{Label: "if5else", Src: TempValue(2), Dst: IntLiteral(0)}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`fizz`)}},
		JMP{"if5elsedone"},
		Label("if5else"),
		CALL{FName: "PrintInt", Args: []Register{LocalValue{1, ast.TypeInfo{8, true}}}},
		Label("if5elsedone"),
		Label("if4elsedone"),
		Label("if3elsedone"),
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		MOV{Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: TempValue(3)},
		ADD{Src: IntLiteral(1), Dst: TempValue(3)},
		MOV{Src: TempValue(3), Dst: LocalValue{1, ast.TypeInfo{8, true}}},
		JL{ConditionalJump{Label: Label("if6else"), Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: IntLiteral(100)}},
		MOV{Src: IntLiteral(1), Dst: LocalValue{0, ast.TypeInfo{1, false}}},
		JMP{"if6elsedone"},
		Label("if6else"),
		Label("if6elsedone"),
		JMP{"loop2cond"},
		Label("loop2end"),
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
	as, ti, c, err := ast.Parse(sampleprograms.SomeMath)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		ADD{
			Src: IntLiteral(2),
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},

		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		SUB{
			Src: IntLiteral(2),
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MUL{
			Left:  IntLiteral(2),
			Right: IntLiteral(3),
			Dst:   TempValue(2),
		},
		MOV{
			Src: TempValue(2),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		DIV{
			Left:  IntLiteral(6),
			Right: IntLiteral(2),
			Dst:   TempValue(3),
		},
		MOV{
			Src: TempValue(3),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(4),
		},
		MUL{
			Left:  IntLiteral(2),
			Right: IntLiteral(3),
			Dst:   TempValue(6),
		},
		DIV{
			Left:  IntLiteral(4),
			Right: IntLiteral(2),
			Dst:   TempValue(7),
		},
		SUB{
			Src: TempValue(7),
			Dst: TempValue(6),
		},
		ADD{
			Src: TempValue(6),
			Dst: TempValue(4),
		},
		MOV{
			Src: TempValue(4),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`Add: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`Sub: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{1, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`Mul: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{2, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`Div: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{3, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`Complex: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{4, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
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
	as, ti, c, err := ast.Parse(sampleprograms.UserDefinedType)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeUint8)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{0, ast.TypeInfo{1, false}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{1, false}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeInt8)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue{0, ast.TypeInfo{1, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{1, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeUint16)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{0, ast.TypeInfo{2, false}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{2, false}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeInt16)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue{0, ast.TypeInfo{2, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{2, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeUint32)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{0, ast.TypeInfo{4, false}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{4, false}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeInt32)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue{0, ast.TypeInfo{4, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{4, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeUint64)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, false}},
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
	as, ti, c, err := ast.Parse(sampleprograms.ConcreteTypeInt64)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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
	as, ti, c, err := ast.Parse(sampleprograms.Fibonacci)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "fib_rec" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, false}, false},
			Dst: TempValue(0),
		},
		ADD{
			Src: FuncArg{1, ast.TypeInfo{8, false}, false},
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		JL{ConditionalJump{Label: "if0else", Src: LocalValue{0, ast.TypeInfo{8, false}}, Dst: IntLiteral(200)}},
		MOV{
			Src: FuncArg{1, ast.TypeInfo{8, false}, false},
			Dst: FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		RET{},
		JMP{"if0elsedone"},
		Label("if0else"),
		Label("if0elsedone"),
		CALL{
			FName: "PrintInt",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
			},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`\n`),
			},
		},
		CALL{
			FName: "fib_rec",
			Args: []Register{
				FuncArg{1, ast.TypeInfo{8, false}, false},
				LocalValue{0, ast.TypeInfo{8, false}},
			},
			TailCall: true,
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "fib_rec", Args: []Register{
			IntLiteral(1),
			IntLiteral(1),
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
	as, ti, c, err := ast.Parse(sampleprograms.EnumType)
	if err != nil {
		t.Fatal(err)
	}

	_, enums, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		JE{ConditionalJump{Label: "match0v0", Src: LocalValue{0, ast.TypeInfo{8, false}}, Dst: IntLiteral(0)}},
		JE{ConditionalJump{Label: "match0v1", Src: LocalValue{0, ast.TypeInfo{8, false}}, Dst: IntLiteral(1)}},
		JMP{"match0done"},
		Label("match0v0"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`I am A!\n`),
			},
		},
		JMP{"match0done"},

		Label("match0v1"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`I am B!\n`),
			},
		},
		JMP{"match0done"},
		Label("match0done"),
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
	as, ti, c, err := ast.Parse(sampleprograms.GenericEnumType)
	if err != nil {
		t.Fatal(err)
	}
	_, enums, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		JLE{ConditionalJump{
			Label: "if0else",
			Src:   FuncArg{0, ast.TypeInfo{8, true}, false},
			Dst:   IntLiteral(3),
		},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		RET{},
		JMP{"if0elsedone"},
		Label("if0else"),
		Label("if0elsedone"),
		// Enum type goes into the first word
		MOV{
			Src: IntLiteral(1),
			Dst: FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		// The concrete parameter is an int, which goes into the
		// next word.
		MOV{
			Src: IntLiteral(5),
			Dst: FuncRetVal{1, ast.TypeInfo{8, true}},
		},
		RET{},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
	i, _, err = Generate(as[2], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		CALL{FName: "DoSomething", Args: []Register{
			IntLiteral(3),
		},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: FuncRetVal{1, ast.TypeInfo{8, true}},
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		JE{ConditionalJump{
			Label: "match1v0",
			Src:   LocalValue{0, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(0),
		},
		},
		JE{ConditionalJump{
			Label: "match1v1",
			Src:   LocalValue{0, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(1),
		},
		},
		JMP{"match1done"},
		Label("match1v0"),
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`I am nothing!\n`),
		},
		},
		JMP{"match1done"},
		Label("match1v1"),
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{1, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},

		JMP{"match1done"},
		Label("match1done"),
		CALL{FName: "DoSomething", Args: []Register{
			IntLiteral(4),
		},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: FuncRetVal{1, ast.TypeInfo{8, true}},
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		JE{ConditionalJump{
			Label: "match2v0",
			Src:   LocalValue{2, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(0),
		},
		},
		JE{ConditionalJump{
			Label: "match2v1",
			Src:   LocalValue{2, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(1),
		},
		},
		JMP{"match2done"},
		Label("match2v0"),
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`I am nothing!\n`),
		},
		},
		JMP{"match2done"},
		Label("match2v1"),
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{3, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			StringLiteral(`\n`),
		},
		},
		JMP{"match2done"},
		Label("match2done"),
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

func TestIRMatchParam(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.MatchParam)
	if err != nil {
		t.Fatal(err)
	}
	_, enums, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		JE{ConditionalJump{
			Label: "match0v0",
			Src:   FuncArg{0, ast.TypeInfo{0, false}, false},
			Dst:   IntLiteral(1),
		},
		},
		JE{ConditionalJump{
			Label: "match0v1",
			Src:   FuncArg{0, ast.TypeInfo{0, false}, false},
			Dst:   IntLiteral(0),
		},
		},
		JMP{"match0done"},
		Label("match0v0"),
		MOV{
			Src: FuncArg{1, ast.TypeInfo{8, true}, false},
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
		JMP{"match0done"},
		Label("match0v1"),
		MOV{
			Src: IntLiteral(0),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
		JMP{"match0done"},
		Label("match0done"),
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", i.Body, expected)
	}

	for j := range expected {
		if !compareOp(expected[j], i.Body[j]) {
			t.Errorf("Unexpected value for opcode %d: got %v want %v", j, i.Body[j], expected[j])
		}
	}
	i, _, err = Generate(as[2], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		CALL{FName: "foo", Args: []Register{
			IntLiteral(1),
			IntLiteral(5),
		},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
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

func TestIRSimpleAlgorithm(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.SimpleAlgorithm)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MUL{
			Left:  FuncArg{0, ast.TypeInfo{8, true}, false},
			Right: IntLiteral(2),
			Dst:   TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		Label("loop0cond"),
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: LocalValue{2, ast.TypeInfo{8, true}}},
		},
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(2), Dst: TempValue(1)},
		JNE{
			ConditionalJump{Label: Label("if1else"), Src: TempValue(1), Dst: IntLiteral(0)},
		},
		MOV{Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: TempValue(2)},
		MUL{
			Left:  LocalValue{1, ast.TypeInfo{8, true}},
			Right: IntLiteral(2),
			Dst:   TempValue(3),
		},
		ADD{Src: TempValue(3), Dst: TempValue(2)},
		MOV{Src: TempValue(2), Dst: LocalValue{0, ast.TypeInfo{8, true}}},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		MOV{Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: TempValue(4)},
		ADD{Src: IntLiteral(1), Dst: TempValue(4)},
		MOV{Src: TempValue(4), Dst: LocalValue{1, ast.TypeInfo{8, true}}},
		JMP{"loop0cond"},
		Label("loop0end"),
		MOV{Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: FuncRetVal{0, ast.TypeInfo{8, true}}},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected = []Opcode{
		CALL{FName: "loop", Args: []Register{
			IntLiteral(10),
		},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
		},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatal(err)
	}

}

func TestIRSimpleArray(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.SimpleArray)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Offset: IntLiteral(3 * 8),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRArrayMutation(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.ArrayMutation)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Offset: IntLiteral(24),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
			// FIXME: This should be
			//	Dst:Offset{
			//		Base: LocalValue{0, ast.TypeInfo{8, true}},
			//		Offset: IntLiteral(24),
			//	},
			// But for now they evaluate to the same address ..
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Offset: IntLiteral(24),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Offset: IntLiteral(16),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRReferenceVariable(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.ReferenceVariable)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: FuncArg{0, ast.TypeInfo{8, true}, true},
		},
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, true}, true},
			Dst: TempValue(0),
		},
		ADD{
			Src: FuncArg{1, ast.TypeInfo{8, true}, false},
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	expected = []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{0, ast.TypeInfo{8, true}}}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{
			FName: "changer",
			Args: []Register{
				Pointer{LocalValue{0, ast.TypeInfo{8, true}}},
				IntLiteral(3),
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{0, ast.TypeInfo{8, true}}}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{1, ast.TypeInfo{8, true}}}},
	}
	i, _, err = Generate(as[1], ti, c, nil)

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSimpleSlice(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.SimpleSlice)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	// LV0 == size of the slice, 1-5=the values
	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Offset: IntLiteral(24),
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSimpleSliceInference(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.SimpleSliceInference)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Offset: IntLiteral(24),
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSimpleSliceMutation(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.SliceMutation)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Offset: IntLiteral(24),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
			/* FIXME: This should be:
				but they're the same value, so it doesn't really matter for now
			Dst: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(24),
			},*/
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Offset: IntLiteral(24),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Offset: IntLiteral(16),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSliceParam(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.SliceParam)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(44),
			Dst: LocalValue{1, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(55),
			Dst: LocalValue{2, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(88),
			Dst: LocalValue{3, ast.TypeInfo{1, false}},
		},
		CALL{
			FName: "PrintASlice",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				FuncArg{0, ast.TypeInfo{8, false}, false},
				FuncArg{1, ast.TypeInfo{8, false}, false},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRWriteSyscall(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.WriteSyscall)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		CALL{
			FName: "Write",
			Args: []Register{
				IntLiteral(1),
				StringLiteral("Stdout!"),
			},
		},
		CALL{
			FName: "Write",
			Args: []Register{
				IntLiteral(2),
				StringLiteral("Stderr!"),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRReadSyscall(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.ReadSyscall)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		CALL{
			FName: "Open",
			Args: []Register{
				StringLiteral("foo.txt"),
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(6),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{2, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{3, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{4, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{5, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{6, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{7, ast.TypeInfo{1, false}},
		},
		CALL{
			FName: "Read",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				LocalValue{1, ast.TypeInfo{8, false}},
				Pointer{LocalValue{2, ast.TypeInfo{1, false}}},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{8, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue{1, ast.TypeInfo{8, false}},
				Pointer{LocalValue{2, ast.TypeInfo{1, false}}},
			},
		},
		CALL{
			FName: "Close",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRIfElseMatch(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.IfElseMatch)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		JGE{
			ConditionalJump{
				"comparison1false",
				LocalValue{0, ast.TypeInfo{8, true}},
				IntLiteral(3),
			},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		JMP{"comparison1done"},
		Label("comparison1false"),
		MOV{
			Src: IntLiteral(0),
			Dst: TempValue(0),
		},
		Label("comparison1done"),
		JE{
			ConditionalJump{
				"match0v0",
				IntLiteral(1),
				TempValue(0),
			},
		},
		JLE{
			ConditionalJump{
				"comparison2false",
				LocalValue{0, ast.TypeInfo{8, true}},
				IntLiteral(3),
			},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		JMP{"comparison2done"},
		Label("comparison2false"),
		MOV{
			Src: IntLiteral(0),
			Dst: TempValue(1),
		},
		Label("comparison2done"),
		JE{
			ConditionalJump{
				"match0v1",
				IntLiteral(1),
				TempValue(1),
			},
		},
		JGE{
			ConditionalJump{
				"comparison3false",
				LocalValue{0, ast.TypeInfo{8, true}},
				IntLiteral(4),
			},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(2),
		},
		JMP{"comparison3done"},
		Label("comparison3false"),
		MOV{
			Src: IntLiteral(0),
			Dst: TempValue(2),
		},
		Label("comparison3done"),
		JE{
			ConditionalJump{
				"match0v2",
				IntLiteral(1),
				TempValue(2),
			},
		},
		JMP{"match0done"},
		Label("match0v0"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`x is less than 3\n`),
			},
		},
		JMP{"match0done"},
		Label("match0v1"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`x is greater than 3\n`),
			},
		},
		JMP{"match0done"},
		Label("match0v2"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`x is less than 4\n`),
			},
		},
		JMP{"match0done"},
		Label("match0done"),
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

/*
func TestIREcho(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.Echo)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "len",
			Args: []Register{
				FuncArg{0, ast.TypeInfo{8, false}, false},
				FuncArg{1, ast.TypeInfo{8, false}, false},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		Label("loop0cond"),
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: LocalValue{1, ast.TypeInfo{8, false}}},
		},
		MUL{
			Left:  IntLiteral(8),
			Right: LocalValue{1, ast.TypeInfo{8, true}},
			Dst:   LocalValue{2, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				FuncArg{1, ast.TypeInfo{0, false}, false},
				Offset{
					Offset: LocalValue{2, ast.TypeInfo{8, true}},
					Base:   FuncArg{1, ast.TypeInfo{8, false}, false},
				},
			},
		},
		ADD{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		ADD{
			Src: IntLiteral(1),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{3, ast.TypeInfo{8, true}},
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		JE{
			ConditionalJump{Label: Label("if1else"), Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: LocalValue{1, ast.TypeInfo{8, false}}},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(" "),
			},
		},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		JMP{"loop0cond"},
		Label("loop0end"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				StringLiteral(`\n`),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}
*/

func TestArrayIndex(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.ArrayIndex)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		// let x = 3
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		// Let statement
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		// mutable statement
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{6, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue{7, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{8, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{9, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{10, ast.TypeInfo{8, true}},
		},
		// Convert index from index into byte offset
		MUL{
			Left:  IntLiteral(8),
			Right: LocalValue{0, ast.TypeInfo{8, true}},
			Dst:   TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Offset: TempValue(0),
				},
			},
		},
		CALL{
			FName: "PrintString",
			Args:  []Register{StringLiteral(`\n`)},
		},
		// Convert x+1 offset from index into byte offset
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		MUL{
			Left:  IntLiteral(8),
			Right: TempValue(1),
			Dst:   TempValue(2),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{6, ast.TypeInfo{8, true}},
					Offset: TempValue(2),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIndexAssignment(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.IndexAssignment)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		// let x []int = { 3, 4, 5 }
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		// Let statement
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(8),
			},
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(16),
			},
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{4, ast.TypeInfo{8, true}},
		}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{5, ast.TypeInfo{8, true}}}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIndexedAddition(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.IndexedAddition)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		// let x []int = { 3, 4, 5 }
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		// Let statement
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(8),
			},
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{4, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		ADD{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(16),
			},
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(16),
			},
			Dst: TempValue(1),
		},
		ADD{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(0),
			},
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{4, ast.TypeInfo{8, true}},
		}},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{5, ast.TypeInfo{8, true}}}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestStringArray(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.StringArray)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, false}},
					Offset: IntLiteral(8),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, false}},
					Offset: IntLiteral(0),
				},
			},
		},
	}
	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPreEcho(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.PreEcho)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("baz"),
			Dst: LocalValue{3, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "len",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{8, false}}},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{5, ast.TypeInfo{8, false}},
		},
		Label("loop0cond"),
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{4, ast.TypeInfo{8, true}}, Dst: LocalValue{5, ast.TypeInfo{8, false}}},
		},
		MUL{
			Left:  IntLiteral(8),
			Right: LocalValue{4, ast.TypeInfo{8, true}},
			Dst:   TempValue(0),
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, false}},
					Offset: TempValue(0),
				},
			},
		},
		MOV{
			Src: LocalValue{4, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		JE{
			ConditionalJump{Label: Label("if1else"), Src: LocalValue{4, ast.TypeInfo{8, true}}, Dst: LocalValue{5, ast.TypeInfo{8, false}}},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(" ")}},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		JMP{"loop0cond"},
		Label("loop0end"),
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPreEcho2(t *testing.T) {
	loopNum = 0
	as, ti, c, err := ast.Parse(sampleprograms.PreEcho2)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "len",
			Args: []Register{
				FuncArg{0, ast.TypeInfo{8, false}, false},
				FuncArg{1, ast.TypeInfo{8, false}, false},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		Label("loop0cond"),
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: LocalValue{1, ast.TypeInfo{8, false}}},
		},
		MUL{
			Left:  IntLiteral(8),
			Right: LocalValue{0, ast.TypeInfo{8, true}},
			Dst:   TempValue(0),
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   FuncArg{1, ast.TypeInfo{8, false}, false},
					Offset: TempValue(0),
				},
			},
		},
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
		JE{
			ConditionalJump{Label: Label("if1else"), Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: LocalValue{1, ast.TypeInfo{8, false}}},
		},
		CALL{FName: "PrintString", Args: []Register{StringLiteral(" ")}},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		JMP{"loop0cond"},
		Label("loop0end"),
		CALL{FName: "PrintString", Args: []Register{StringLiteral(`\n`)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i2, _, err := Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("baz"),
			Dst: LocalValue{3, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintSlice",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{8, false}}},
			},
		},
	}
	if err := compareIR(i2.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}
