package vm

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/sampleprograms"
)

func compileAndTestWithArgs(t *testing.T, prog string, args []string, estdout, estderr string) {
	ctx, err := Parse(prog)
	if err != nil {
		t.Fatal(err)
	}

	args = append([]string{prog}, args...)

	// Create a new environment to put the StringLiterals in, so that they can be
	// treated the same as any other call in how they're passed
	env := NewContext()
	env.stdout = ctx.stdout
	env.stderr = ctx.stderr

	for i, s := range args {
		env.localValues[hlir.LocalValue(2*i)] = len([]byte(s))
		env.localValues[hlir.LocalValue(2*i+1)] = s //hlir.StringLiteral(s)

	}
	ctx.pointers[hlir.Pointer{hlir.FuncArg{1, false}}] = Pointer{hlir.Pointer{hlir.LocalValue(0)}, env}

	ctx.funcArg[hlir.FuncArg{0, false}] = int(len(args))
	//ctx.funcArg[hlir.FuncArg{uint(1), false}] = hlir.Pointer{hlir.FuncArg{1, false}}

	stdout, stderr, err := RunWithSideEffects("main", ctx)
	if err != nil {
		t.Fatal(err)
	}
	stdo, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Error(err)
	}
	if string(stdo) != estdout {
		t.Errorf("Unexpected stdout: got %s want %s", stdo, estdout)
	}
	stde, err := ioutil.ReadAll(stderr)
	if err != nil {
		t.Error(err)
	}
	if string(stde) != estderr {
		t.Errorf("Unexpected stdeut: got %s want %s", stde, estderr)
	}
}

func compileAndTest(t *testing.T, prog, estdout, estderr string) {
	compileAndTestWithArgs(t, prog, nil, estdout, estderr)
}

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
	if !debug {
		defer os.RemoveAll(dir)
	}

	// Make sure foo.txt gets created in dir, so that the defer cleans it up..
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	compileAndTest(t, sampleprograms.CreateSyscall, "", "")

	content, err := ioutil.ReadFile("foo.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	if string(content) != "Hello\n" {
		t.Errorf("Unexpected content of file foo.txt: got %v want %v", string(content), "Hello\n")
	}
}

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

	if !debug {
		defer os.RemoveAll(dir)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("foo.txt", []byte("Hello, world!"), 0755); err != nil {
		t.Fatal(err)
	}

	compileAndTest(t, sampleprograms.ReadSyscall, "Hello,", "")

	// Run it again with different file content.
	if err := ioutil.WriteFile("foo.txt", []byte("Goodbye"), 0755); err != nil {
		t.Fatal(err)
	}
	compileAndTest(t, sampleprograms.ReadSyscall, "Goodby", "")

}

// Echo is the simplest program that takes arguments
func TestEchoProgram(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	compileAndTestWithArgs(t, sampleprograms.Echo, []string{"foo", "bar"}, "foo bar\n", "")
	compileAndTestWithArgs(t, sampleprograms.Echo, []string{"other", "params"}, "other params\n", "")
}

