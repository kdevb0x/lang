package vm

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/driusan/lang/compiler/hlir"
)

func compileAndTestWithArgs(t *testing.T, name string, prog io.Reader, args []string, estdout, estderr string) {
	t.Helper()
	ctx, err := ParseFromReader(prog)
	if err != nil {
		t.Fatal(err)
	}

	args = append([]string{name}, args...)

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
		if _, assert := err.(assertionError); !assert {
			t.Fatal(err)
		}
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
		t.Errorf("Unexpected stderr: got %s want %s", stde, estderr)
	}
}

func compileAndTestFromFile(t *testing.T, prog, estdout, estderr string, setup func()) {
	t.Helper()
	f, err := os.Open("../../../testsuite/" + prog + ".l")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	setup()

	compileAndTestWithArgs(t, prog, f, nil, estdout, estderr)
}

func compileAndTestFromFileWithArgs(t *testing.T, prog string, args []string, estdout, estderr string, setup func()) {
	t.Helper()
	f, err := os.Open("../../../testsuite/" + prog + ".l")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	setup()

	compileAndTestWithArgs(t, prog, f, args, estdout, estderr)
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
	compileAndTestFromFile(t, "createsyscall", "", "", func() {
		// Make sure foo.txt gets created in dir, so that the defer cleans it up..
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}

	})

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

	compileAndTestFromFile(t, "readsyscall", "Hello,", "", func() {
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile("foo.txt", []byte("Hello, world!"), 0755); err != nil {
			t.Fatal(err)
		}

	})

	os.Chdir(pwd)

	compileAndTestFromFile(t, "readsyscall", "Goodby", "", func() {
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		// Run it again with different file content.
		if err := ioutil.WriteFile("foo.txt", []byte("Goodbye"), 0755); err != nil {
			t.Fatal(err)
		}
	})
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

	compileAndTestFromFileWithArgs(t, "echo", []string{"foo", "bar"}, "foo bar\n", "", func() {})
	compileAndTestFromFileWithArgs(t, "echo", []string{"other", "params"}, "other params\n", "", func() {})
}

func TestCatPrograms(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	cats := []string{"unbufferedcat", "unbufferedcat2", "unbufferedcat3"}
	for _, c := range cats {
		// We need to be in the original directory for loading the test to work.
		if err := os.Chdir(pwd); err != nil {
			t.Fatal(err)
		}
		// Set up a test directory
		dir, err := ioutil.TempDir("", "langtestCat")
		if err != nil {
			t.Fatal(err)
		}
		if !debug {
			defer os.RemoveAll(dir)
		}

		setup := func() {
			if err := os.Chdir(pwd); err != nil {
				t.Fatal(err)
			}
			if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
				t.Fatal(err)
			}
			if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
				t.Fatal(err)
			}
		}

		compileAndTestFromFileWithArgs(t, c, []string{"foo.tmp", "bar.tmp"}, "FooBar", "", setup)
		compileAndTestFromFileWithArgs(t, c, []string{"bar.tmp", "foo.tmp"}, "BarFoo", "", setup)
	}
}

func TestAssertions(t *testing.T) {
	tests := []struct {
		Name           string
		Stdout, Stderr string
	}{
		{
			"AssertionFail", "", "assertion false failed",
		},
		{
			"AssertionFailWithMessage", "", "assertion false failed: This always fails",
		},
		{
			"AssertionPass", "", "",
		},
		{
			"AssertionPassWithMessage", "", "",
		},
		{
			"AssertionFailWithVariable", "", "assertion x > 3 failed",
		},
		{
			// The first PrintInt should succeed, the second not.
			"AssertionFailWithContext", "0", "assertion false failed",
		},
	}

	for _, tc := range tests {
		if testing.Verbose() {
			println("\tRunning test: ", tc.Name)
		}
		compileAndTestFromFile(t, tc.Name, tc.Stdout, tc.Stderr, func() {})
	}
}

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
		{"swap", "", ""},
		{"reverse", "", ""},
		{"digitsinto", "", ""},
	}

	for _, tc := range tests {
		compileAndTestFromFile(t, tc.Name, tc.Stdout, tc.Stderr, func() {})
	}
}
