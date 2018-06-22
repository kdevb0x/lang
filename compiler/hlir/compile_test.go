package hlir

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
	case IF:
		b1, ok := b.(IF)
		if !ok {
			return false
		}
		if a1.Condition.Register != b1.Condition.Register {
			return false
		}
		if err := compareIR(a1.Condition.Body, b1.Condition.Body); err != nil {
			return false
		}
		if err := compareIR(a1.Body, b1.Body); err != nil {
			return false
		}
		if err := compareIR(a1.ElseBody, b1.ElseBody); err != nil {
			return false
		}
		return true
	case LOOP:
		b1, ok := b.(LOOP)
		if !ok {
			return false
		}
		if !compareOp(a1.Condition, b1.Condition) {
			return false
		}
		if err := compareIR(a1.Body, b1.Body); err != nil {
			return false
		}
		if err := compareIR(a1.Initializer, b1.Initializer); err != nil {
			return false
		}
		return true
	case Condition:
		b1, ok := b.(Condition)
		if !ok {
			return false
		}
		if err := compareIR(a1.Body, b1.Body); err != nil {
			return false
		}
		return a1.Register == b1.Register
	case JumpTable:
		b1, ok := b.(JumpTable)
		if !ok {
			return false
		}
		if len(a1) != len(b1) {
			return false
		}
		for i := range a1 {
			if !compareOp(a1[i].Condition, b1[i].Condition) {
				return false
			}
			if err := compareIR(a1[i].Body, b1[i].Body); err != nil {
				return false
			}
		}
		return true
	case ASSERT:
		b1, ok := b.(ASSERT)
		if !ok {
			return false
		}
		if !compareOp(a1.Predicate, b1.Predicate) {
			return false
		}
		return a1.Message == b1.Message
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

	i, _, _, err := Generate(ast[0], ti, c, nil)
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

	i, _, _, err := Generate(ast[0], ti, c, nil)
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(1),
		},
		MOV{
			Src: StringLiteral("hello"),
			Dst: LocalValue(2),
		},
		CALL{FName: "PrintString", Args: []Register{
			LocalValue(1),
			LocalValue(2),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal(0),
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

	i, _, _, err = Generate(as[1], ti, c, nil)
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
				LastFuncCallRetVal{0, 0},
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		CALL{FName: "foo", Args: []Register{}},
		CALL{FName: "PrintInt", Args: []Register{
			LastFuncCallRetVal{0, 0},
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

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected = []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		ADD{
			Left:  LocalValue(0),
			Right: IntLiteral(1),
			Dst:   TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue(1),
		},
		ADD{
			Left:  LocalValue(1),
			Right: IntLiteral(1),
			Dst:   TempValue(1),
		},
		ADD{
			Left:  LocalValue(0),
			Right: TempValue(1),
			Dst:   TempValue(2),
		},
		MOV{
			Src: TempValue(2),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", PrettyPrint(0, i.Body), PrettyPrint(0, expected))
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "foo" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "foo")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: FuncRetVal(0),
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

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "foo", Args: []Register{}},
		CALL{FName: "PrintInt", Args: []Register{
			LastFuncCallRetVal{0, 0},
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "sum")
	}
	expected := []Opcode{
		MOV{
			Src: FuncArg{0, false},
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(1),
		},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					GT{LocalValue(0), IntLiteral(0), TempValue(0)},
				},
				Register: TempValue(0),
			},
			Body: []Opcode{
				ADD{
					Left:  LocalValue(1),
					Right: LocalValue(0),
					Dst:   TempValue(1),
				},
				MOV{
					Src: TempValue(1),
					Dst: LocalValue(1),
				},
				SUB{
					Left:  LocalValue(0),
					Right: IntLiteral(1),
					Dst:   TempValue(2),
				},
				MOV{
					Src: TempValue(2),
					Dst: LocalValue(0),
				},
			},
		},
		MOV{
			Src: LocalValue(1),
			Dst: FuncRetVal(0),
		},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i, _, _, err = Generate(as[1], ti, c, nil)
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
			LastFuncCallRetVal{0, 0},
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

	i, _, _, err := Generate(as[0], ti, c, nil)
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
				FuncArg{0, false},
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

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "partial_sum" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "partial_sum")
	}
	expected = []Opcode{
		IF{
			ControlFlow: ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{1, false},
							Right: IntLiteral(0),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{
					MOV{
						Src: FuncArg{0, false},
						Dst: FuncRetVal(0),
					},
					RET{},
				},
			},
			ElseBody: nil,
		},
		ADD{
			Left:  FuncArg{0, false},
			Right: FuncArg{1, false},
			Dst:   TempValue(1),
		},
		SUB{
			Left:  FuncArg{1, false},
			Right: IntLiteral(1),
			Dst:   TempValue(2),
		},
		CALL{
			FName: "partial_sum",
			Args: []Register{
				TempValue(1), TempValue(2),
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

	i, _, _, err = Generate(as[2], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected = []Opcode{
		CALL{FName: "sum", Args: []Register{IntLiteral(10)}},
		CALL{FName: "PrintInt", Args: []Register{
			LastFuncCallRetVal{0, 0},
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{Src: IntLiteral(0), Dst: LocalValue(0)},
		MOV{Src: IntLiteral(1), Dst: LocalValue(1)},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					NEQ{
						Left:  LocalValue(0),
						Right: IntLiteral(1),
						Dst:   TempValue(0),
					},
				},
				Register: TempValue(0),
			},
			Body: []Opcode{
				IF{
					ControlFlow: ControlFlow{
						Condition: Condition{
							Body: []Opcode{MOD{
								Left:  LocalValue(1),
								Right: IntLiteral(15),
								Dst:   TempValue(1),
							},
								EQ{
									Left:  TempValue(1),
									Right: IntLiteral(0),
									Dst:   TempValue(2),
								},
							},
							Register: TempValue(2),
						},
						Body: []Opcode{
							CALL{FName: "PrintString", Args: []Register{
								IntLiteral(8),
								StringLiteral("fizzbuzz"),
							},
							},
						},
					},
					ElseBody: []Opcode{
						IF{
							ControlFlow: ControlFlow{
								Condition: Condition{
									Body: []Opcode{MOD{
										Left:  LocalValue(1),
										Right: IntLiteral(5),
										Dst:   TempValue(3),
									},
										EQ{
											Left:  TempValue(3),
											Right: IntLiteral(0),
											Dst:   TempValue(4),
										},
									},
									Register: TempValue(4),
								},
								Body: []Opcode{
									CALL{FName: "PrintString", Args: []Register{
										IntLiteral(4),
										StringLiteral("buzz"),
									},
									},
								},
							},
							ElseBody: []Opcode{
								IF{
									ControlFlow: ControlFlow{
										Condition: Condition{
											Body: []Opcode{MOD{
												Left:  LocalValue(1),
												Right: IntLiteral(3),
												Dst:   TempValue(5),
											},
												EQ{
													Left:  TempValue(5),
													Right: IntLiteral(0),
													Dst:   TempValue(6),
												},
											},
											Register: TempValue(6),
										},
										Body: []Opcode{
											CALL{FName: "PrintString", Args: []Register{
												IntLiteral(4),
												StringLiteral("fizz"),
											},
											},
										},
									},
									ElseBody: []Opcode{
										CALL{FName: "PrintInt", Args: []Register{
											LocalValue(1),
										},
										},
									},
								},
							},
						},
					},
				},
				CALL{FName: "PrintString", Args: []Register{
					IntLiteral(1),
					StringLiteral(`\n`),
				},
				},
				ADD{
					Left:  LocalValue(1),
					Right: IntLiteral(1),
					Dst:   TempValue(7),
				},
				MOV{
					Src: TempValue(7),
					Dst: LocalValue(1),
				},
				IF{
					ControlFlow: ControlFlow{
						Condition: Condition{
							Body: []Opcode{
								GEQ{
									Left:  LocalValue(1),
									Right: IntLiteral(100),
									Dst:   TempValue(8),
								},
							},
							Register: TempValue(8),
						},
						Body: []Opcode{
							MOV{
								Src: IntLiteral(1),
								Dst: LocalValue(0),
							},
						},
					},
				},
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

func TestIRGenSomeMathStatement(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SomeMath)
	if err != nil {
		t.Fatal(err)
	}

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		ADD{
			Left:  IntLiteral(1),
			Right: IntLiteral(2),
			Dst:   TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue(0),
		},
		SUB{
			Left:  IntLiteral(1),
			Right: IntLiteral(2),
			Dst:   TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue(1),
		},
		MUL{
			Left:  IntLiteral(2),
			Right: IntLiteral(3),
			Dst:   TempValue(2),
		},
		MOV{
			Src: TempValue(2),
			Dst: LocalValue(2),
		},
		DIV{
			Left:  IntLiteral(6),
			Right: IntLiteral(2),
			Dst:   TempValue(3),
		},
		MOV{
			Src: TempValue(3),
			Dst: LocalValue(3),
		},

		// x
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
		SUB{
			Left:  TempValue(4),
			Right: TempValue(5),
			Dst:   TempValue(6),
		},
		ADD{
			Left:  IntLiteral(1),
			Right: TempValue(6),
			Dst:   TempValue(7),
		},
		MOV{
			Src: TempValue(7),
			Dst: LocalValue(4),
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(5),
			StringLiteral(`Add: `),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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
			LocalValue(1),
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
			LocalValue(2),
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
			LocalValue(3),
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
			LocalValue(4),
		},
		},
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral(`\n`),
		},
		},
	}
	if len(i.Body) != len(expected) {
		t.Fatalf("Unexpected body: got %v want %v\n", PrettyPrint(0, i.Body), PrettyPrint(0, expected))
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

	i, _, _, err := Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(-4),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(0),
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
	as, ti, c, err := ast.Parse(sampleprograms.Fibonacci)
	if err != nil {
		t.Fatal(err)
	}

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "fib_rec" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		ADD{
			Left:  FuncArg{0, false},
			Right: FuncArg{1, false},
			Dst:   TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue(0),
		},
		IF{
			ControlFlow: ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						GEQ{
							Left:  LocalValue(0),
							Right: IntLiteral(200),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					MOV{
						Src: FuncArg{1, false},
						Dst: FuncRetVal(0),
					},
					RET{},
				},
			},
			ElseBody: nil,
		},
		CALL{
			FName: "PrintInt",
			Args:  []Register{LocalValue(0)},
		},
		CALL{
			FName: "PrintString",
			Args:  []Register{IntLiteral(1), StringLiteral(`\n`)},
		},
		CALL{
			FName: "fib_rec",
			Args: []Register{
				FuncArg{1, false},
				LocalValue(0),
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

	i, _, _, err = Generate(as[1], ti, c, nil)
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
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(0),
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

	_, enums, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	if i.Name != "main" {
		t.Errorf("Unexpected name: got %v want %v", i.Name, "main")
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(0),
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(0),
							Right: IntLiteral(0),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(8),
							StringLiteral(`I am A!\n`),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(0),
							Right: IntLiteral(1),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(8),
							StringLiteral(`I am B!\n`),
						},
					},
				},
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

func TestIRGenericEnumType(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.GenericEnumType)
	if err != nil {
		t.Fatal(err)
	}
	_, enums, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		IF{
			ControlFlow: ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						GT{
							Left:  FuncArg{0, false},
							Right: IntLiteral(3),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{
					MOV{
						Src: IntLiteral(0),
						Dst: FuncRetVal(0),
					},
					RET{},
				},
			},
		},
		MOV{
			Src: IntLiteral(1),
			Dst: FuncRetVal(0),
		},
		// The concrete parameter is an int, which goes into the
		// next word.
		MOV{
			Src: IntLiteral(5),
			Dst: FuncRetVal(1),
		},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i, _, rd, err := Generate(as[2], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		CALL{FName: "DoSomething", Args: []Register{
			IntLiteral(3),
		},
		},
		MOV{
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(0),
		},
		MOV{
			Src: LastFuncCallRetVal{0, 1},
			Dst: LocalValue(1),
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(0),
							Right: IntLiteral(0),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(14),
							StringLiteral(`I am nothing!\n`),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(0),
							Right: IntLiteral(1),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							LocalValue(1),
						},
					},
					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(1),
							StringLiteral(`\n`),
						},
					},
				},
			},
		},
		CALL{FName: "DoSomething", Args: []Register{
			IntLiteral(4),
		},
		},
		MOV{
			Src: LastFuncCallRetVal{4, 0},
			Dst: LocalValue(2),
		},
		MOV{
			Src: LastFuncCallRetVal{4, 1},
			Dst: LocalValue(3),
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(2),
							Right: IntLiteral(0),
							Dst:   TempValue(2),
						},
					},
					Register: TempValue(2),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(14),
							StringLiteral(`I am nothing!\n`),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(2),
							Right: IntLiteral(1),
							Dst:   TempValue(3),
						},
					},
					Register: TempValue(3),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							LocalValue(3),
						},
					},
					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(1),
							StringLiteral(`\n`),
						},
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	lv1, ok := rd[LocalValue(1)]
	if !ok {
		t.Fatal("No type information for LocalValue(1)")
	}
	if ti := lv1.TypeInfo; ti != (ast.TypeInfo{0, true}) {
		t.Fatalf("Unexpected type info for LocalValue(1): got %v want %v", ti, ast.TypeInfo{0, true})
	}
}