func TestCatProgram(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	// Set up a test directory
	dir, err := ioutil.TempDir("", "langtestCat")
	if err != nil {
		t.Fatal(err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
		t.Fatal(err)
	}

	compileAndTestWithArgs(t, sampleprograms.UnbufferedCat, []string{"foo.tmp", "bar.tmp"}, "FooBar", "")
	compileAndTestWithArgs(t, sampleprograms.UnbufferedCat, []string{"bar.tmp", "foo.tmp"}, "BarFoo", "")

}

func TestUnbufferedCat2(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	// Set up a test directory
	dir, err := ioutil.TempDir("", "langtestCat")
	if err != nil {
		t.Fatal(err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
		t.Fatal(err)
	}

	compileAndTestWithArgs(t, sampleprograms.UnbufferedCat2, []string{"foo.tmp", "bar.tmp"}, "FooBar", "")
	compileAndTestWithArgs(t, sampleprograms.UnbufferedCat2, []string{"bar.tmp", "foo.tmp"}, "BarFoo", "")

}

func TestUnbufferedCat3(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	// Set up a test directory
	dir, err := ioutil.TempDir("", "langtestCat")
	if err != nil {
		t.Fatal(err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
		t.Fatal(err)
	}

	compileAndTestWithArgs(t, sampleprograms.UnbufferedCat3, []string{"foo.tmp", "bar.tmp"}, "FooBar", "")
	compileAndTestWithArgs(t, sampleprograms.UnbufferedCat3, []string{"bar.tmp", "foo.tmp"}, "BarFoo", "")

}

func TestEmptyMain(t *testing.T) {
	compileAndTest(t, sampleprograms.EmptyMain, "", "")
}

func TestHelloWorld(t *testing.T) {
	compileAndTest(t, sampleprograms.HelloWorld, "Hello, world!\n", "")
}

func TestGenLetStatement(t *testing.T) {
	compileAndTest(t, sampleprograms.LetStatement, "5", "")
}

func TestGenLetStatementShadow(t *testing.T) {
	compileAndTest(t, sampleprograms.LetStatementShadow, "5\nhello", "")
}

func TestGenTwoProcs(t *testing.T) {
	compileAndTest(t, sampleprograms.TwoProcs, "3", "")
}

func TestOutOfOrder(t *testing.T) {
	compileAndTest(t, sampleprograms.OutOfOrder, "3", "")
}

func TestMutAddition(t *testing.T) {
	compileAndTest(t, sampleprograms.MutAddition, "8", "")
}

func TestSimpleFunc(t *testing.T) {
	compileAndTest(t, sampleprograms.SimpleFunc, "3", "")
}

func TestSumToTen(t *testing.T) {
	compileAndTest(t, sampleprograms.SumToTen, "55", "")
}

func TestSumToTenRecursive(t *testing.T) {
	compileAndTest(t, sampleprograms.SumToTenRecursive, "55", "")
}

func TestFizzBuzz(t *testing.T) {
	compileAndTest(t, sampleprograms.Fizzbuzz, `1
2
fizz
4
buzz
fizz
7
8
fizz
buzz
11
fizz
13
14
fizzbuzz
16
17
fizz
19
buzz
fizz
22
23
fizz
buzz
26
fizz
28
29
fizzbuzz
31
32
fizz
34
buzz
fizz
37
38
fizz
buzz
41
fizz
43
44
fizzbuzz
46
47
fizz
49
buzz
fizz
52
53
fizz
buzz
56
fizz
58
59
fizzbuzz
61
62
fizz
64
buzz
fizz
67
68
fizz
buzz
71
fizz
73
74
fizzbuzz
76
77
fizz
79
buzz
fizz
82
83
fizz
buzz
86
fizz
88
89
fizzbuzz
91
92
fizz
94
buzz
fizz
97
98
fizz
`, "")
}

func TestSomeMath(t *testing.T) {
	compileAndTest(t, sampleprograms.SomeMath,
		`Add: 3
Sub: -1
Mul: 6
Div: 3
Complex: 5
`, "")
}

func TestEqualComparison(t *testing.T) {
	compileAndTest(t, sampleprograms.EqualComparison,
		`true
3
`, "")
}

func TestNotEqualComparison(t *testing.T) {
	compileAndTest(t, sampleprograms.NotEqualComparison, "false\n", "")
}

func TestGreaterComparison(t *testing.T) {
	compileAndTest(t, sampleprograms.GreaterComparison, "true\n4\n", "")
}

func TestGreaterOrEqualComparison(t *testing.T) {
	compileAndTest(t, sampleprograms.GreaterOrEqualComparison, "true\n4\n3\n", "")
}

func TestLessThanComparison(t *testing.T) {
	compileAndTest(t, sampleprograms.LessThanComparison, "false\n", "")
}

func TestLessThanOrEqualComparison(t *testing.T) {
	compileAndTest(t, sampleprograms.LessThanOrEqualComparison, "true\n1\n2\n3\n", "")
}

func TestUserDefinedType(t *testing.T) {
	compileAndTest(t, sampleprograms.UserDefinedType, "4", "")
}

func TestTypeInference(t *testing.T) {
	compileAndTest(t, sampleprograms.TypeInference, "0, 4\n", "")
}

func TestConcreteTypeUint8(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeUint8, "4", "")
}

func TestConcreteTypeInt8(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeInt8, "-4", "")
}

func TestConcreteTypeUint16(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeUint16, "4", "")
}

func TestConcreteTypeInt16(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeInt16, "-4", "")
}

func TestConcreteTypeUint32(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeUint32, "4", "")
}

func TestConcreteTypeInt32(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeInt32, "-4", "")
}

func TestConcreteTypeUint64(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeUint64, "4", "")
}

func TestConcreteTypeInt64(t *testing.T) {
	compileAndTest(t, sampleprograms.ConcreteTypeInt64, "-4", "")
}

func TestFibonacci(t *testing.T) {
	compileAndTest(t, sampleprograms.Fibonacci, `2
3
5
8
13
21
34
55
89
144
`, "")
}

func TestEnumType(t *testing.T) {
	compileAndTest(t, sampleprograms.EnumType, "I am A!\n", "")
}

func TestEnumTypeInferred(t *testing.T) {
	compileAndTest(t, sampleprograms.EnumTypeInferred, "I am B!\n", "")
}

func TestSimpleMatch(t *testing.T) {
	compileAndTest(t, sampleprograms.SimpleMatch, "I am 3\n", "")
}

func TestIfElseMatch(t *testing.T) {
	compileAndTest(t, sampleprograms.IfElseMatch, "x is less than 4\n", "")
}

