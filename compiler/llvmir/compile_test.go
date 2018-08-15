package llvmir

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func loadTestFile(t *testing.T, testcase string) io.ReadCloser {
	t.Helper()
	f, err := os.Open("../../testsuite/" + testcase + ".l")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func runWithArgs(t *testing.T, name, dir string, args []string, estdout, estderr string) {
	t.Helper()
	cmd := exec.Command(dir+"/main", args...)
	cstdout := &strings.Builder{}
	cstderr := &strings.Builder{}
	cmd.Stdout = cstdout
	cmd.Stderr = cstderr

	if err := cmd.Start(); err != nil {
		t.Errorf("%v: %v", name, err)
		return
	}
	err := cmd.Wait()
	switch e := err.(type) {
	case nil:
		// Success
	case *exec.ExitError:
		if e.Exited() {
			// If it exited with an error code, we don't consider it an
			// error
		} else {
			t.Errorf("%v: %v", name, err)
		}
	default:
		// cmd.Wait is documented as returning an *exec.ExitError, so if anything
		// else happened it's fatal.
		t.Fatal(e)
	}

	if estdout != cstdout.String() {
		t.Errorf("Unexpected stdout for %v: got %v want %v", name, cstdout.String(), estdout)
	}
	if estderr != cstderr.String() {
		t.Errorf("Unexpected stderr for %v: got %v want %v", name, cstderr.String(), estderr)
	}
}

func compileAndRun(t *testing.T, name string, estdout, estderr string) {
	t.Helper()
	r := loadTestFile(t, name)
	defer r.Close()

	dir, err := Compile(r, false)
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	if err != nil {
		t.Fatal(err)
	}
	runWithArgs(t, name, dir, nil, estdout, estderr)
}

// TestAssertions tests that assertions work in various contexts.
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

	for _, tst := range tests {
		if testing.Verbose() {
			println("\tRunning test: ", tst.Name)
		}
		compileAndRun(t, tst.Name, tst.Stdout, tst.Stderr)
	}
}

// TestTestSuite tests everything in the testsuite
func TestTestSuite(t *testing.T) {
	tests := []struct {
		Name           string
		Stdout, Stderr string
	}{
		/*
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
		{"enumarray", "", ""},
		{"enumarrayexplicit", "", ""},
		{"sumtypearray", "", ""},
		*/
		{"sumtypearraymatch", "stringy33", ""},
		// needs "default" keyword and variable slice indexing
		//{"printint", "67\n-568", ""},
	}

	for _, tst := range tests {
		if testing.Verbose() {
			println("\tRunning test: ", tst.Name)
		}
		compileAndRun(t, tst.Name, tst.Stdout, tst.Stderr)
	}
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

	r := loadTestFile(t, "createsyscall")
	dir, err := Compile(r, false)
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	if err != nil {
		t.Fatal(err)
	}

	// Make sure foo gets created in dir, so that the defer cleans it up.
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	runWithArgs(t, "createsyscall", dir, nil, "", "")

	content, err := ioutil.ReadFile("foo.txt")
	if err != nil {
		t.Errorf("%v", err)
	}
	if string(content) != "Hello\n" {
		t.Errorf("Unexpected content of file foo.txt: got %v want %v", string(content), "Hello\n")
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

	r := loadTestFile(t, "readsyscall")
	dir, err := Compile(r, false)
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Make sure foo gets created in dir, so that the defer cleans it up.
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("foo.txt", []byte("Hello, world!"), 0755); err != nil {
		t.Fatal(err)
	}
	runWithArgs(t, "readsyscall", dir, nil, "Hello,", "")

	if err := ioutil.WriteFile("foo.txt", []byte("Goodbye"), 0755); err != nil {
		t.Fatal(err)
	}
	runWithArgs(t, "readsyscall", dir, nil, "Goodby", "")
}

// Argc prints the number of arguments passed.
func TestArgc(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	r := loadTestFile(t, "argc")
	dir, err := Compile(r, false)
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	runWithArgs(t, "argc", dir, nil, "1", "")
	runWithArgs(t, "argc", dir, []string{"other", "params"}, "3", "")
}

// ArgLens prints the length of each argument passed to ensure argc was converted from
// a C string correctly.
func TestArgLens(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	r := loadTestFile(t, "arglens")
	dir, err := Compile(r, false)
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	runWithArgs(t, "arglens", dir, nil, "\n", "")
	runWithArgs(t, "arglens", dir, []string{"other", "params"}, "5 6\n", "")
}

// Echo is one of the simplest useful program that takes arguments
func TestEchoProgram(t *testing.T) {
	// We chdir below, defer a cleanup that resets it after the test finishes.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pwd string) {
		os.Chdir(pwd)
	}(pwd)

	r := loadTestFile(t, "echo")
	dir, err := Compile(r, false)
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	if err != nil {
		t.Fatal(err)
	}

	runWithArgs(t, "echo", dir, []string{"foo", "bar"}, "foo bar\n", "")
	runWithArgs(t, "echo", dir, []string{"other", "params"}, "other params\n", "")
}

// TestCatPrograms tests various implementations of cat of varying complexity.
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
		// Compile needs to be in the default directory so that it can
		// find the test program
		if err := os.Chdir(pwd); err != nil {
			t.Fatal(err)
		}

		r := loadTestFile(t, c)
		dir, err := Compile(r, false)
		if dir != "" {
			defer os.RemoveAll(dir)
		}
		if err != nil {
			t.Fatal(err)
		}
		// Make sure foo gets created in dir, so that the defer cleans it up, now
		// that we've parsed the program.
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}

		if err := ioutil.WriteFile("foo.tmp", []byte("Foo"), 0666); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile("bar.tmp", []byte("Bar"), 0666); err != nil {
			t.Fatal(err)
		}
		runWithArgs(t, c, dir, []string{"foo.tmp", "bar.tmp"}, "FooBar", "")
		runWithArgs(t, c, dir, []string{"bar.tmp", "foo.tmp"}, "BarFoo", "")
	}
}