func TestIRMatchParam(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.MatchParam)
	if err != nil {
		t.Fatal(err)
	}
	_, enums, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{0, false},
							Right: IntLiteral(1),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{
					MOV{
						Src: FuncArg{1, false},
						Dst: FuncRetVal(0),
					},
					RET{},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{0, false},
							Right: IntLiteral(0),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					MOV{
						Src: IntLiteral(0),
						Dst: FuncRetVal(0),
					},
					RET{},
				},
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
	i, _, _, err = Generate(as[2], ti, c, enums)
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
			LastFuncCallRetVal{0, 0},
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
	_, enums, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[1], ti, c, enums)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		CALL{FName: "PrintString", Args: []Register{
			IntLiteral(1),
			StringLiteral("x"),
		},
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{0, false},
							Right: IntLiteral(1),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{
					MOV{
						Src: FuncArg{1, false},
						Dst: FuncRetVal(0),
					},
					RET{},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{0, false},
							Right: IntLiteral(0),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					MOV{
						Src: IntLiteral(0),
						Dst: FuncRetVal(0),
					},
					RET{},
				},
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
	i, _, _, err = Generate(as[2], ti, c, enums)
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
			LastFuncCallRetVal{0, 0},
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
	as, ti, c, err := ast.Parse(sampleprograms.SimpleAlgorithm)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(1),
		},
		MUL{
			Left:  FuncArg{0, false},
			Right: IntLiteral(2),
			Dst:   TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(1),
		},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					LT{
						Left:  LocalValue(1),
						Right: LocalValue(2),
						Dst:   TempValue(1),
					},
				},
				Register: TempValue(1),
			},
			Body: []Opcode{
				IF{
					ControlFlow: ControlFlow{
						Condition: Condition{
							Body: []Opcode{
								MOD{
									Left:  LocalValue(1),
									Right: IntLiteral(2),
									Dst:   TempValue(2),
								},
								EQ{
									Left:  TempValue(2),
									Right: IntLiteral(0),
									Dst:   TempValue(3),
								},
							},
							Register: TempValue(3),
						},
						Body: []Opcode{
							MUL{
								Left:  LocalValue(1),
								Right: IntLiteral(2),
								Dst:   TempValue(4),
							},
							ADD{
								Left:  LocalValue(0),
								Right: TempValue(4),
								Dst:   TempValue(5),
							},
							MOV{
								Src: TempValue(5),
								Dst: LocalValue(0),
							},
						},
					},
				},
				ADD{
					Left:  LocalValue(1),
					Right: IntLiteral(1),
					Dst:   TempValue(6),
				},
				MOV{
					Src: TempValue(6),
					Dst: LocalValue(1),
				},
			},
		},
		MOV{
			Src: LocalValue(0),
			Dst: FuncRetVal(0),
		},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected = []Opcode{
		CALL{FName: "loop", Args: []Register{
			IntLiteral(10),
		},
		},
		CALL{FName: "PrintInt", Args: []Register{
			LastFuncCallRetVal{0, 0},
		},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatal(err)
	}

}

