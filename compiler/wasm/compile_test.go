package wasm

import (
	"fmt"
	"testing"

	"github.com/driusan/lang/parser/sampleprograms"
)

func compareModule(got, want Module) error {
	if err := compareImports(got.Imports, want.Imports); err != nil {
		return err
	}

	if err := compareGlobals(got.Globals, want.Globals); err != nil {
		return err
	}
	if err := compareData(got.Data, want.Data); err != nil {
		return err
	}
	if got.Memory != want.Memory {
		return fmt.Errorf("Memory does not match. got %v want %v", got.Memory, want.Memory)
	}
	return compareFuncs(got.Funcs, want.Funcs)
}

func compareImports(got, want []Import) error {
	if len(got) != len(want) {
		return fmt.Errorf("Import lengths do not match: got %v want %v", len(got), len(want))
	}
	for i := range got {
		if err := compareFunc(got[i].Func, want[i].Func); err != nil {
			return fmt.Errorf("Import %d does not match: %v", i, err)
		}
	}
	return nil
}
func compareGlobals(got, want []Global) error {
	if len(got) != len(want) {
		return fmt.Errorf("Global lengths do not match: got %v want %v", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			return fmt.Errorf("Global %d does not match: got %v want %v", i, got, want)
		}
	}
	return nil
}

func compareData(got, want []DataSection) error {
	if len(got) != len(want) {
		return fmt.Errorf("Data section lengths do not match: got %v want %v", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			return fmt.Errorf("Data section %d does not match: got %v want %v", i, got[i], want[i])
		}
	}
	return nil
}

func compareFunc(got, want Func) error {
	if got.Name != want.Name {
		return fmt.Errorf("Function name does not match: got %v want %v", got.Name, want.Name)
	}
	if err := compareSignature(got.Signature, want.Signature); err != nil {
		return fmt.Errorf("Function signatures do not match. %v: got %v want %v", err, got.Signature, want.Signature)
	}
	if err := compareBody(got.Body, want.Body); err != nil {
		return fmt.Errorf("Function bodies do not match. %v: got %v want %v", err, got.Body, want.Body)
	}
	return nil
}

func compareSignature(got, want Signature) error {
	if len(got) != len(want) {
		return fmt.Errorf("Signature lengths do not match: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			return fmt.Errorf("Signatures do not match: got %v want %v", got, want)
		}
	}
	return nil
}

func compareBody(got, want []Instruction) error {
	if len(got) != len(want) {
		return fmt.Errorf("Body lengths do not match: got %v want %v", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			return fmt.Errorf("Instruction %d does not match: got %v want %v", i, got[i], want[i])
		}
	}
	return nil
}
func compareFuncs(got, want []Func) error {
	if len(got) != len(want) {
		return fmt.Errorf("Incorrect number of funcs: got %v want %v", len(got), len(want))
	}
	for i := range want {
		if err := compareFunc(got[i], want[i]); err != nil {
			return fmt.Errorf("Function %d: %v", i, err)
		}
	}
	return nil
}

func TestEmptyMain(t *testing.T) {
	module, err := Parse(sampleprograms.EmptyMain)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Func{
		Func{
			Name: "main",
		},
	}

	if err := compareFuncs(module.Funcs, expected); err != nil {
		t.Fatal(err)
	}
}