func TestGenericEnumType(t *testing.T) {
	compileAndTest(t, sampleprograms.GenericEnumType, "5\nI am nothing!\n", "")
}

func TestMatchParam(t *testing.T) {
	compileAndTest(t, sampleprograms.MatchParam, "5", "")
}

func TestMatchParam2(t *testing.T) {
	compileAndTest(t, sampleprograms.MatchParam2, "x5", "")
}

func TestSimpleAlgorithm(t *testing.T) {
	compileAndTest(t, sampleprograms.SimpleAlgorithm, "180", "")
}

func TestSimpleArray(t *testing.T) {
	compileAndTest(t, sampleprograms.SimpleArray, "4", "")
}

func TestArrayMutation(t *testing.T) {
	compileAndTest(t, sampleprograms.ArrayMutation, "4\n2\n3", "")
}

func TestReferenceVariable(t *testing.T) {
	compileAndTest(t, sampleprograms.ReferenceVariable, "3\n4\n7", "")
}

func TestSimpleSlice(t *testing.T) {
	compileAndTest(t, sampleprograms.SimpleSlice, "4", "")
}

func TestSimpleSliceInference(t *testing.T) {
	compileAndTest(t, sampleprograms.SimpleSliceInference, "4", "")
}

func TestSliceMutation(t *testing.T) {
	compileAndTest(t, sampleprograms.SliceMutation, "4\n2\n3", "")
}

func TestSliceParam(t *testing.T) {
	compileAndTest(t, sampleprograms.SliceParam, ",7X", "")
}

func TestSliceStringParam(t *testing.T) {
	compileAndTest(t, sampleprograms.SliceStringParam, "bar", "")
}

func TestSliceStringVariableParam(t *testing.T) {
	compileAndTest(t, sampleprograms.SliceStringVariableParam, "bar", "")
}

func TestPrintString(t *testing.T) {
	compileAndTest(t, sampleprograms.PrintString, "Success!", "")
}

func TestWriteSyscall(t *testing.T) {
	compileAndTest(t, sampleprograms.WriteSyscall, "Stdout!", "Stderr!")
}

func TestSliceLength(t *testing.T) {
	compileAndTest(t, sampleprograms.SliceLength2, "5", "")
}

func TestArrayIndex(t *testing.T) {
	compileAndTest(t, sampleprograms.ArrayIndex, "4\n5", "")
}

func TestIndexAssignment(t *testing.T) {
	compileAndTest(t, sampleprograms.IndexAssignment, "4\n5", "")
}

func TestIndexedAddition(t *testing.T) {
	compileAndTest(t, sampleprograms.IndexedAddition, "9\n8", "")
}

func TestStringArray(t *testing.T) {
	compileAndTest(t, sampleprograms.StringArray, "bar\nfoo", "")
}

func TestPreEcho(t *testing.T) {
	compileAndTest(t, sampleprograms.PreEcho, "bar baz\n", "")
}

func TestPreEcho2(t *testing.T) {
	compileAndTest(t, sampleprograms.PreEcho2, "bar baz\n", "")
}

func TestPrecedence(t *testing.T) {
	compileAndTest(t, sampleprograms.Precedence, "-3", "")
}

func TestLetCondition(t *testing.T) {
	compileAndTest(t, sampleprograms.LetCondition, "1-112", "")
}

func TestMethodSyntax(t *testing.T) {
	compileAndTest(t, sampleprograms.MethodSyntax, "10", "")
}

func TestAssignmentToConstantIndex(t *testing.T) {
	compileAndTest(t, sampleprograms.AssignmentToConstantIndex, "365", "")
}

func TestAssignmentToVariableIndex(t *testing.T) {
	compileAndTest(t, sampleprograms.AssignmentToVariableIndex, "64", "")
}

func TestAssignmentToSliceVariableIndex(t *testing.T) {
	compileAndTest(t, sampleprograms.AssignmentToSliceVariableIndex, "64", "")
}

func TestWriteStringByte(t *testing.T) {
	compileAndTest(t, sampleprograms.WriteStringByte, "hellohello", "")
}

func TestStringArg(t *testing.T) {
	compileAndTest(t, sampleprograms.StringArg, "foobar", "")
}

func TestCastBuiltin(t *testing.T) {
	compileAndTest(t, sampleprograms.CastBuiltin, "Foo", "")
}

func TestCastBuiltin2(t *testing.T) {
	compileAndTest(t, sampleprograms.CastBuiltin2, "bar", "")
}

func TestCastIntVariable(t *testing.T) {
	compileAndTest(t, sampleprograms.CastIntVariable, "65", "")
}

func TestEmptyReturn(t *testing.T) {
	compileAndTest(t, sampleprograms.EmptyReturn, "", "")
}