func TestIRSimpleArray(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SimpleArray)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(4),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.ArrayType{
							ast.TypeLiteral("int"),
							5,
						},
						false,
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v: %v", err, i.Body)
	}
}

func TestIRArrayMutation(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.ArrayMutation)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(4),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.ArrayType{
							ast.TypeLiteral("int"),
							5,
						},
						false,
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		MOV{
			Src: IntLiteral(2),
			Dst: Offset{
				Base:   LocalValue(0),
				Offset: IntLiteral(3),
				Scale:  0,
				Container: ast.VarWithType{
					"n", ast.ArrayType{
						ast.TypeLiteral("int"),
						5,
					},
					false,
				},
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.ArrayType{
							ast.TypeLiteral("int"),
							5,
						},
						false,
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: IntLiteral(2),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.ArrayType{
							ast.TypeLiteral("int"),
							5,
						},
						false,
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRReferenceVariable(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.ReferenceVariable)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(4),
			Dst: FuncArg{0, true},
		},
		ADD{
			Left:  FuncArg{0, true},
			Right: FuncArg{1, false},
			Dst:   TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: FuncRetVal(0),
		},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	expected = []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		CALL{FName: "PrintInt", Args: []Register{LocalValue(0)}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "changer",
			Args: []Register{
				Pointer{LocalValue(0)},
				IntLiteral(3),
			},
		},
		MOV{
			Src: LastFuncCallRetVal{2, 0},
			Dst: LocalValue(1),
		},
		CALL{FName: "PrintInt", Args: []Register{LocalValue(0)}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue(1)}},
	}
	i, _, _, err = Generate(as[1], ti, c, nil)

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSimpleSlice(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SimpleSlice)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	// LV0 == size of the slice, 1-5=the values
	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(5),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Base:   LocalValue(1),
					Container: ast.VarWithType{
						"n", ast.SliceType{
							ast.TypeLiteral("int"),
						},
						false,
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v: %v", i.Body, err)
	}
}