func TestHelloWorld(t *testing.T) {
	module, err := Parse(sampleprograms.HelloWorld)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintString", Func{
				Name: "PrintString",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\0e\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`Hello, world!\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(0),
					Call{"PrintString"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatalf("%v: %v", module, err)
	}
}

func TestSimpleFunc(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleFunc)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			}},
		},
		Funcs: []Func{
			Func{
				Name: "foo",
				Signature: Signature{
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					I32Const(3),
					Return{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					Call{"foo"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestTypeInference(t *testing.T) {
	module, err := Parse(sampleprograms.TypeInference)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			}},
			{"stdlib", "PrintString", Func{
				Name: "PrintString",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			}},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\02\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`, `, // str content
			},
			DataSection{
				16,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "foo",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Result, ""},
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					GetLocal(0),
					SetLocal(1),
					GetLocal(1),
					I32Const(1),
					I32Add{},
					SetLocal(1),
					GetLocal(1),
					I32Const(1),
					I32Add{},
					SetLocal(2),
					GetLocal(2),
					I32Const(3),
					I32GT_S{},
					If{},
					GetLocal(1),
					Return{},
					End{},
					I32Const(0),
					Return{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(1),
					Call{"foo"},
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					I32Const(3),
					Call{"foo"},
					Call{"PrintInt"},
					I32Const(16),
					Call{"PrintString"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSumToTenRecursive(t *testing.T) {
	module, err := Parse(sampleprograms.SumToTenRecursive)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			}},
		},
		Funcs: []Func{
			Func{
				Name: "sum",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					I32Const(0),
					GetLocal(0),
					Call{"partial_sum"},
					Return{},
				},
			},
			Func{
				Name: "partial_sum",
				Signature: Signature{
					Variable{i32, Param, "partial"},
					Variable{i32, Param, "x"},
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					GetLocal(1),
					I32Const(0),
					I32EQ{},
					If{},
					GetLocal(0),
					Return{},
					End{},
					GetLocal(0),
					GetLocal(1),
					I32Add{},
					GetLocal(1),
					I32Const(1),
					I32Sub{},
					Call{"partial_sum"},
					Return{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(10),
					Call{"sum"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSumToTen(t *testing.T) {
	module, err := Parse(sampleprograms.SumToTen)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			}},
		},
		Funcs: []Func{
			Func{
				Name: "sum",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Result, ""},
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					GetLocal(0),
					SetLocal(1),
					I32Const(0),
					SetLocal(2),
					Block{},
					Loop{},
					GetLocal(1),
					I32Const(0),
					I32GT_S{},
					I32EQZ{},
					BrIf(1), // Break out of the block if the condition isn't met
					GetLocal(2),
					GetLocal(1),
					I32Add{},
					SetLocal(2),
					GetLocal(1),
					I32Const(1),
					I32Sub{},
					SetLocal(1),
					Br(0), // Branch back to the start of the loop
					End{},
					End{},
					GetLocal(2),
					Return{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(10),
					Call{"sum"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestFizzbuzz(t *testing.T) {
	module, err := Parse(sampleprograms.Fizzbuzz)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt", Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\08\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`fizzbuzz`, // str content
			},
			DataSection{
				16,
				`\04\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`buzz`, // str content
			},
			DataSection{
				32,
				`\04\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`fizz`, // str content
			},
			DataSection{
				48,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				56,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(0),
					SetLocal(0),
					I32Const(1),
					SetLocal(1),
					Block{},
					Loop{},
					GetLocal(0),
					I32Const(1),
					I32NE{},
					I32EQZ{},
					BrIf(1),
					GetLocal(1),
					I32Const(15),
					I32Rem_S{},
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					GetLocal(1),
					I32Const(5),
					I32Rem_S{},
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(16),
					Call{"PrintString"},
					Else{},
					GetLocal(1),
					I32Const(3),
					I32Rem_S{},
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(32),
					Call{"PrintString"},
					Else{},
					GetLocal(1),
					Call{"PrintInt"},
					End{},
					End{},
					End{},
					I32Const(48),
					Call{"PrintString"},
					GetLocal(1),
					I32Const(1),
					I32Add{},
					SetLocal(1),
					GetLocal(1),
					I32Const(100),
					I32GE_S{},
					If{},
					I32Const(1),
					SetLocal(0),
					End{},
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestLetStatement(t *testing.T) {
	module, err := Parse(sampleprograms.LetStatement)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt", Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(5),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestLetStatementShadow(t *testing.T) {
	module, err := Parse(sampleprograms.LetStatementShadow)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
			DataSection{
				16,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`hello`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(5),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					I32Const(16),
					SetLocal(1),
					GetLocal(1),
					Call{"PrintString"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestMutAddition(t *testing.T) {
	module, err := Parse(sampleprograms.MutAddition)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					GetLocal(0),
					I32Const(1),
					I32Add{},
					SetLocal(1),
					GetLocal(1),
					I32Const(1),
					I32Add{},
					GetLocal(0),
					I32Add{},
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSomeMath(t *testing.T) {
	module, err := Parse(sampleprograms.SomeMath)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`Add: `, // str content
			},
			DataSection{
				16,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`\n`, // str content
			},
			DataSection{
				32,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`Sub: `, // str content
			},
			DataSection{
				48,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				56,
				`Mul: `, // str content
			},
			DataSection{
				64,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				72,
				`Div: `, // str content
			},
			DataSection{
				80,
				`\09\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				88,
				`Complex: `, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
				},
				Body: []Instruction{
					I32Const(1),
					I32Const(2),
					I32Add{},
					SetLocal(0),
					I32Const(1),
					I32Const(2),
					I32Sub{},
					SetLocal(1),
					I32Const(2),
					I32Const(3),
					I32Mul{},
					SetLocal(2),
					I32Const(6),
					I32Const(2),
					I32Div_S{},
					SetLocal(3),
					I32Const(2),
					I32Const(3),
					I32Mul{},
					I32Const(4),
					I32Const(2),
					I32Div_S{},
					I32Sub{},
					I32Const(1),
					I32Add{},
					SetLocal(4),
					I32Const(0),
					Call{"PrintString"},
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(16),
					Call{"PrintString"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(1),
					Call{"PrintInt"},
					I32Const(16),
					Call{"PrintString"},
					I32Const(48),
					Call{"PrintString"},
					GetLocal(2),
					Call{"PrintInt"},
					I32Const(16),
					Call{"PrintString"},
					I32Const(64),
					Call{"PrintString"},
					GetLocal(3),
					Call{"PrintInt"},
					I32Const(16),
					Call{"PrintString"},
					I32Const(80),
					Call{"PrintString"},
					GetLocal(4),
					Call{"PrintInt"},
					I32Const(16),
					Call{"PrintString"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestEqualComparison(t *testing.T) {
	module, err := Parse(sampleprograms.EqualComparison)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`true\n`, // str content
			},
			DataSection{
				16,
				`\06\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`false\n`, // str content
			},
			DataSection{
				32,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					GetLocal(0),
					GetLocal(1),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					I32Const(16),
					Call{"PrintString"},
					End{},
					Block{},
					Loop{},
					GetLocal(0),
					GetLocal(1),
					I32EQ{},
					I32EQZ{},
					BrIf(1),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(0),
					I32Const(1),
					I32Add{},
					SetLocal(0),
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestNotEqualComparison(t *testing.T) {
	module, err := Parse(sampleprograms.NotEqualComparison)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`true\n`, // str content
			},
			DataSection{
				16,
				`\06\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`false\n`, // str content
			},
			DataSection{
				32,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					GetLocal(0),
					GetLocal(1),
					I32NE{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					I32Const(16),
					Call{"PrintString"},
					End{},
					Block{},
					Loop{},
					GetLocal(0),
					GetLocal(1),
					I32NE{},
					I32EQZ{},
					BrIf(1),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(0),
					I32Const(1),
					I32Add{},
					SetLocal(0),
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestGreaterComparison(t *testing.T) {
	module, err := Parse(sampleprograms.GreaterComparison)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`true\n`, // str content
			},
			DataSection{
				16,
				`\06\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`false\n`, // str content
			},
			DataSection{
				32,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					GetLocal(0),
					GetLocal(1),
					I32GT_S{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					I32Const(16),
					Call{"PrintString"},
					End{},
					Block{},
					Loop{},
					GetLocal(0),
					GetLocal(1),
					I32GT_S{},
					I32EQZ{},
					BrIf(1),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(0),
					I32Const(1),
					I32Sub{},
					SetLocal(0),
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestGreaterOrEqualComparison(t *testing.T) {
	module, err := Parse(sampleprograms.GreaterOrEqualComparison)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`true\n`, // str content
			},
			DataSection{
				16,
				`\06\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`false\n`, // str content
			},
			DataSection{
				32,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					GetLocal(0),
					GetLocal(1),
					I32GE_S{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					I32Const(16),
					Call{"PrintString"},
					End{},
					Block{},
					Loop{},
					GetLocal(0),
					GetLocal(1),
					I32GE_S{},
					I32EQZ{},
					BrIf(1),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(0),
					I32Const(1),
					I32Sub{},
					SetLocal(0),
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestLessThanComparison(t *testing.T) {
	module, err := Parse(sampleprograms.LessThanComparison)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`true\n`, // str content
			},
			DataSection{
				16,
				`\06\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`false\n`, // str content
			},
			DataSection{
				32,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					GetLocal(0),
					GetLocal(1),
					I32LT_S{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					I32Const(16),
					Call{"PrintString"},
					End{},
					Block{},
					Loop{},
					GetLocal(0),
					GetLocal(1),
					I32LT_S{},
					I32EQZ{},
					BrIf(1),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(0),
					I32Const(1),
					I32Add{},
					SetLocal(0),
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestLessThanOrEqualComparison(t *testing.T) {
	module, err := Parse(sampleprograms.LessThanOrEqualComparison)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`true\n`, // str content
			},
			DataSection{
				16,
				`\06\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`false\n`, // str content
			},
			DataSection{
				32,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					I32Const(1),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					GetLocal(0),
					GetLocal(1),
					I32LE_S{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					I32Const(16),
					Call{"PrintString"},
					End{},
					Block{},
					Loop{},
					GetLocal(0),
					GetLocal(1),
					I32LE_S{},
					I32EQZ{},
					BrIf(1),
					GetLocal(0),
					Call{"PrintInt"},
					I32Const(32),
					Call{"PrintString"},
					GetLocal(0),
					I32Const(1),
					I32Add{},
					SetLocal(0),
					Br(0),
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeUint8(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeUint8)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeInt8(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeInt8)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(-4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeUint16(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeUint16)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeInt16(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeInt16)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(-4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeUint32(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeUint32)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeInt32(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeInt32)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(-4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeUint64(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeUint64)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i64, Local, "LV0"},
				},
				Body: []Instruction{
					I64Const(4),
					SetLocal(0),
					GetLocal(0),
					I32WrapI64{},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestConcreteTypeInt64(t *testing.T) {
	module, err := Parse(sampleprograms.ConcreteTypeInt64)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i64, Local, "LV0"},
				},
				Body: []Instruction{
					I64Const(-4),
					SetLocal(0),
					GetLocal(0),
					I32WrapI64{},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestFibonacci(t *testing.T) {
	module, err := Parse(sampleprograms.Fibonacci)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "fib_rec",
				Signature: Signature{
					Variable{i64, Param, "n"},
					Variable{i64, Param, "n1"},
					Variable{i64, Result, ""},
					Variable{i64, Local, "LV0"},
				},
				Body: []Instruction{
					GetLocal(0),
					GetLocal(1),
					I64Add{},
					SetLocal(2),
					GetLocal(2),
					I64Const(200),
					I64GE_U{},
					If{},
					GetLocal(1),
					Return{},
					End{},
					GetLocal(2),
					I32WrapI64{},
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetLocal(1),
					GetLocal(2),
					Call{"fib_rec"},
					Return{},
				},
			},
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i64, Local, "LV0"},
				},
				Body: []Instruction{
					I64Const(1),
					I64Const(1),
					Call{"fib_rec"},
					Drop{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleAlgorithm(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleAlgorithm)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "loop",
				Signature: Signature{
					Variable{i32, Param, "high"},
					Variable{i32, Result, ""},
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
				},
				Body: []Instruction{
					I32Const(0),
					SetLocal(1),
					I32Const(0),
					SetLocal(2),
					GetLocal(0),
					I32Const(2),
					I32Mul{},
					SetLocal(3),
					I32Const(1),
					SetLocal(2),
					Block{},
					Loop{},
					GetLocal(2),
					GetLocal(3),
					I32LT_S{},
					I32EQZ{},
					BrIf(1),
					GetLocal(2),
					I32Const(2),
					I32Rem_S{},
					I32Const(0),
					I32EQ{},
					If{},
					GetLocal(2),
					I32Const(2),
					I32Mul{},
					GetLocal(1),
					I32Add{},
					SetLocal(1),
					End{},
					GetLocal(2),
					I32Const(1),
					I32Add{},
					SetLocal(2),
					Br(0),

					End{},
					End{},

					GetLocal(1),
					Return{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(10),
					Call{"loop"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestTwoProcs(t *testing.T) {
	module, err := Parse(sampleprograms.TwoProcs)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "foo",
				Signature: Signature{
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					I32Const(3),
					Return{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					Call{"foo"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestOutOfOrder(t *testing.T) {
	module, err := Parse(sampleprograms.OutOfOrder)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Body: []Instruction{
					Call{"foo"},
					Call{"PrintInt"},
				},
			},
			Func{
				Name: "foo",
				Signature: Signature{
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					I32Const(3),
					Return{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleMatch(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleMatch)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\07\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`I am 1\n`, // str content
			},
			DataSection{
				16,
				`\07\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`I am 2\n`, // str content
			},
			DataSection{
				32,
				`\07\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`I am 4\n`, // str content
			},
			DataSection{
				48,
				`\07\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				56,
				`I am 3\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					/* FIXME: Should this use br_table instead of converting to an if/else chain? */
					// let x = 3
					I32Const(3),
					SetLocal(0),

					// match case 1
					GetLocal(0),
					I32Const(1),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},

					// match case 2
					Else{},
					GetLocal(0),
					I32Const(2),
					I32EQ{},
					If{},
					I32Const(16),
					Call{"PrintString"},

					// match case 4 (third case)
					Else{},
					GetLocal(0),
					I32Const(4),
					I32EQ{},
					If{},
					I32Const(32),
					Call{"PrintString"},

					// match case 3 (last case)
					Else{},
					GetLocal(0),
					I32Const(3),
					I32EQ{},
					If{},
					I32Const(48),
					Call{"PrintString"},

					End{},
					End{},
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestIfElseMatch(t *testing.T) {
	module, err := Parse(sampleprograms.IfElseMatch)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\11\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`x is less than 3\n`, // str conten2
			},
			DataSection{
				32,
				`\14\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`x is greater than 3\n`, // str content
			},
			DataSection{
				64,
				`\11\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				72,
				`x is less than 4\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					// let x = 3
					I32Const(3),
					SetLocal(0),

					// match case 1
					GetLocal(0),
					I32Const(3),
					I32LT_S{},
					If{},
					I32Const(0),
					Call{"PrintString"},

					// match case 2
					Else{},
					GetLocal(0),
					I32Const(3),
					I32GT_S{},
					If{},
					I32Const(32),
					Call{"PrintString"},

					// match case 4 (third case)
					Else{},
					GetLocal(0),
					I32Const(4),
					I32LT_S{},
					If{},
					I32Const(64),
					Call{"PrintString"},

					End{},
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestUserDefinedType(t *testing.T) {
	module, err := Parse(sampleprograms.UserDefinedType)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestEnumType(t *testing.T) {
	module, err := Parse(sampleprograms.EnumType)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\08\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`I am A!\n`, // str conten2
			},
			DataSection{
				24,
				`\08\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				32,
				`I am B!\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(0),
					SetLocal(0),
					GetLocal(0),
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					GetLocal(0),
					I32Const(1),
					I32EQ{},
					If{},
					I32Const(24),
					Call{"PrintString"},
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestEnumTypeInferred(t *testing.T) {
	module, err := Parse(sampleprograms.EnumTypeInferred)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\08\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`I am A!\n`, // str conten2
			},
			DataSection{
				24,
				`\08\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				32,
				`I am B!\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(1),
					SetLocal(0),
					GetLocal(0),
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					GetLocal(0),
					I32Const(1),
					I32EQ{},
					If{},
					I32Const(24),
					Call{"PrintString"},
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestMatchParam(t *testing.T) {
	module, err := Parse(sampleprograms.MatchParam)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "foo",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Param, ""},
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					// match case 1
					GetLocal(0),
					I32Const(1),
					I32EQ{},
					If{},
					GetLocal(1),
					Return{},
					// match case 2
					Else{},
					GetLocal(0),
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Return{},
					End{},
					End{},
					Unreachable{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(1),
					I32Const(5),
					Call{"foo"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestMatchParam2(t *testing.T) {
	module, err := Parse(sampleprograms.MatchParam2)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`x`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "foo",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Param, ""},
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					I32Const(0),
					Call{"PrintString"},
					// match case 1
					GetLocal(0),
					I32Const(1),
					I32EQ{},
					If{},
					GetLocal(1),
					Return{},
					// match case 2
					Else{},
					GetLocal(0),
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Return{},
					End{},
					End{},
					Unreachable{},
				},
			},
			Func{
				Name: "main",
				Body: []Instruction{
					I32Const(1),
					I32Const(5),
					Call{"foo"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestGenericEnumType(t *testing.T) {
	module, err := Parse(sampleprograms.GenericEnumType)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\0e\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`I am nothing!\n`, // str content
			},
			DataSection{
				24,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				32,
				`\n`, // str content
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 40,
			},
		},
		Funcs: []Func{
			Func{
				Name: "DoSomething",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					GetLocal(0),
					I32Const(3),
					I32GT_S{},
					If{},
					GetGlobal(0),
					I32Const(0),
					I32Store{},
					GetGlobal(0),
					Return{},
					End{},
					GetGlobal(0),
					I32Const(1),
					I32Store{},
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Const(5),
					I32Store{},
					GetGlobal(0),
					Return{},
				},
			},
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
				},
				Body: []Instruction{
					I32Const(3),
					Call{"DoSomething"},
					Drop{},
					GetGlobal(0),
					I32Load{},
					SetLocal(0),
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Load{},
					SetLocal(1),
					GetLocal(0),
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					GetLocal(0),
					I32Const(1),
					I32EQ{},
					If{},
					GetLocal(1),
					Call{"PrintInt"},
					I32Const(24),
					Call{"PrintString"},
					End{},
					End{},
					I32Const(4),
					Call{"DoSomething"},
					Drop{},
					GetGlobal(0),
					I32Load{},
					SetLocal(2),
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Load{},
					SetLocal(3),
					GetLocal(2),
					I32Const(0),
					I32EQ{},
					If{},
					I32Const(0),
					Call{"PrintString"},
					Else{},
					GetLocal(2),
					I32Const(1),
					I32EQ{},
					If{},
					GetLocal(3),
					Call{"PrintInt"},
					I32Const(24),
					Call{"PrintString"},
					End{},
					End{},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleArray(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleArray)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
				},
				Body: []Instruction{
					I32Const(1),
					SetLocal(0),
					I32Const(2),
					SetLocal(1),
					I32Const(3),
					SetLocal(2),
					I32Const(4),
					SetLocal(3),
					I32Const(5),
					SetLocal(4),
					GetLocal(3),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleArrayInference(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleArrayInference)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
				},
				Body: []Instruction{
					I32Const(1),
					SetLocal(0),
					I32Const(2),
					SetLocal(1),
					I32Const(3),
					SetLocal(2),
					I32Const(4),
					SetLocal(3),
					I32Const(5),
					SetLocal(4),
					GetLocal(3),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestArrayMutation(t *testing.T) {
	module, err := Parse(sampleprograms.ArrayMutation)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
				},
				Body: []Instruction{
					I32Const(1),
					SetLocal(0),
					I32Const(2),
					SetLocal(1),
					I32Const(3),
					SetLocal(2),
					I32Const(4),
					SetLocal(3),
					I32Const(5),
					SetLocal(4),
					GetLocal(3),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					I32Const(2),
					SetLocal(3),
					GetLocal(3),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetLocal(2),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleSlice(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleSlice)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
					Variable{i32, Local, "LV5"},
				},
				Body: []Instruction{
					I32Const(5),
					SetLocal(0),
					I32Const(1),
					SetLocal(1),
					I32Const(2),
					SetLocal(2),
					I32Const(3),
					SetLocal(3),
					I32Const(4),
					SetLocal(4),
					I32Const(5),
					SetLocal(5),
					GetLocal(4),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleSliceInference(t *testing.T) {
	module, err := Parse(sampleprograms.SimpleSliceInference)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
					Variable{i32, Local, "LV5"},
				},
				Body: []Instruction{
					I32Const(5),
					SetLocal(0),
					I32Const(1),
					SetLocal(1),
					I32Const(2),
					SetLocal(2),
					I32Const(3),
					SetLocal(3),
					I32Const(4),
					SetLocal(4),
					I32Const(5),
					SetLocal(5),
					GetLocal(4),
					Call{"PrintInt"},
				},
			},
		},
	}
	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSliceMutation(t *testing.T) {
	module, err := Parse(sampleprograms.SliceMutation)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},
					Variable{i32, Local, "LV4"},
					Variable{i32, Local, "LV5"},
				},
				Body: []Instruction{
					I32Const(5),
					SetLocal(0),
					I32Const(1),
					SetLocal(1),
					I32Const(2),
					SetLocal(2),
					I32Const(3),
					SetLocal(3),
					I32Const(4),
					SetLocal(4),
					I32Const(5),
					SetLocal(5),
					GetLocal(4),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					I32Const(2),
					SetLocal(4),
					GetLocal(4),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetLocal(3),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestIndexAssignment(t *testing.T) {
	module, err := Parse(sampleprograms.IndexAssignment)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"}, // x
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},

					Variable{i32, Local, "LV4"}, // n
					Variable{i32, Local, "LV5"}, // n2
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					I32Const(4),
					SetLocal(2),
					I32Const(5),
					SetLocal(3),

					GetLocal(2),
					SetLocal(4),
					GetLocal(3),
					SetLocal(5),
					GetLocal(4),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetLocal(5),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestIndexedAddition(t *testing.T) {
	module, err := Parse(sampleprograms.IndexedAddition)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"}, // x
					Variable{i32, Local, "LV1"},
					Variable{i32, Local, "LV2"},
					Variable{i32, Local, "LV3"},

					Variable{i32, Local, "LV4"}, // n
					Variable{i32, Local, "LV5"}, // n2
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					I32Const(3),
					SetLocal(1),
					I32Const(4),
					SetLocal(2),
					I32Const(5),
					SetLocal(3),

					GetLocal(2),
					SetLocal(4),

					GetLocal(4),
					GetLocal(3),
					I32Add{},
					SetLocal(4),

					GetLocal(3),
					GetLocal(1),
					I32Add{},
					SetLocal(5),
					GetLocal(4),
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetLocal(5),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestArrayIndex(t *testing.T) {
	module, err := Parse(sampleprograms.ArrayIndex)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			Import{
				"stdlib",
				"PrintInt",
				Func{
					Name: "PrintInt",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
			Import{
				"stdlib",
				"PrintString",
				Func{
					Name: "PrintString",
					Signature: Signature{
						Variable{i32, Param, ""},
					},
				},
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 16,
			},
		},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					// let x = 3
					I32Const(3), // 0
					SetLocal(0), // 1

					// let n = { 1, 2, 3, 4, 5}
					GetGlobal(0), // 2
					I32Const(1),  // 3
					I32Store{},   // 4

					GetGlobal(0), // 5
					I32Const(4),
					I32Add{},
					I32Const(2),
					I32Store{},

					GetGlobal(0), // 10
					I32Const(8),
					I32Add{},
					I32Const(3),
					I32Store{},

					GetGlobal(0), // 15
					I32Const(12),
					I32Add{},
					I32Const(4),
					I32Store{},

					GetGlobal(0), // 20
					I32Const(16),
					I32Add{},
					I32Const(5),
					I32Store{},

					// mut n = { 1, 2, 3, 4, 5}
					GetGlobal(0), // 25
					I32Const(20),
					I32Add{},
					I32Const(1),
					I32Store{},

					GetGlobal(0), // 30
					I32Const(24),
					I32Add{},
					I32Const(2),
					I32Store{},

					GetGlobal(0), // 35
					I32Const(28),
					I32Add{},
					I32Const(3),
					I32Store{},

					GetGlobal(0), // 40
					I32Const(32),
					I32Add{},
					I32Const(4),
					I32Store{},

					GetGlobal(0), // 45
					I32Const(36),
					I32Add{},
					I32Const(5),
					I32Store{},

					// PrintInt(n[x])
					GetLocal(0), // 51
					I32Const(4), // 52
					I32Mul{},
					GetGlobal(0), // 50
					I32Add{},
					I32Load{}, // 55
					Call{"PrintInt"},

					// PrintString
					I32Const(0),
					Call{"PrintString"},

					// PrintInt(n2[x+1])
					GetLocal(0), // x+1
					I32Const(1),
					I32Add{},
					I32Const(4), // offset scale
					I32Mul{},
					GetGlobal(0), // 59
					I32Const(20), // n2 base
					I32Add{},
					I32Add{},  // sp = n2base+scaled offset
					I32Load{}, // load the value
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSliceParam(t *testing.T) {
	module, err := Parse(sampleprograms.SliceParam)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintByteSlice", Func{
				Name: "PrintByteSlice",
				Signature: Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 0,
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					// let x = 3
					I32Const(3),
					SetLocal(0),

					// let b = { 44, 55, 88}
					GetGlobal(0),
					I32Const(44),
					I32Store8{},

					GetGlobal(0),
					I32Const(1),
					I32Add{},
					I32Const(55),
					I32Store8{},

					GetGlobal(0),
					I32Const(2),
					I32Add{},
					I32Const(88),
					I32Store8{},

					// PrintASlice(b)
					GetLocal(0),
					GetGlobal(0),
					Call{"PrintASlice"},
				},
			},
			Func{
				Name: "PrintASlice",
				Signature: Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
				},
				Body: []Instruction{
					GetLocal(0),
					GetLocal(1),
					Call{"PrintByteSlice"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSliceStringParam(t *testing.T) {
	module, err := Parse(sampleprograms.SliceStringParam)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintString", Func{
				Name: "PrintString",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Data: []DataSection{
			DataSection{
				0,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`foo`, // str content
			},
			DataSection{
				16,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`bar`, // str content
			},
			DataSection{
				32,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`baz`, // str content
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 48,
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Funcs: []Func{
			Func{
				Name: "PrintSecond",
				Signature: Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
				},
				Body: []Instruction{
					GetLocal(1),
					I32Const(4),
					I32Add{},
					I32Load{},
					Call{"PrintString"},
				},
			},
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					GetGlobal(0),
					I32Const(0),
					I32Store{},
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Const(16),
					I32Store{},
					GetGlobal(0),
					I32Const(8),
					I32Add{},
					I32Const(32),
					I32Store{},
					GetLocal(0),
					GetGlobal(0),
					Call{"PrintSecond"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSliceStringVariableParam(t *testing.T) {
	module, err := Parse(sampleprograms.SliceStringVariableParam)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintString", Func{
				Name: "PrintString",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Data: []DataSection{
			DataSection{
				0,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`foo`, // str content
			},
			DataSection{
				16,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`bar`, // str content
			},
			DataSection{
				32,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`baz`, // str content
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 48,
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Funcs: []Func{
			Func{
				Name: "PrintSecond",
				Signature: Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(1),
					SetLocal(2),
					GetLocal(2),
					I32Const(4),
					I32Mul{},
					GetLocal(1),
					I32Add{},
					I32Load{},
					Call{"PrintString"},
				},
			},
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(3),
					SetLocal(0),
					GetGlobal(0),
					I32Const(0),
					I32Store{},
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Const(16),
					I32Store{},
					GetGlobal(0),
					I32Const(8),
					I32Add{},
					I32Const(32),
					I32Store{},
					GetLocal(0),
					GetGlobal(0),
					Call{"PrintSecond"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestReferenceVariable(t *testing.T) {
	module, err := Parse(sampleprograms.ReferenceVariable)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
			{"stdlib", "PrintString", Func{
				Name: "PrintString",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`\n`, // str content
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 16,
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Funcs: []Func{
			Func{
				Name: "changer",
				Signature: Signature{
					Variable{i32, Param, "x"},
					Variable{i32, Param, "y"},
					Variable{i32, Result, ""},
				},
				Body: []Instruction{
					GetLocal(0),
					I32Const(4),
					I32Store{},
					GetLocal(0),
					I32Load{},
					GetLocal(1),
					I32Add{},
					Return{},
				},
			},
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV1"},
				},
				Body: []Instruction{
					GetGlobal(0),
					I32Const(3),
					I32Store{},
					GetGlobal(0),
					I32Load{},
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetGlobal(0),
					I32Const(3),
					Call{"changer"},
					SetLocal(0),
					GetGlobal(0),
					I32Load{},
					Call{"PrintInt"},
					I32Const(0),
					Call{"PrintString"},
					GetLocal(0),
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSliceLength(t *testing.T) {
	module, err := Parse(sampleprograms.SliceLength)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "len", Func{
				Name: "len",
				Signature: Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
					Variable{i32, Result, ""},
				},
			},
			},
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Data: []DataSection{
			DataSection{
				0,
				`\01\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				8,
				`3`, // str content
			},
			DataSection{
				16,
				`\03\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				24,
				`foo`, // str content
			},
			DataSection{
				32,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				40,
				`hello`, // str content
			},
			DataSection{
				48,
				`\05\00\00\00\00\00\00\00`, // str size
			},
			DataSection{
				56,
				`world`, // str content
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 64,
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(4),
					SetLocal(0),
					GetGlobal(0),
					I32Const(0),
					I32Store{},
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Const(16),
					I32Store{},
					GetGlobal(0),
					I32Const(8),
					I32Add{},
					I32Const(32),
					I32Store{},
					GetGlobal(0),
					I32Const(12),
					I32Add{},
					I32Const(48),
					I32Store{},
					GetLocal(0),
					GetGlobal(0),
					Call{"len"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}

func TestSliceLength2(t *testing.T) {
	module, err := Parse(sampleprograms.SliceLength2)
	if err != nil {
		t.Fatal(err)
	}

	expected := Module{
		Imports: []Import{
			{"stdlib", "len", Func{
				Name: "len",
				Signature: Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
					Variable{i32, Result, ""},
				},
			},
			},
			{"stdlib", "PrintInt", Func{
				Name: "PrintInt",
				Signature: Signature{
					Variable{i32, Param, ""},
				},
			},
			},
		},
		Globals: []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: 0,
			},
		},
		Memory: Memory{Name: "mem", Size: 1},
		Funcs: []Func{
			Func{
				Name: "main",
				Signature: Signature{
					Variable{i32, Local, "LV0"},
				},
				Body: []Instruction{
					I32Const(5),
					SetLocal(0),
					GetGlobal(0),
					I32Const(3),
					I32Store{},
					GetGlobal(0),
					I32Const(4),
					I32Add{},
					I32Const(3),
					I32Store{},
					GetGlobal(0),
					I32Const(8),
					I32Add{},
					I32Const(3),
					I32Store{},
					GetGlobal(0),
					I32Const(12),
					I32Add{},
					I32Const(3),
					I32Store{},
					GetGlobal(0),
					I32Const(16),
					I32Add{},
					I32Const(3),
					I32Store{},
					GetLocal(0),
					GetGlobal(0),
					Call{"len"},
					Call{"PrintInt"},
				},
			},
		},
	}

	if err := compareModule(module, expected); err != nil {
		t.Fatal(err)
	}
}
