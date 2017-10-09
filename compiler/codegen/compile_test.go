package codegen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/driusan/lang/compiler/ir"
	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/sampleprograms"
)

func RunProgram(name, p string) error {
	dir, err := ioutil.TempDir("", "langtest"+name)
	if err != nil {
		return err
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	exe, err := BuildProgram(dir, strings.NewReader(p))
	if err != nil {
		return err
	}
	cmd := exec.Command(dir + "/" + exe)
	val, err := cmd.Output()
	fmt.Println(string(val))
	return err
}

func TestCompileHelloWorld(t *testing.T) {
	if debug {
		t.Fatal("Can not run helloworld test in debug mode.")
	}
	prgAst, types, callables, err := ast.Parse(sampleprograms.HelloWorld)
	if err != nil {
		t.Fatal(err)
	}

	var w bytes.Buffer
	fnc, _, err := ir.Generate(prgAst[0], types, callables, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := Compile(&w, fnc); err != nil {
		t.Fatal(err)
	}

	// HelloWorld is simple enough that we can hand compile it. Most other
	// programs we just compile it, run it and check the output to ensure
	// that they work.
	expected := `TEXT main(SB), 4+16, $40
	DATA string0<>+0(SB)/8, $14
	DATA string0<>+8(SB)/8, $"Hello, w"
	DATA string0<>+16(SB)/8, $"orld!\n\000\000"
	GLOBL string0<>+0(SB), 8+16, $24
		MOVQ $string0<>+0(SB), BX
	MOVQ BX, 0(SP)
	CALL PrintString+0(SB)
	RET
`

	if val := w.Bytes(); string(val) != expected {
		t.Errorf("Unexpected ASM output compiling HelloWorld. got %s; want %s", val, expected)
	}
}

// Test that Create/Write work correctly. This isn't done as an Example test because we
// need to know a little outside context (ie. what the current directory is, we also
// need to be able to read the file from the test, to make sure it got written correctly
// and cleaned up correctly.
func TestCreateSyscall(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)
	// Inline the RunProgram call, because we want to cd to the directory so that the foo.txt
	// file isn't created as garbage in the cwd..
	dir, err := ioutil.TempDir("", "langtestCreateSyscall")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	exe, err := BuildProgram(dir, strings.NewReader(sampleprograms.CreateSyscall))
	if err != nil {
		t.Fatal(err)
	}

	// Make sure foo.txt gets created in dir, so that the defer cleans it up..
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("./" + exe)
	if _, err := cmd.Output(); err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile("foo.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	if string(content) != "Hello\n" {
		t.Errorf("Unexpected content of file foo.txt: got %v want %v", string(content), "Hello\n")
	}
}

// Test that Open/Readwork correctly. This isn't done as an Example test because we
// need to know a little outside context (ie. what the current directory is, we also
// need to be able to read the file from the test, to make sure it got written correctly
// and cleaned up correctly.
func TestReadSyscall(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)
	// Inline the RunProgram call, because we want to cd to the directory so that the foo.txt
	// file isn't created as garbage in the cwd..
	dir, err := ioutil.TempDir("", "langtestReadSyscall")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// There's currently no way to define constants, no bitwise operations, and only decimal
	// numbers, so modes and flags are hardcoded in decimal.
	exe, err := BuildProgram(dir, strings.NewReader(sampleprograms.ReadSyscall))
	if err != nil {
		t.Fatal(err)
	}

	// Make sure foo.txt gets created in dir, so that the defer cleans it up..
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("foo.txt", []byte("Hello, world!"), 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("./" + exe)

	content, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Hello," {
		t.Errorf("Unexpected content of file foo.txt: got %v want %v", string(content), "Hello,")
	}

	// Run it again with different file content.
	if err := ioutil.WriteFile("foo.txt", []byte("Goodbye"), 0755); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("./" + exe)

	content, err = cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Goodby" {
		t.Errorf("Unexpected content of file foo.txt: got %v want %v", string(content), "Goodby")
	}
}

func ExampleCompileEmptyMain() {
	if err := RunProgram("emptymain", sampleprograms.EmptyMain); err != nil {
		fmt.Println(err.Error())
	}
	// Output:
}

func ExampleHelloWorld() {
	if err := RunProgram("helloworld", sampleprograms.HelloWorld); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Hello, world!
}

func ExampleLetStatement() {
	if err := RunProgram("letstatement", sampleprograms.LetStatement); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 5
}

func ExampleLetStatementShadow() {
	if err := RunProgram("letstatementshadow", sampleprograms.LetStatementShadow); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 5
	// hello
}