func TestIRSimpleSliceInference(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SimpleSliceInference)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(5),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Base:   LocalValue(1),
					Container: ast.VarWithType{
						"n2", ast.SliceType{
							ast.TypeLiteral("int"),
						},
						false,
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSimpleSliceMutation(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SliceMutation)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(5),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(1),
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.SliceType{
							ast.TypeLiteral("int"),
						},
						false,
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		MOV{
			Src: IntLiteral(2),
			Dst: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(3),
				Scale:  0,
				Container: ast.VarWithType{
					"n", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(1),
					Offset: IntLiteral(3),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.SliceType{
							ast.TypeLiteral("int"),
						},
						false,
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(1),
					Offset: IntLiteral(2),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.SliceType{
							ast.TypeLiteral("int"),
						},
						false,
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRSliceParam(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SliceParam)
	if err != nil {
		t.Fatal(err)
	}

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(44),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(55),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(88),
			Dst: LocalValue(3),
		},
		CALL{
			FName: "PrintASlice",
			Args: []Register{
				LocalValue(0),
				Pointer{LocalValue(1)},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				FuncArg{0, false},
				FuncArg{1, false},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRWriteSyscall(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.WriteSyscall)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
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
	as, ti, c, err := ast.Parse(sampleprograms.ReadSyscall)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
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
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(6),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(5),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(6),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(7),
		},
		CALL{
			FName: "Read",
			Args: []Register{
				LocalValue(0),
				LocalValue(1),
				Pointer{LocalValue(2)},
			},
		},
		MOV{
			Src: LastFuncCallRetVal{1, 0},
			Dst: LocalValue(8),
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue(1),
				Pointer{LocalValue(2)},
			},
		},
		CALL{
			FName: "Close",
			Args: []Register{
				LocalValue(0),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIRIfElseMatch(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.IfElseMatch)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						LT{
							Left:  LocalValue(0),
							Right: IntLiteral(3),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(17),
							StringLiteral(`x is less than 3\n`),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						GT{
							Left:  LocalValue(0),
							Right: IntLiteral(3),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(20),
							StringLiteral(`x is greater than 3\n`),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						LT{
							Left:  LocalValue(0),
							Right: IntLiteral(4),
							Dst:   TempValue(2),
						},
					},
					Register: TempValue(2),
				},
				Body: []Opcode{

					CALL{
						FName: "PrintString",
						Args: []Register{
							IntLiteral(17),
							StringLiteral(`x is less than 4\n`),
						},
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIREcho(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.Echo)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		CALL{
			FName: "len",
			Args: []Register{
				FuncArg{0, false},
				FuncArg{1, false},
			},
		},
		MOV{
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(1),
		},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					LT{
						Left:  LocalValue(0),
						Right: LocalValue(1),
						Dst:   TempValue(0),
					},
				},
				Register: TempValue(0),
			},
			Body: []Opcode{
				CALL{
					FName: "PrintString",
					Args: []Register{
						Offset{
							Base:   FuncArg{1, false},
							Offset: LocalValue(0),
							Scale:  IntLiteral(16),
							Container: ast.VarWithType{
								"args", ast.SliceType{
									ast.TypeLiteral("string"),
								},
								false,
							},
						},
					},
				},
				ADD{
					Left:  LocalValue(0),
					Right: IntLiteral(1),
					Dst:   TempValue(1),
				},
				MOV{
					Src: TempValue(1),
					Dst: LocalValue(0),
				},
				IF{
					ControlFlow: ControlFlow{
						Condition: Condition{
							Body: []Opcode{
								NEQ{
									Left:  LocalValue(0),
									Right: LocalValue(1),
									Dst:   TempValue(2),
								},
							},
							Register: TempValue(2),
						},
						Body: []Opcode{
							CALL{
								FName: "PrintString",
								Args: []Register{
									IntLiteral(1),
									StringLiteral(" "),
								},
							},
						},
					},
				},
			},
		},
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
	as, ti, c, err := ast.Parse(sampleprograms.ArrayIndex)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		// let x = 3
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		// Let statement
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(5),
		},
		// mutable statement
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(6),
		},
		MOV{
			Src: IntLiteral(2),
			Dst: LocalValue(7),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(8),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(9),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(10),
		},
		// Convert index from index into byte offset
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(1),
					Offset: LocalValue(0),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n", ast.ArrayType{
							Base: ast.TypeLiteral("int"),
							Size: 5,
						},
						false,
					},
				},
			},
		},
		CALL{
			FName: "PrintString",
			Args:  []Register{IntLiteral(1), StringLiteral(`\n`)},
		},
		// Convert x+1 offset from index into byte offset
		ADD{
			Left:  LocalValue(0),
			Right: IntLiteral(1),
			Dst:   TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(6),
					Offset: TempValue(0),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"n2", ast.ArrayType{
							Base: ast.TypeLiteral("int"),
							Size: 5,
						},
						false,
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIndexAssignment(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.IndexAssignment)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		// let x []int = { 3, 4, 5 }
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		// Let statement
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(3),
		},
		MOV{
			Src: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(1),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
			Dst: LocalValue(4),
		},
		MOV{
			Src: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(2),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
			Dst: LocalValue(5),
		},
		CALL{FName: "PrintInt", Args: []Register{
			LocalValue(4),
		}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue(5)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestIndexedAddition(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.IndexedAddition)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		// let x []int = { 3, 4, 5 }
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		// Let statement
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(3),
		},
		MOV{
			Src: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(1),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
			Dst: LocalValue(4),
		},
		ADD{
			Left: LocalValue(4),
			Right: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(2),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
			Dst: TempValue(0),
		},
		MOV{
			Src: TempValue(0),
			Dst: LocalValue(4),
		},
		ADD{
			Left: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(2),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
			Right: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(0),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						ast.TypeLiteral("int"),
					},
					false,
				},
			},
			Dst: TempValue(1),
		},
		MOV{
			Src: TempValue(1),
			Dst: LocalValue(5),
		},
		CALL{FName: "PrintInt", Args: []Register{LocalValue(4)}},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{FName: "PrintInt", Args: []Register{LocalValue(5)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestStringArray(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.StringArray)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(2),
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue(3),
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: IntLiteral(1),
					Scale:  IntLiteral(16),
					Container: ast.VarWithType{
						"args", ast.ArrayType{
							ast.TypeLiteral("string"),
							2,
						},
						false,
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
		CALL{
			FName: "PrintString",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: IntLiteral(0),
					Scale:  IntLiteral(16),
					Container: ast.VarWithType{
						"args", ast.ArrayType{
							ast.TypeLiteral("string"),
							2,
						},
						false,
					},
				},
			},
		},
	}
	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPreEcho(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.PreEcho)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(1),
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(3),
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(5),
		},
		MOV{
			Src: StringLiteral("baz"),
			Dst: LocalValue(6),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(7),
		},
		CALL{
			FName: "len",
			Args: []Register{
				LocalValue(0),
				Pointer{LocalValue(1)},
			},
		},
		MOV{
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(8),
		},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					LT{
						Left:  LocalValue(7),
						Right: LocalValue(8),
						Dst:   TempValue(0),
					},
				},
				Register: TempValue(0),
			},
			Body: []Opcode{
				CALL{
					FName: "PrintString",
					Args: []Register{
						Offset{
							Base:   LocalValue(1),
							Scale:  IntLiteral(16),
							Offset: LocalValue(7),
							Container: ast.VarWithType{
								"args",
								ast.SliceType{
									ast.TypeLiteral("string"),
								},
								false,
							},
						},
					},
				},
				ADD{
					Left:  LocalValue(7),
					Right: IntLiteral(1),
					Dst:   TempValue(1),
				},
				MOV{
					Src: TempValue(1),
					Dst: LocalValue(7),
				},
				IF{
					ControlFlow: ControlFlow{
						Condition: Condition{
							Body: []Opcode{
								NEQ{
									Left:  LocalValue(7),
									Right: LocalValue(8),
									Dst:   TempValue(2),
								},
							},
							Register: TempValue(2),
						},
						Body: []Opcode{

							CALL{
								FName: "PrintString",
								Args: []Register{
									IntLiteral(1),
									StringLiteral(" "),
								},
							},
						},
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPreEcho2(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.PreEcho2)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		CALL{
			FName: "len",
			Args: []Register{
				FuncArg{0, false},
				FuncArg{1, false},
			},
		},
		MOV{
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(1),
		},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					LT{
						Left:  LocalValue(0),
						Right: LocalValue(1),
						Dst:   TempValue(0),
					},
				},
				Register: TempValue(0),
			},
			Body: []Opcode{
				CALL{
					FName: "PrintString",
					Args: []Register{
						Offset{
							Base:   FuncArg{1, false},
							Scale:  IntLiteral(16),
							Offset: LocalValue(0),
							Container: ast.VarWithType{
								"args",
								ast.SliceType{
									ast.TypeLiteral("string"),
								},
								false,
							},
						},
					},
				},
				ADD{
					Left:  LocalValue(0),
					Right: IntLiteral(1),
					Dst:   TempValue(1),
				},
				MOV{
					Src: TempValue(1),
					Dst: LocalValue(0),
				},
				IF{
					ControlFlow: ControlFlow{
						Condition: Condition{
							Body: []Opcode{
								NEQ{
									Left:  LocalValue(0),
									Right: LocalValue(1),
									Dst:   TempValue(2),
								},
							},
							Register: TempValue(2),
						},
						Body: []Opcode{

							CALL{
								FName: "PrintString",
								Args: []Register{
									IntLiteral(1),
									StringLiteral(" "),
								},
							},
						},
					},
				},
			},
		},
		CALL{FName: "PrintString", Args: []Register{IntLiteral(1), StringLiteral(`\n`)}},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	i2, _, _, err := Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(1),
		},
		MOV{
			Src: StringLiteral("foo"),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(3),
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue(4),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(5),
		},
		MOV{
			Src: StringLiteral("baz"),
			Dst: LocalValue(6),
		},
		CALL{
			FName: "PrintSlice",
			Args: []Register{
				LocalValue(0),
				Pointer{LocalValue(1)},
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
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(2),
		},
		CALL{
			FName: "len",
			Args: []Register{
				FuncArg{0, false},
				FuncArg{1, false},
			},
		},
		MOV{
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(3),
		},
		LOOP{
			Condition: Condition{
				Body: []Opcode{
					LT{Left: LocalValue(2), Right: LocalValue(3), Dst: TempValue(0)},
				},
				Register: TempValue(0),
			},
			Body: []Opcode{
				CALL{
					FName: "Open",
					Args: []Register{
						Offset{
							Base:   FuncArg{1, false},
							Offset: LocalValue(2),
							Scale:  IntLiteral(16),
							Container: ast.VarWithType{
								"args",
								ast.SliceType{
									ast.TypeLiteral("string"),
								},
								false,
							},
						},
					},
				},
				MOV{
					Src: LastFuncCallRetVal{1, 0},
					Dst: LocalValue(4),
				},
				CALL{
					FName: "Read",
					Args: []Register{
						LocalValue(4),
						LocalValue(0),
						Pointer{LocalValue(1)},
					},
				},
				MOV{
					Src: LastFuncCallRetVal{2, 0},
					Dst: LocalValue(5),
				},
				CALL{
					FName: "PrintByteSlice",
					Args: []Register{
						LocalValue(0),
						Pointer{LocalValue(1)},
					},
				},
				LOOP{
					Condition: Condition{
						Body: []Opcode{
							GT{Left: LocalValue(5), Right: IntLiteral(0), Dst: TempValue(1)},
						},
						Register: TempValue(1),
					},
					Body: []Opcode{
						CALL{
							FName: "Read",
							Args: []Register{
								LocalValue(4),
								LocalValue(0),
								Pointer{LocalValue(1)},
							},
						},
						MOV{
							Src: LastFuncCallRetVal{4, 0},
							Dst: LocalValue(5),
						},
						IF{
							ControlFlow: ControlFlow{
								Condition: Condition{
									Body: []Opcode{
										GT{Left: LocalValue(5), Right: IntLiteral(0), Dst: TempValue(2)},
									},
									Register: TempValue(2),
								},
								Body: []Opcode{
									CALL{
										FName: "PrintByteSlice",
										Args: []Register{
											LocalValue(0),
											Pointer{LocalValue(1)},
										},
									},
								},
							},
						},
					},
				},
				CALL{
					FName: "Close",
					Args: []Register{
						LocalValue(4),
					},
				},
				ADD{Left: LocalValue(2), Right: IntLiteral(1), Dst: TempValue(3)},
				MOV{
					Src: TempValue(3),
					Dst: LocalValue(2),
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestPrecedence(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.Precedence)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		ADD{
			Left:  IntLiteral(1),
			Right: IntLiteral(2),
			Dst:   TempValue(0),
		},
		SUB{
			Left:  IntLiteral(3),
			Right: IntLiteral(4),
			Dst:   TempValue(1),
		},
		MUL{
			Left:  TempValue(0),
			Right: TempValue(1),
			Dst:   TempValue(2),
		},
		MOV{
			Src: TempValue(2),
			Dst: LocalValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				LocalValue(0),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestLetCondition(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.LetCondition)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(0),
		},
		IF{
			ControlFlow: ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						ADD{
							Left:  LocalValue(0),
							Right: IntLiteral(1),
							Dst:   TempValue(0),
						},
						MOV{
							Src: TempValue(0),
							Dst: LocalValue(1),
						},
						EQ{
							Left:  LocalValue(1),
							Right: IntLiteral(1),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							LocalValue(1),
						},
					},
				},
			},
			ElseBody: []Opcode{
				CALL{
					FName: "PrintInt",
					Args: []Register{
						IntLiteral(-1),
					},
				},
			},
		},
		IF{
			ControlFlow: ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						ADD{
							Left:  LocalValue(0),
							Right: IntLiteral(1),
							Dst:   TempValue(2),
						},
						MOV{
							Src: TempValue(2),
							Dst: LocalValue(2),
						},
						NEQ{
							Left:  LocalValue(2),
							Right: IntLiteral(1),
							Dst:   TempValue(3),
						},
					},
					Register: TempValue(3),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							LocalValue(2),
						},
					},
				},
			},
			ElseBody: []Opcode{
				CALL{
					FName: "PrintInt",
					Args: []Register{
						IntLiteral(-1),
					},
				},
			},
		},
		LOOP{
			Initializer: []Opcode{
				MOV{
					Src: LocalValue(0),
					Dst: LocalValue(3),
				},
			},
			Condition: Condition{
				Body: []Opcode{
					ADD{
						Left:  LocalValue(3),
						Right: IntLiteral(1),
						Dst:   TempValue(4),
					},
					MOV{
						Src: TempValue(4),
						Dst: LocalValue(3),
					},
					LT{LocalValue(3), IntLiteral(3), TempValue(5)},
				},
				Register: TempValue(5),
			},
			Body: []Opcode{
				CALL{
					FName: "PrintInt",
					Args: []Register{
						LocalValue(3),
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestSliceStringVariableParam(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SliceStringVariableParam)

	if err != nil {
		t.Fatal(err)
	}

	f, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}
	if f.NumArgs != 2 {
		t.Errorf("Got %d arguments, expected 2. Slices take up 2 argument slots.", f.NumArgs)
	}
}

func TestUnbufferedCat2(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.UnbufferedCat2)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(2),
		},
		LOOP{
			Initializer: []Opcode{
				MOV{
					Src: LocalValue(2),
					Dst: LocalValue(3),
				},
			},
			Condition: Condition{
				Body: []Opcode{
					ADD{
						Left:  LocalValue(3),
						Right: IntLiteral(1),
						Dst:   TempValue(0),
					},
					MOV{
						Src: TempValue(0),
						Dst: LocalValue(3),
					},
					CALL{
						FName: "len",
						Args: []Register{
							FuncArg{0, false},
							FuncArg{1, false},
						},
					},
					LT{Left: LocalValue(3), Right: LastFuncCallRetVal{0, 0}, Dst: TempValue(1)},
				},
				Register: TempValue(1),
			},
			Body: []Opcode{
				CALL{
					FName: "Open",
					Args: []Register{
						Offset{
							Base:   FuncArg{1, false},
							Offset: LocalValue(3),
							Scale:  IntLiteral(16),
							Container: ast.VarWithType{
								"args",
								ast.SliceType{
									ast.TypeLiteral("string"),
								},
								false,
							},
						},
					},
				},
				MOV{
					Src: LastFuncCallRetVal{1, 0},
					Dst: LocalValue(4),
				},
				LOOP{
					Condition: Condition{
						Body: []Opcode{
							CALL{
								FName: "Read",
								Args: []Register{
									LocalValue(4),
									LocalValue(0),
									Pointer{LocalValue(1)},
								},
							},
							MOV{
								Src: LastFuncCallRetVal{2, 0},
								Dst: LocalValue(5),
							},
							GT{Left: LocalValue(5), Right: IntLiteral(0), Dst: TempValue(2)},
						},
						Register: TempValue(2),
					},
					Body: []Opcode{
						CALL{
							FName: "PrintByteSlice",
							Args: []Register{
								LocalValue(0),
								Pointer{LocalValue(1)},
							},
						},
					},
				},
				CALL{
					FName: "Close",
					Args: []Register{
						LocalValue(4),
					},
				},
			},
		},
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
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{ // mutable x = { 1, 3, 4, 5 }
			Src: IntLiteral(1),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(3),
		},
		MOV{ // let y = x[0]
			Src: Offset{
				Base:   LocalValue(0),
				Offset: IntLiteral(0),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.ArrayType{
						Base: ast.TypeLiteral("int"),
						Size: 4,
					},
					false,
				},
			},
			Dst: LocalValue(4),
		},
		MOV{ // x[y] = 6
			Src: IntLiteral(6),
			Dst: Offset{
				Base:   LocalValue(0),
				Offset: LocalValue(4),
				Scale:  IntLiteral(0),
				Container: ast.VarWithType{
					"x", ast.ArrayType{
						Base: ast.TypeLiteral("int"),
						Size: 4,
					},
					false,
				},
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: LocalValue(4),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"x", ast.ArrayType{
							Base: ast.TypeLiteral("int"),
							Size: 4,
						},
						false,
					},
				},
			},
		},
		ADD{
			Left:  LocalValue(4),
			Right: IntLiteral(1),
			Dst:   TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(0),
					Offset: TempValue(0),
					Scale:  IntLiteral(0),
					Container: ast.VarWithType{
						"x", ast.ArrayType{
							Base: ast.TypeLiteral("int"),
							Size: 4,
						},
						false,
					},
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
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{ // mutable x = { 1, 3, 4, 5 }
			Src: IntLiteral(4),
			Dst: LocalValue(0),
		},
		MOV{ // mutable x = { 1, 3, 4, 5 }
			Src: IntLiteral(1),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: LocalValue(3),
		},
		MOV{
			Src: IntLiteral(5),
			Dst: LocalValue(4),
		},
		MOV{ // let y = x[0]
			Src: Offset{
				Base:   LocalValue(1),
				Offset: IntLiteral(0),
				Scale:  IntLiteral(1),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						Base: ast.TypeLiteral("byte"),
					},
					false,
				},
			},
			Dst: LocalValue(5),
		},
		MOV{ // x[y] = 6
			Src: IntLiteral(6),
			Dst: Offset{
				Base:   LocalValue(1),
				Offset: LocalValue(5),
				Scale:  IntLiteral(1),
				Container: ast.VarWithType{
					"x", ast.SliceType{
						Base: ast.TypeLiteral("byte"),
					},
					false,
				},
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(1),
					Offset: LocalValue(5),
					Scale:  IntLiteral(1),
					Container: ast.VarWithType{
						"x", ast.SliceType{
							Base: ast.TypeLiteral("byte"),
						},
						false,
					},
				},
			},
		},
		ADD{
			Left:  LocalValue(5),
			Right: IntLiteral(1),
			Dst:   TempValue(0),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				Offset{
					Base:   LocalValue(1),
					Offset: TempValue(0),
					Scale:  IntLiteral(1),
					Container: ast.VarWithType{
						"x", ast.SliceType{
							Base: ast.TypeLiteral("byte"),
						},
						false,
					},
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
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(6),
			Dst: LocalValue(0),
		},
		MOV{
			Src: StringLiteral("foobar"),
			Dst: LocalValue(1),
		},
		CALL{
			FName: "PrintAString",
			Args: []Register{
				LocalValue(0),
				LocalValue(1),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}

	expected = []Opcode{
		CALL{
			FName: "PrintString",
			Args: []Register{
				FuncArg{0, false},
				FuncArg{1, false},
			},
		},
	}

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCastBuiltin(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.CastBuiltin)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(70),
			Dst: LocalValue(1),
		},
		MOV{
			Src: IntLiteral(111),
			Dst: LocalValue(2),
		},
		MOV{
			Src: IntLiteral(111),
			Dst: LocalValue(3),
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				LocalValue(0),
				Pointer{LocalValue(1)},
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
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: StringLiteral("bar"),
			Dst: LocalValue(1),
		},
		CALL{
			FName: "PrintByteSlice",
			Args: []Register{
				LocalValue(0),
				LocalValue(1),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestCastIntVariable(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.CastIntVariable)
	if err != nil {
		t.Fatal(err)
	}
	i, _, rd, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(65),
			Dst: LocalValue(0),
		},
		MOV{
			Src: LocalValue(0),
			Dst: LocalValue(1),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				LocalValue(1),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
	if got := rd[LocalValue(1)].TypeInfo; got != (ast.TypeInfo{1, false}) {
		t.Errorf("Incorrect type info for byte cast: got %v", got)
	}
	if got := rd[LocalValue(0)].TypeInfo; got != (ast.TypeInfo{0, true}) {
		t.Errorf("Incorrect type info for uncast integer: got %v", got)
	}
}

func TestAssertionFail(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.AssertionFail)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		ASSERT{
			Predicate: Condition{nil, IntLiteral(0)},
			Message:   "",
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}

func TestAssertionPass(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.AssertionPass)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		ASSERT{
			Predicate: Condition{nil, IntLiteral(1)},
			Message:   "",
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}

func TestAssertionFailWithMessage(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.AssertionFailWithMessage)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		ASSERT{
			Predicate: Condition{nil, IntLiteral(0)},
			Message:   "This always fails",
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}

func TestAssertionPassWithMessage(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.AssertionPassWithMessage)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		ASSERT{
			Predicate: Condition{nil, IntLiteral(1)},
			Message:   "You should never see this",
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}

func TestSumTypeFuncCall(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SumTypeFuncCall)
	if err != nil {
		t.Fatal(err)
	}
	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{0, false},
							Right: IntLiteral(0),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							FuncArg{1, false},
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  FuncArg{0, false},
							Right: IntLiteral(1),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintString",
						Args: []Register{
							FuncArg{1, false},
							FuncArg{2, false},
						},
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected = []Opcode{
		CALL{
			FName: "foo",
			Args: []Register{
				IntLiteral(1),
				IntLiteral(3),
				StringLiteral("bar"),
			},
		},
		CALL{
			FName: "foo",
			Args: []Register{
				IntLiteral(0),
				IntLiteral(3),
			},
		},
	}
	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}

func TestSumTypeFuncReturn(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.SumTypeFuncReturn)
	if err != nil {
		t.Fatal(err)
	}

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		IF{
			ControlFlow: ControlFlow{
				Condition: Condition{
					Body:     []Opcode{},
					Register: FuncArg{0, false},
				},
				Body: []Opcode{
					MOV{
						Src: IntLiteral(0),
						Dst: FuncRetVal(0),
					},
					MOV{
						Src: IntLiteral(3),
						Dst: FuncRetVal(1),
					},
					RET{},
				},
			},
			ElseBody: nil,
		},
		MOV{
			Src: IntLiteral(1),
			Dst: FuncRetVal(0),
		},
		MOV{
			Src: IntLiteral(4),
			Dst: FuncRetVal(1),
		},
		MOV{
			Src: StringLiteral("not3"),
			Dst: FuncRetVal(2),
		},
		RET{},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}

	i, _, _, err = Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected = []Opcode{
		CALL{
			FName: "foo",
			Args:  []Register{IntLiteral(0)},
		},
		MOV{
			Src: LastFuncCallRetVal{0, 0},
			Dst: LocalValue(0),
		},
		MOV{
			Src: LastFuncCallRetVal{0, 1},
			Dst: LocalValue(1),
		},
		MOV{
			Src: LastFuncCallRetVal{0, 2},
			Dst: LocalValue(2),
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(0),
							Right: IntLiteral(0),
							Dst:   TempValue(0),
						},
					},
					Register: TempValue(0),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							LocalValue(1),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(0),
							Right: IntLiteral(1),
							Dst:   TempValue(1),
						},
					},
					Register: TempValue(1),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintString",
						Args: []Register{
							LocalValue(1),
							LocalValue(2),
						},
					},
				},
			},
		},
		CALL{
			FName: "foo",
			Args:  []Register{IntLiteral(1)},
		},
		MOV{
			Src: LastFuncCallRetVal{3, 0},
			Dst: LocalValue(3),
		},
		MOV{
			Src: LastFuncCallRetVal{3, 1},
			Dst: LocalValue(4),
		},
		MOV{
			Src: LastFuncCallRetVal{3, 2},
			Dst: LocalValue(5),
		},
		JumpTable{
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(3),
							Right: IntLiteral(0),
							Dst:   TempValue(2),
						},
					},
					Register: TempValue(2),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintInt",
						Args: []Register{
							LocalValue(4),
						},
					},
				},
			},
			ControlFlow{
				Condition: Condition{
					Body: []Opcode{
						EQ{
							Left:  LocalValue(3),
							Right: IntLiteral(1),
							Dst:   TempValue(3),
						},
					},
					Register: TempValue(3),
				},
				Body: []Opcode{
					CALL{
						FName: "PrintString",
						Args: []Register{
							LocalValue(4),
							LocalValue(5),
						},
					},
				},
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}

}

func TestProductTypeValue(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.ProductTypeValue)
	if err != nil {
		t.Fatal(err)
	}

	i, _, _, err := Generate(as[0], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(0),
			Dst: LocalValue(1),
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				LocalValue(0),
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
			FName: "PrintInt",
			Args: []Register{
				LocalValue(1),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}

func TestUserProductTypeValue(t *testing.T) {
	as, ti, c, err := ast.Parse(sampleprograms.UserProductTypeValue)
	if err != nil {
		t.Fatal(err)
	}

	i, _, _, err := Generate(as[1], ti, c, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Opcode{
		MOV{
			Src: IntLiteral(3),
			Dst: LocalValue(0),
		},
		MOV{
			Src: IntLiteral(6),
			Dst: LocalValue(1),
		},
		MOV{
			Src: StringLiteral(`hello\n`),
			Dst: LocalValue(2),
		},
		CALL{
			FName: "PrintString",
			Args: []Register{
				LocalValue(1),
				LocalValue(2),
			},
		},
		CALL{
			FName: "PrintInt",
			Args: []Register{
				LocalValue(0),
			},
		},
	}

	if err := compareIR(i.Body, expected); err != nil {
		t.Errorf("%v", err)
	}
}
