package mlir

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
		return fmt.Errorf("Unexpected body (%d != %d): got %v want %v\n", len(i), len(expected), i, expected)
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
		CALL{FName: "PrintString", Args: []Register{IntLiteral(14), StringLiteral(`Hello, world!\n`)}},
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
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("hello"),
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		CALL{FName: "PrintString", Args: []Register{
			LocalValue{1, ast.TypeInfo{8, false}},
			LocalValue{2, ast.TypeInfo{8, false}},
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
		CALL{FName: "PrintInt",
			Args: []Register{
				FuncRetVal{0, ast.TypeInfo{8, true}},
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
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
			Src: LocalValue{1, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(2),
		},
		ADD{
			Src: TempValue(1),
			Dst: TempValue(2),
		},
		MOV{
			Src: TempValue(2),
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
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
			Dst: TempValue(1),
		},
		ADD{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{1, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(2),
		},
		SUB{
			Src: IntLiteral(1),
			Dst: TempValue(2),
		},
		MOV{
			Src: TempValue(2),
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
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
				Label: "if0else",
				Src:   FuncArg{1, ast.TypeInfo{8, true}, false},
				Dst:   IntLiteral(0),
			},
		},
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, true}, false},
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
		JMP{"if0elsedone"},
		Label("if0else"),
		Label("if0elsedone"),
		MOV{
			Src: FuncArg{0, ast.TypeInfo{8, true}, false},
			Dst: TempValue(1),
		},
		ADD{
			Src: FuncArg{1, ast.TypeInfo{8, true}, false},
			Dst: TempValue(1),
		},
		MOV{
			Src: FuncArg{1, ast.TypeInfo{8, true}, false},
			Dst: TempValue(2),
		},
		SUB{
			Src: IntLiteral(1),
			Dst: TempValue(2),
		},
		CALL{
			FName: "partial_sum",
			Args: []Register{
				TempValue(1),
				TempValue(2),
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
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
		Label("loop0cond"),
		JE{ConditionalJump{Label: "loop0end", Src: LocalValue{0, ast.TypeInfo{1, false}}, Dst: IntLiteral(1)}},
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(15), Dst: TempValue(1)},
		JNE{ConditionalJump{Label: "if1else", Src: TempValue(1), Dst: IntLiteral(0)}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(8), StringLiteral(`fizzbuzz`)}},
		JMP{"if1elsedone"},
		Label("if1else"),
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(5), Dst: TempValue(3)},
		JNE{ConditionalJump{Label: "if2else", Src: TempValue(3), Dst: IntLiteral(0)}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(4), StringLiteral(`buzz`)}},
		JMP{"if2elsedone"},
		Label("if2else"),
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(3), Dst: TempValue(5)},
		JNE{ConditionalJump{Label: "if3else", Src: TempValue(5), Dst: IntLiteral(0)}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(4), StringLiteral(`fizz`)}},
		JMP{"if3elsedone"},
		Label("if3else"),
		CALL{FName: "PrintInt", Args: []Register{LocalValue{1, ast.TypeInfo{8, true}}}},
		Label("if3elsedone"),
		Label("if2elsedone"),
		Label("if1elsedone"),
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		MOV{Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: TempValue(7)},
		ADD{Src: IntLiteral(1), Dst: TempValue(7)},
		MOV{Src: TempValue(7), Dst: LocalValue{1, ast.TypeInfo{8, true}}},
		JL{ConditionalJump{Label: Label("if4else"), Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: IntLiteral(100)}},
		MOV{Src: IntLiteral(1), Dst: LocalValue{0, ast.TypeInfo{1, false}}},
		JMP{"if4elsedone"},
		Label("if4else"),
		Label("if4elsedone"),
		JMP{"loop0cond"},
		Label("loop0end"),
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
		MUL{
			Left:  IntLiteral(2),
			Right: IntLiteral(3),
			Dst:   TempValue(4),
		},
		DIV{
			Left:  IntLiteral(4),
			Right: IntLiteral(2),
			Dst:   TempValue(5),
		},
		MOV{
			Src: TempValue(4),
			Dst: TempValue(6),
		},
		SUB{
			Src: TempValue(5),
			Dst: TempValue(6),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: TempValue(7),
		},
		ADD{
			Src: TempValue(6),
			Dst: TempValue(7),
		},
		MOV{
			Src: TempValue(7),
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(5),
			StringLiteral(`Add: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{0, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(5),
			StringLiteral(`Sub: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{1, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(5),
			StringLiteral(`Mul: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{2, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(5),
			StringLiteral(`Div: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{3, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},

		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(9),
			StringLiteral(`Complex: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{4, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body (%d != %d): got %v want %v\n", len(i.Body), len(expected), i.Body, expected)
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
	branchNum = 0
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
				IntLiteral(1),
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
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
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
	as, ti, c, err := ast.Parse(sampleprograms.EnumType)
	if err != nil {
		t.Fatal(err)
	}

	_, enums, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(enums) == 0 {
		t.Fatalf("No enums returned from %v", as[0])
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
				IntLiteral(8),
				StringLiteral(`I am A!\n`),
			},
		},
		JMP{"match0done"},

		Label("match0v1"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				IntLiteral(8),
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
	branchNum = 0
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
			Label: "match0v0",
			Src:   LocalValue{0, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(0),
		},
		},
		JE{ConditionalJump{
			Label: "match0v1",
			Src:   LocalValue{0, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(1),
		},
		},
		JMP{"match0done"},
		Label("match0v0"),
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(14),
			StringLiteral(`I am nothing!\n`),
		},
		},
		JMP{"match0done"},
		Label("match0v1"),
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{1, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},

		JMP{"match0done"},
		Label("match0done"),
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
			Label: "match1v0",
			Src:   LocalValue{2, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(0),
		},
		},
		JE{ConditionalJump{
			Label: "match1v1",
			Src:   LocalValue{2, ast.TypeInfo{8, false}},
			Dst:   IntLiteral(1),
		},
		},
		JMP{"match1done"},
		Label("match1v0"),
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(14),
			StringLiteral(`I am nothing!\n`),
		},
		},
		JMP{"match1done"},
		Label("match1v1"),
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{3, ast.TypeInfo{8, true}},
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},
		JMP{"match1done"},
		Label("match1done"),
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
	branchNum = 0
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
			Src:   FuncArg{0, ast.TypeInfo{8, false}, false},
			Dst:   IntLiteral(1),
		},
		},
		JE{ConditionalJump{
			Label: "match0v1",
			Src:   FuncArg{0, ast.TypeInfo{8, false}, false},
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
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
func TestIRMatchParam2(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.MatchParam2)
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
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral("x"),
		},
		},
		JE{ConditionalJump{
			Label: "match0v0",
			Src:   FuncArg{0, ast.TypeInfo{8, false}, false},
			Dst:   IntLiteral(1),
		},
		},
		JE{ConditionalJump{
			Label: "match0v1",
			Src:   FuncArg{0, ast.TypeInfo{8, false}, false},
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
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
	branchNum = 0
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
		MOD{Left: LocalValue{1, ast.TypeInfo{8, true}}, Right: IntLiteral(2), Dst: TempValue(2)},
		JNE{
			ConditionalJump{Label: Label("if1else"), Src: TempValue(2), Dst: IntLiteral(0)},
		},
		MUL{
			Left:  LocalValue{1, ast.TypeInfo{8, true}},
			Right: IntLiteral(2),
			Dst:   TempValue(4),
		},
		MOV{Src: LocalValue{0, ast.TypeInfo{8, true}}, Dst: TempValue(5)},
		ADD{Src: TempValue(4), Dst: TempValue(5)},
		MOV{Src: TempValue(5), Dst: LocalValue{0, ast.TypeInfo{8, true}}},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		MOV{Src: LocalValue{1, ast.TypeInfo{8, true}}, Dst: TempValue(6)},
		ADD{Src: IntLiteral(1), Dst: TempValue(6)},
		MOV{Src: TempValue(6), Dst: LocalValue{1, ast.TypeInfo{8, true}}},
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
		CALL{FName: "PrintInt", Args: []Register{
			FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatal(err)
	}

}

func TestIRSimpleArray(t *testing.T) {
	branchNum = 0
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
					Scale:  8,
					Offset: IntLiteral(3),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRArrayMutation(t *testing.T) {
	branchNum = 0
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
					Scale:  8,
					Offset: IntLiteral(3),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		MOV{
			Src: IntLiteral(2),
			Dst: Offset{
				Base:   LocalValue{0, ast.TypeInfo{8, true}},
				Offset: IntLiteral(3),
				Scale:  8,
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Scale:  8,
					Offset: IntLiteral(3),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Scale:  8,
					Offset: IntLiteral(2),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRReferenceVariable(t *testing.T) {
	branchNum = 0
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
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
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
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{1, ast.TypeInfo{8, true}}}},
	}
	i, _, err = Generate(as[1], ti, c, nil)

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSimpleSlice(t *testing.T) {
	branchNum = 0
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
					Offset: IntLiteral(3),
					Scale:  8,
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
	branchNum = 0
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
					Offset: IntLiteral(3),
					Scale:  8,
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
	branchNum = 0
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
					Scale:  8,
					Offset: IntLiteral(3),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		MOV{
			Src: IntLiteral(2),
			Dst: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Offset: IntLiteral(3),
				Scale:  8,
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Scale:  8,
					Offset: IntLiteral(3),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Scale:  8,
					Offset: IntLiteral(2),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSliceParam(t *testing.T) {
	branchNum = 0
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
	branchNum = 0
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
				IntLiteral(7),
				StringLiteral("Stdout!"),
			},
		},
		CALL{
			FName: "Write",
			Args: []Register{
				IntLiteral(2),
				IntLiteral(7),
				StringLiteral("Stderr!"),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRReadSyscall(t *testing.T) {
	branchNum = 0
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
				IntLiteral(7),
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
	branchNum = 0
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
		JL{
			ConditionalJump{
				"match0v0",
				LocalValue{0, ast.TypeInfo{8, true}},
				IntLiteral(3),
			},
		},
		JG{
			ConditionalJump{
				"match0v1",
				LocalValue{0, ast.TypeInfo{8, true}},
				IntLiteral(3),
			},
		},
		JL{
			ConditionalJump{
				"match0v2",
				LocalValue{0, ast.TypeInfo{8, true}},
				IntLiteral(4),
			},
		},
		JMP{"match0done"},
		Label("match0v0"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				IntLiteral(17),
				StringLiteral(`x is less than 3\n`),
			},
		},
		JMP{"match0done"},
		Label("match0v1"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				IntLiteral(20),
				StringLiteral(`x is greater than 3\n`),
			},
		},
		JMP{"match0done"},
		Label("match0v2"),
		CALL{
			FName: "PrintString",
			Args: []Register{
				IntLiteral(17),
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

func TestIREcho(t *testing.T) {
	branchNum = 0
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
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Offset: LocalValue{0, ast.TypeInfo{8, true}},
					Scale:  16,
					Base:   FuncArg{1, ast.TypeInfo{8, false}, false},
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
		CALL{
			FName: "PrintString",
			Args: []Register{
				IntLiteral(1),
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
				IntLiteral(1),
				StringLiteral(`\n`),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestArrayIndex(t *testing.T) {
	branchNum = 0
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
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, true}},
					Scale:  8,
					Offset: LocalValue{0, ast.TypeInfo{8, true}},
				},
			},
		},
		CALL{
			FName: "PrintString",
			Args:  []Register{IntLiteral(1), StringLiteral(`\n`)},
		},
		// Convert x+1 offset from index into byte offset
		MOV{
			Src: LocalValue{0, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{6, ast.TypeInfo{8, true}},
					Scale:  8,
					Offset: TempValue(0),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIndexAssignment(t *testing.T) {
	branchNum = 0
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
				Scale:  8,
				Offset: IntLiteral(1),
			},
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Scale:  8,
				Offset: IntLiteral(2),
			},
			Dst: LocalValue{5, ast.TypeInfo{8, true}},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue{4, ast.TypeInfo{8, true}},
		}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{5, ast.TypeInfo{8, true}}}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIndexedAddition(t *testing.T) {
	branchNum = 0
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
				Scale:  8,
				Offset: IntLiteral(1),
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
				Scale:  8,
				Offset: IntLiteral(2),
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
				Scale:  8,
				Offset: IntLiteral(2),
			},
			Dst: TempValue(1),
		},
		ADD{
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{8, true}},
				Scale:  8,
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
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue{5, ast.TypeInfo{8, true}}}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestStringArray(t *testing.T) {
	branchNum = 0
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
			Src: IntLiteral(3),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue{3, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, false}},
					Scale:  16,
					Offset: IntLiteral(1),
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, false}},
					Scale:  16,
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
	branchNum = 0
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
			Src: IntLiteral(3),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{3, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue{4, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{5, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("baz"),
			Dst: LocalValue{6, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{7, ast.TypeInfo{8, true}},
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
			Dst: LocalValue{8, ast.TypeInfo{8, false}},
		},
		Label("loop0cond"),
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{7, ast.TypeInfo{8, true}}, Dst: LocalValue{8, ast.TypeInfo{8, false}}},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{8, false}},
					Scale:  16,
					Offset: LocalValue{7, ast.TypeInfo{8, true}},
				},
			},
		},
		MOV{
			Src: LocalValue{7, ast.TypeInfo{8, true}},
			Dst: TempValue(1),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue{7, ast.TypeInfo{8, true}},
		},
		JE{
			ConditionalJump{Label: Label("if1else"), Src: LocalValue{7, ast.TypeInfo{8, true}}, Dst: LocalValue{8, ast.TypeInfo{8, false}}},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(" ")}},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		JMP{"loop0cond"},
		Label("loop0end"),
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPreEcho2(t *testing.T) {
	branchNum = 0
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
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   FuncArg{1, ast.TypeInfo{8, false}, false},
					Scale:  16,
					Offset: LocalValue{0, ast.TypeInfo{8, true}},
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
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(" ")}},
		JMP{"if1elsedone"},
		Label("if1else"),
		Label("if1elsedone"),
		JMP{"loop0cond"},
		Label("loop0end"),
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
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
			Src: IntLiteral(3),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue{2, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{3, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue{4, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{5, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("baz"),
			Dst: LocalValue{6, ast.TypeInfo{8, false}},
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

func TestUnbufferedCat(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.UnbufferedCat)
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
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{1, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
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
			Dst: LocalValue{3, ast.TypeInfo{8, false}},
		},
		Label("loop0cond"),
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{2, ast.TypeInfo{8, true}}, Dst: LocalValue{3, ast.TypeInfo{8, false}}},
		},
		CALL{
			FName: "Open",
			Args: []Register{
				Offset{
					Base:   FuncArg{1, ast.TypeInfo{8, false}, false},
					Scale:  16,
					Offset: LocalValue{2, ast.TypeInfo{8, true}},
				},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{4, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "Read",
			Args: []Register{
				LocalValue{4, ast.TypeInfo{8, false}},
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{5, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
		Label("loop1cond"),
		JLE{
			ConditionalJump{
				"loop1end",
				LocalValue{5, ast.TypeInfo{8, false}},
				IntLiteral(0),
			},
		},
		CALL{
			FName: "Read",
			Args: []Register{
				LocalValue{4, ast.TypeInfo{8, false}},
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{5, ast.TypeInfo{8, false}},
		},
		JLE{
			ConditionalJump{
				"if2else",
				LocalValue{5, ast.TypeInfo{8, false}},
				IntLiteral(0),
			},
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
		JMP{"if2elsedone"},
		Label("if2else"),
		Label("if2elsedone"),
		JMP{"loop1cond"},
		Label("loop1end"),
		CALL{
			FName: "Close",
			Args: []Register{
				LocalValue{4, ast.TypeInfo{8, false}},
			},
		},
		MOV{
			Src: LocalValue{2, ast.TypeInfo{8, true}},
			Dst: TempValue(3),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(3),
		},
		MOV{
			Src: TempValue(3),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		JMP{"loop0cond"},
		Label("loop0end"),
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUnbufferedCat2(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.UnbufferedCat2)
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
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{1, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue{2, ast.TypeInfo{8, true}},
		},
		MOV{
			Src: LocalValue{2, ast.TypeInfo{8, true}},
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		Label("loop0cond"),
		MOV{
			Src: LocalValue{3, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue{3, ast.TypeInfo{8, true}},
		},
		CALL{
			FName: "len",
			Args: []Register{
				FuncArg{0, ast.TypeInfo{8, false}, false},
				FuncArg{1, ast.TypeInfo{8, false}, false},
			},
		},
		JGE{
			ConditionalJump{Label: Label("loop0end"), Src: LocalValue{3, ast.TypeInfo{8, true}}, Dst: FuncRetVal{0, ast.TypeInfo{8, false}}},
		},
		CALL{
			FName: "Open",
			Args: []Register{
				Offset{
					Base:   FuncArg{1, ast.TypeInfo{8, false}, false},
					Scale:  16,
					Offset: LocalValue{3, ast.TypeInfo{8, true}},
				},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{4, ast.TypeInfo{8, false}},
		},
		Label("loop1cond"),
		CALL{
			FName: "Read",
			Args: []Register{
				LocalValue{4, ast.TypeInfo{8, false}},
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
		MOV{
			Src: FuncRetVal{0, ast.TypeInfo{8, false}},
			Dst: LocalValue{5, ast.TypeInfo{8, false}},
		},
		JLE{
			ConditionalJump{
				"loop1end",
				LocalValue{5, ast.TypeInfo{8, false}},
				IntLiteral(0),
			},
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
		JMP{"loop1cond"},
		Label("loop1end"),
		CALL{
			FName: "Close",
			Args: []Register{
				LocalValue{4, ast.TypeInfo{8, false}},
			},
		},
		JMP{"loop0cond"},
		Label("loop0end"),
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestAssignmentToVariableIndex(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.AssignmentToVariableIndex)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{ // mutable x = { 1, 3, 4, 5 }
			Src: IntLiteral(1),
			Dst: LocalValue{0, ast.TypeInfo{8, true}},
		},
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
		MOV{ // let y = x[0]
			Src: Offset{
				Base:   LocalValue{0, ast.TypeInfo{8, true}},
				Offset: IntLiteral(0),
				Scale:  8,
			},
			Dst: LocalValue{4, ast.TypeInfo{8, true}},
		},
		MOV{ // x[y] = 6
			Src: IntLiteral(6),
			Dst: Offset{
				Base:   LocalValue{0, ast.TypeInfo{8, true}},
				Offset: LocalValue{4, ast.TypeInfo{8, true}},
				Scale:  8,
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Offset: LocalValue{4, ast.TypeInfo{8, true}},
					Scale:  8,
				},
			},
		},
		MOV{
			Src: LocalValue{4, ast.TypeInfo{8, true}},
			Dst: TempValue(0),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{0, ast.TypeInfo{8, true}},
					Offset: TempValue(0),
					Scale:  8,
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestAssignmentToSliceVariableIndex(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.AssignmentToSliceVariableIndex)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{ // mutable x []byte = { 1, 3, 4, 5 }
			Src: IntLiteral(4),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue{1, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue{2, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue{3, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue{4, ast.TypeInfo{1, false}},
		},
		MOV{ // let y = x[0]
			Src: Offset{
				Base:   LocalValue{1, ast.TypeInfo{1, false}},
				Offset: IntLiteral(0),
				Scale:  1,
			},
			Dst: LocalValue{5, ast.TypeInfo{1, false}},
		},
		MOV{ // x[y] = 6
			Src: IntLiteral(6),
			Dst: Offset{
				Base:   LocalValue{1, ast.TypeInfo{1, false}},
				Offset: LocalValue{5, ast.TypeInfo{1, false}},
				Scale:  1,
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{1, false}},
					Offset: LocalValue{5, ast.TypeInfo{1, false}},
					Scale:  1,
				},
			},
		},
		MOV{
			Src: LocalValue{5, ast.TypeInfo{1, false}},
			Dst: TempValue(0),
		},
		ADD{
			Src: IntLiteral(1),
			Dst: TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue{1, ast.TypeInfo{1, false}},
					Offset: TempValue(0),
					Scale:  1,
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestStringArg(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.StringArg)
	if err != nil {
		t.Fatal(err)
	}
	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{ // mutable x []byte = { 1, 3, 4, 5 }
			Src: IntLiteral(6),
			Dst: LocalValue{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("foobar"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintAString",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				LocalValue{1, ast.TypeInfo{8, false}},
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
			FName: "PrintString",
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

func TestCastBuiltin(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.CastBuiltin)
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
			Src: IntLiteral(70),
			Dst: LocalValue{1, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(111),
			Dst: LocalValue{2, ast.TypeInfo{1, false}},
		},
		MOV{
			Src: IntLiteral(111),
			Dst: LocalValue{3, ast.TypeInfo{1, false}},
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				Pointer{LocalValue{1, ast.TypeInfo{1, false}}},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCastBuiltin2(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.CastBuiltin2)
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
			Src: StringLiteral("bar"),
			Dst: LocalValue{1, ast.TypeInfo{8, false}},
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue{0, ast.TypeInfo{8, false}},
				LocalValue{1, ast.TypeInfo{8, false}},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestSumtypeFuncReturn(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SumTypeFuncReturn)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		// if x
		JE{
			ConditionalJump{
				Label: "if0else",
				Src:   FuncArg{0, ast.TypeInfo{1, false}, false},
				Dst:   IntLiteral(0),
			},
		},
		// return 3
		MOV{
			Src: IntLiteral(0),
			Dst: FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal{1, ast.TypeInfo{8, false}},
		},
		RET{},
		JMP{"if0elsedone"},
		Label("if0else"),
		Label("if0elsedone"),
		// return "not3"
		MOV{
			Src: IntLiteral(1),
			Dst: FuncRetVal{0, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: IntLiteral(4),
			Dst: FuncRetVal{1, ast.TypeInfo{8, false}},
		},
		MOV{
			Src: StringLiteral("not3"),
			Dst: FuncRetVal{2, ast.TypeInfo{8, false}},
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
}

func TestIfBool(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.IfBool)
	if err != nil {
		t.Fatal(err)
	}

	i, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		// if x
		JE{
			ConditionalJump{
				Label: "if0else",
				Src:   FuncArg{0, ast.TypeInfo{1, false}, false},
				Dst:   IntLiteral(0),
			},
		},
		// return 3
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
		},
		RET{},
		JMP{"if0elsedone"},
		Label("if0else"),
		Label("if0elsedone"),
		// return 7
		MOV{
			Src: IntLiteral(7),
			Dst: FuncRetVal{0, ast.TypeInfo{8, true}},
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
}
