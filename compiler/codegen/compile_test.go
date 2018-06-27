package codegen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/driusan/lang/compiler/mlir"
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

func compileAndTest(t *testing.T, name, p string, stdout, stderr string) error {
	dir, err := ioutil.TempDir("", "langtest"+name)
	if err != nil {
		return err
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	exe, err := BuildProgram(dir, strings.NewReader(p))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(dir + "/" + exe)
	cstdout := &strings.Builder{}
	cstderr := &strings.Builder{}
	cmd.Stdout = cstdout
	cmd.Stderr = cstderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if stdout != cstdout.String() {
		t.Errorf("Unexpected standard out for %v: want %v got %v", name, stdout, cstdout.String())
	}
	if stderr != cstderr.String() {
		t.Errorf("Unexpected standard err for %v: want %v got %v", name, stderr, cstderr.String())
	}
	return err
}

func runTest(t *testing.T, filename string, stdout, stderr string) error {
	t.Helper()
	dir, err := ioutil.TempDir("", "langtest"+filename)
	if err != nil {
		t.Fatalf("%v: %v", filename, err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	f, err := os.Open("../../testsuite/" + filename + ".l")
	if err != nil {
		t.Fatalf("%v: %v", filename, err)
	}
	exe, err := BuildProgram(dir, f)
	if err != nil {
		t.Fatalf("%v: %v", filename, err)
	}
	cmd := exec.Command(dir + "/" + exe)
	cstdout := &strings.Builder{}
	cstderr := &strings.Builder{}
	cmd.Stdout = cstdout
	cmd.Stderr = cstderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	if stdout != cstdout.String() {
		t.Errorf("Unexpected standard out for %v: want %v got %v", filename, stdout, cstdout.String())
	}
	if stderr != cstderr.String() {
		t.Errorf("Unexpected standard err for %v: want %v got %v", filename, stderr, cstderr.String())
	}
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
	fnc, _, err := mlir.Generate(prgAst[0], types, callables, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := Compile(&w, fnc); err != nil {
		t.Fatal(err)
	}

	// HelloWorld is simple enough that we can hand compile it. Most other
	// programs we just compile it, run it and check the output to ensure
	// that they work.
	expected := `TEXT main(SB), 4+16, $24
	DATA string0<>+0(SB)/8, $14
	DATA string0<>+8(SB)/8, $"Hello, w"
	DATA string0<>+16(SB)/8, $"orld!\n\000\000"
	GLOBL string0<>+0(SB), 8+16, $24
	MOVQ $14, 0(SP)
	MOVQ $string0<>+8(SB), BX
	MOVQ BX, 8(SP)
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
	// We chdir below, defer a cleanup that resets it after the teste finishes.
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

// Test that Open/Read work correctly. This isn't done as an Example test because we
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
	if !debug {
		defer os.RemoveAll(dir)
	}

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

// Echo is the simplest program that takes arguments
func TestEchoProgram(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	mlir.Debug = false
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)
	// Inline the RunProgram call, because we want to cd to the directory so that the foo.txt
	// file isn't created as garbage in the cwd..
	dir, err := ioutil.TempDir("", "langtestEcho")
	if err != nil {
		t.Fatal(err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}

	// There's currently no way to define constants, no bitwise operations, and only decimal
	// numbers, so modes and flags are hardcoded in decimal.
	exe, err := BuildProgram(dir, strings.NewReader(sampleprograms.Echo))
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(dir+"/"+exe, "foo", "bar")

	content, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected := "foo bar\n"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}

	cmd = exec.Command(dir+"/"+exe, "other", "params")

	content, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected = "other params\n"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}
}

func TestCatProgram(t *testing.T) {
	mlir.Debug = false
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

	// There's currently no way to define constants, no bitwise operations, and only decimal
	// numbers, so modes and flags are hardcoded in decimal.
	exe, err := BuildProgram(dir, strings.NewReader(sampleprograms.UnbufferedCat))
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(dir+"/"+exe, "foo.tmp", "bar.tmp")

	content, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected := "FooBar"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}

	cmd = exec.Command(dir+"/"+exe, "bar.tmp", "foo.tmp")

	content, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected = "BarFoo"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}
}

func TestUnbufferedCat2(t *testing.T) {
	mlir.Debug = false
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
	dir, err := ioutil.TempDir("", "langtestunbufferedcat2")
	if err != nil {
		t.Fatal(err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// There's currently no way to define constants, no bitwise operations, and only decimal
	// numbers, so modes and flags are hardcoded in decimal.
	exe, err := BuildProgram(dir, strings.NewReader(sampleprograms.UnbufferedCat2))
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(dir+"/"+exe, "foo.tmp", "bar.tmp")

	content, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected := "FooBar"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}

	cmd = exec.Command(dir+"/"+exe, "bar.tmp", "foo.tmp")

	content, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected = "BarFoo"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}
}

func TestUnbufferedCat3(t *testing.T) {
	mlir.Debug = false
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
	dir, err := ioutil.TempDir("", "langtestunbufferedcat3")
	if err != nil {
		t.Fatal(err)
	}
	if !debug {
		defer os.RemoveAll(dir)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// There's currently no way to define constants, no bitwise operations, and only decimal
	// numbers, so modes and flags are hardcoded in decimal.
	exe, err := BuildProgram(dir, strings.NewReader(sampleprograms.UnbufferedCat3))
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(dir+"/"+exe, "foo.tmp", "bar.tmp")

	content, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected := "FooBar"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}

	cmd = exec.Command(dir+"/"+exe, "bar.tmp", "foo.tmp")

	content, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	expected = "BarFoo"
	if got := string(content); got != expected {
		t.Errorf("Unexpected value: got %v want %v", got, expected)
	}
}

func TestAssertions(t *testing.T) {
	tests := []struct {
		Name           string
		Program        string
		Stdout, Stderr string
	}{
		{
			"AssertionFail", sampleprograms.AssertionFail, "", "assertion false failed",
		},
		{
			"AssertionFailWithMessage", sampleprograms.AssertionFailWithMessage, "", "assertion false failed: This always fails",
		},
		{
			"AssertionPass", sampleprograms.AssertionPass, "", "",
		},
		{
			"AssertionPassWithMessage", sampleprograms.AssertionPassWithMessage, "", "",
		},
		{
			"AssertionFailWithVariable", sampleprograms.AssertionFailWithVariable, "", "assertion x > 3 failed",
		},
	}

	for _, test := range tests {
		compileAndTest(t, test.Name, test.Program, test.Stdout, test.Stderr)
	}

}

// TestTestSuite tests everything in the testsuite
func TestTestSuite(t *testing.T) {
	tests := []struct {
		Name           string
		Stdout, Stderr string
	}{
		{"emptymain", "", ""},
		{"emptyreturn", "", ""},
		{"helloworld", "Hello, world!\n", ""},
		{"letstatement", "5", ""},
		{"letstatementshadow", "5\nhello", ""},
		{"twoprocs", "3", ""},
		{"outoforder", "3", ""},
		{"mutaddition", "8", ""},
		{"sumtoten", "55", ""},
		{"sumtotenrecursive", "55", ""},
		{"fizzbuzz", "1\n2\nfizz\n4\nbuzz\nfizz\n7\n8\nfizz\nbuzz\n11\nfizz\n13\n14\nfizzbuzz\n16\n17\nfizz\n19\nbuzz\nfizz\n22\n23\nfizz\nbuzz\n26\nfizz\n28\n29\nfizzbuzz\n31\n32\nfizz\n34\nbuzz\nfizz\n37\n38\nfizz\nbuzz\n41\nfizz\n43\n44\nfizzbuzz\n46\n47\nfizz\n49\nbuzz\nfizz\n52\n53\nfizz\nbuzz\n56\nfizz\n58\n59\nfizzbuzz\n61\n62\nfizz\n64\nbuzz\nfizz\n67\n68\nfizz\nbuzz\n71\nfizz\n73\n74\nfizzbuzz\n76\n77\nfizz\n79\nbuzz\nfizz\n82\n83\nfizz\nbuzz\n86\nfizz\n88\n89\nfizzbuzz\n91\n92\nfizz\n94\nbuzz\nfizz\n97\n98\nfizz\n", ""},
		{"somemath", "Add: 3\nSub: -1\nMul: 6\nDiv: 3\nComplex: 5\n", ""},
		{"equalcomparison", "true\n3\n", ""},
		{"notequalcomparison", "", ""},
		{"greatercomparison", "4\n", ""},
		{"greaterorequalcomparison", "true\n4\n3\n", ""},
		{"lessthancomparison", "false\n", ""},
		{"lessthanorequalcomparison", "true\n1\n2\n3\n", ""},
		{"userdefinedtype", "4", ""},
		{"typeinference", "0, 4\n", ""},
		{"concreteuint8", "4", ""},
		{"concreteint8", "-4", ""},
		{"concreteuint16", "4", ""},
		{"concreteint16", "-4", ""},
		{"concreteuint32", "4", ""},
		{"concreteint32", "-4", ""},
		{"concreteuint64", "4", ""},
		{"concreteint64", "-4", ""},
		{"fibonacci", "2\n3\n5\n8\n13\n21\n34\n55\n89\n144\n", ""},
		{"enumtype", "I am A!\n", ""},
		{"enumtypeinferred", "I am B!\n", ""},
		{"simplematch", "I am 3\n", ""},
		{"ifelsematch", "x is less than 4\n", ""},
		{"genericenumtype", "5\nI am nothing!\n", ""},
		{"matchparam", "5", ""},
		{"matchparam2", "x5", ""},
		{"simplealgorithm", "180", ""},
		{"simplearray", "4", ""},
		{"arraymutation", "4\n2\n3", ""},
		{"referencevariable", "3\n4\n7", ""},
		{"simpleslice", "4", ""},
		{"simplesliceinference", "4", ""},
		{"slicemutation", "4\n2\n3", ""},
		{"sliceparam", ",7X", ""},
		{"slicestringparam", "bar", ""},
		{"slicestringvariableparam", "bar", ""},
		{"printstring", "Success!", ""},
		{"write", "Stdout!", "Stderr!"},
		{"slicelength", "4", ""},
		{"slicelength2", "5", ""},
		{"arrayindex", "4\n5", ""},
		{"indexassignment", "4\n5", ""},
		{"indexedaddition", "9\n8", ""},
		{"stringarray", "bar\nfoo", ""},
		{"preecho", "bar baz\n", ""},
		{"preecho2", "bar baz\n", ""},
		{"precedence", "-3", ""},
		{"letcondition", "1-112", ""},
		{"methodsyntax", "10", ""},
		{"assignmenttoconstantindex", "365", ""},
		{"assignmenttovariableindex", "64", ""},
		{"assignmenttosliceconstantindex", "365", ""},
		{"assignmenttoslicevariableindex", "64", ""},
		{"writestringbyte", "hellohello", ""},
		{"stringarg", "foobar", ""},
		{"castbuiltin", "Foo", ""},
		{"castbuiltin2", "bar", ""},
		{"castintvariable", "65", ""},
		{"sumtypefunccall", "bar3", ""},
		{"sumtypefuncreturn", "not33", ""},
		{"ifbool", "73", ""},
		{"linecomment", "3", ""},
		{"blockcomment", "3", ""},
		{"producttypevalue", "3\n0", ""},
		{"userproducttypevalue", "hello\n3", ""},
		{"slicefromarray", "34", ""},
		{"slicefromslice", "45", ""},
		{"mutslicefromarray", "34", ""},
		{"mutslicefromslice", "45", ""},
		{"sliceprint", "Foo", "Bar"},
		{"arrayparam", "16", ""},
		{"mutarrayparam", "", ""},
	}

	for _, tst := range tests {
		if testing.Verbose() {
			println("\tRunning test: ", tst.Name)
		}
		runTest(t, tst.Name, tst.Stdout, tst.Stderr)
	}
}