func ExampleCompileTwoProcs() {
	if err := RunProgram("twoprocs", sampleprograms.TwoProcs); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 3
}

func ExampleOutOfOrder() {
	if err := RunProgram("outoforder", sampleprograms.OutOfOrder); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 3
}

func ExampleMutAddition() {
	if err := RunProgram("mutaddition", sampleprograms.MutAddition); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 8
}
func ExampleSimpleFunc() {
	if err := RunProgram("simplefunc", sampleprograms.SimpleFunc); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 3
}

func ExampleSumToTen() {
	if err := RunProgram("sumtoten", sampleprograms.SumToTen); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 55
}

func ExampleSumToTenRecursive() {
	if err := RunProgram("sumtotenrecursive", sampleprograms.SumToTenRecursive); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 55
}

func ExampleFizzBuzz() {
	if err := RunProgram("fizzbuzz", sampleprograms.Fizzbuzz); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 1
	// 2
	// fizz
	// 4
	// buzz
	// fizz
	// 7
	// 8
	// fizz
	// buzz
	// 11
	// fizz
	// 13
	// 14
	// fizzbuzz
	// 16
	// 17
	// fizz
	// 19
	// buzz
	// fizz
	// 22
	// 23
	// fizz
	// buzz
	// 26
	// fizz
	// 28
	// 29
	// fizzbuzz
	// 31
	// 32
	// fizz
	// 34
	// buzz
	// fizz
	// 37
	// 38
	// fizz
	// buzz
	// 41
	// fizz
	// 43
	// 44
	// fizzbuzz
	// 46
	// 47
	// fizz
	// 49
	// buzz
	// fizz
	// 52
	// 53
	// fizz
	// buzz
	// 56
	// fizz
	// 58
	// 59
	// fizzbuzz
	// 61
	// 62
	// fizz
	// 64
	// buzz
	// fizz
	// 67
	// 68
	// fizz
	// buzz
	// 71
	// fizz
	// 73
	// 74
	// fizzbuzz
	// 76
	// 77
	// fizz
	// 79
	// buzz
	// fizz
	// 82
	// 83
	// fizz
	// buzz
	// 86
	// fizz
	// 88
	// 89
	// fizzbuzz
	// 91
	// 92
	// fizz
	// 94
	// buzz
	// fizz
	// 97
	// 98
	// fizz
}

func ExampleSomeMath() {
	if err := RunProgram("somemath", sampleprograms.SomeMath); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Add: 3
	// Sub: -1
	// Mul: 6
	// Div: 3
	// Complex: 5
}

func ExampleEqualComparison() {
	if err := RunProgram("equalcompare", sampleprograms.EqualComparison); err != nil {
		fmt.Println(err.Error())
	}
	// Output: true
	// 3
}

func ExampleNotEqualComparison() {
	if err := RunProgram("notequalcompare", sampleprograms.NotEqualComparison); err != nil {
		fmt.Println(err.Error())
	}
	// Output: false
}

func ExampleGreaterComparison() {
	if err := RunProgram("greatercompare", sampleprograms.GreaterComparison); err != nil {
		fmt.Println(err.Error())
	}
	// Output: true
	// 4
}

func ExampleGreaterOrEqualComparison() {
	if err := RunProgram("greaterorequalcompare", sampleprograms.GreaterOrEqualComparison); err != nil {
		fmt.Println(err.Error())
	}
	// Output: true
	// 4
	// 3
}

func ExampleLessThanEqualComparison() {
	if err := RunProgram("lessthanequalcompare", sampleprograms.LessThanComparison); err != nil {
		fmt.Println(err.Error())
	}
	// Output: false
}

func ExampleLessThanOrEqualComparison() {
	if err := RunProgram("lessthanorequalcompare", sampleprograms.LessThanOrEqualComparison); err != nil {
		fmt.Println(err.Error())
	}
	// Output: true
	// 1
	// 2
	// 3
}

func ExampleUserDefinedType() {
	if err := RunProgram("usertype", sampleprograms.UserDefinedType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleTypeInference() {
	if err := RunProgram("typeinference", sampleprograms.TypeInference); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 0, 4
}

func ExampleConcreteUint8() {
	if err := RunProgram("concreteuint8", sampleprograms.ConcreteTypeUint8); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleConcreteInt8() {
	if err := RunProgram("concreteint8", sampleprograms.ConcreteTypeInt8); err != nil {
		fmt.Println(err.Error())
	}
	// Output: -4
}

func ExampleConcreteUint16() {
	if err := RunProgram("concreteuint16", sampleprograms.ConcreteTypeUint16); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleConcreteInt16() {
	if err := RunProgram("concreteint16", sampleprograms.ConcreteTypeInt16); err != nil {
		fmt.Println(err.Error())
	}
	// Output: -4
}

func ExampleConcreteUint32() {
	if err := RunProgram("concreteuint32", sampleprograms.ConcreteTypeUint32); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleConcreteInt32() {
	if err := RunProgram("concreteint32", sampleprograms.ConcreteTypeInt32); err != nil {
		fmt.Println(err.Error())
	}
	// Output: -4
}

func ExampleConcreteUint64() {
	if err := RunProgram("concreteuint64", sampleprograms.ConcreteTypeUint64); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleConcreteInt64() {
	if err := RunProgram("concreteint64", sampleprograms.ConcreteTypeInt64); err != nil {
		fmt.Println(err.Error())
	}
	// Output: -4
}

func ExampleFibonacci() {
	if err := RunProgram("fibonacci", sampleprograms.Fibonacci); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 2
	// 3
	// 5
	// 8
	// 13
	// 21
	// 34
	// 55
	// 89
	// 144
}

func ExampleEnumType() {
	if err := RunProgram("enumtype", sampleprograms.EnumType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: I am A!
}

func ExampleEnumTypeInferred() {
	if err := RunProgram("enumtypeinferred", sampleprograms.EnumTypeInferred); err != nil {
		fmt.Println(err.Error())
	}
	// Output: I am B!
}

func ExampleSimpleMatch() {
	if err := RunProgram("simplematch", sampleprograms.SimpleMatch); err != nil {
		fmt.Println(err.Error())
	}
	// Output: I am 3
}

func ExampleIfElseMatch() {
	if err := RunProgram("ifelsematch", sampleprograms.IfElseMatch); err != nil {
		fmt.Println(err.Error())
	}
	// Output: x is less than 4
}

func ExampleGenericEnumType() {
	if err := RunProgram("genericenumtype", sampleprograms.GenericEnumType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 5
	// I am nothing!
}

func ExampleMatchParam() {
	if err := RunProgram("matchparam", sampleprograms.MatchParam); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 5
}

func ExampleMatchParam2() {
	if err := RunProgram("matchparam2", sampleprograms.MatchParam2); err != nil {
		fmt.Println(err.Error())
	}
	// Output: x5
}

func ExampleSimpleAlgorithm() {
	if err := RunProgram("simplealgorithm", sampleprograms.SimpleAlgorithm); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 180
}

func ExampleSimpleArray() {
	if err := RunProgram("simplearray", sampleprograms.SimpleArray); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleMutationArray() {
	if err := RunProgram("arraymutation", sampleprograms.ArrayMutation); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
	// 2
	// 3
}

func ExampleReference() {
	if err := RunProgram("referencevariable", sampleprograms.ReferenceVariable); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 3
	// 4
	// 7
}

func ExampleSimpleSlice() {
	if err := RunProgram("simpleslice", sampleprograms.SimpleSlice); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleSimpleSliceInference() {
	if err := RunProgram("simplesliceinference", sampleprograms.SimpleSliceInference); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleSliceMutation() {
	if err := RunProgram("slicemutation", sampleprograms.SliceMutation); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
	// 2
	// 3
}

func ExampleSliceParam() {
	if err := RunProgram("sliceparam", sampleprograms.SliceParam); err != nil {
		fmt.Println(err.Error())
	}
	// Output: ,7X
}

func ExamplePrintString() {
	if err := RunProgram("printstring", sampleprograms.PrintString); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Success!
}

func ExampleWriteSyscall() {
	if err := RunProgram("writestring", sampleprograms.WriteSyscall); err != nil {
		fmt.Println(err.Error())
	}
	// Example output only compares stdout, not stderr, so the Write(2, "Stderr!") call
	// isn't automatically tested.

	// Output: Stdout!
}

func ExampleSliceLength() {
	if err := RunProgram("slicelength", sampleprograms.SliceLength); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
}

func ExampleSliceLength2() {
	if err := RunProgram("slicelength2", sampleprograms.SliceLength2); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 5
}

func ExampleArrayIndex() {
	if err := RunProgram("arrayindex", sampleprograms.ArrayIndex); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
	// 5
}

func ExampleIndexAssignment() {
	if err := RunProgram("indexassignment", sampleprograms.IndexAssignment); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 4
	// 5
}

func ExampleIndexedAddition() {
	if err := RunProgram("indexaddition", sampleprograms.IndexedAddition); err != nil {
		fmt.Println(err.Error())
	}
	// Output: 9
	// 8
}
