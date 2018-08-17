package gobackend

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/driusan/lang/parser/ast"
	//"github.com/driusan/lang/parser/token"
)

func compileAndRun(t *testing.T, name string, nodes []ast.Node, stdout, stderr string) {
	t.Helper()

	dirname, err := ioutil.TempDir("", "langgobackendtest"+name)
	if err != nil {
		t.Fatal(err)
	}
	//defer os.RemoveAll(dirname)

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(pwd)

	if err := os.Chdir(dirname); err != nil {
		t.Fatal(err)
	}

	partialf, err := os.Create("main.partial")
	if err != nil {
		t.Fatal(err)
	}
	defer partialf.Close()

	imports := make(map[string]bool)
	for _, n := range nodes {
		newimports, err := Convert(partialf, n)
		if err != nil {
			t.Errorf("%v: %v", name, err)
			return
		}
		for imp := range newimports {
			imports[imp] = true
		}
	}

	f, err := os.Create("main.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "package main\n")
	for imp := range imports {
		fmt.Fprintf(f, "import \"%v\"\n", imp)
	}

	partialf.Seek(0, io.SeekStart)
	if _, err := io.Copy(f, partialf); err != nil {
		t.Errorf("%v: %v", name, err)
		return
	}

	buildAndRun(t, name, stdout, stderr)
}

func buildAndRun(t *testing.T, name string, estdout, estderr string) {
	t.Helper()

	build := exec.Command("go", "build", "-o", "main")
	if err := build.Run(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr strings.Builder
	ex := exec.Command("./main")
	ex.Stdout = &stdout
	ex.Stderr = &stderr

	if err := ex.Run(); err != nil {
		t.Fatalf("%v: %v", name, err)
	}
	if stdout.String() != estdout {
		t.Errorf("%v: unexpected stdout got %v want %v", name, stdout.String(), estdout)
	}
	if stderr.String() != estderr {
		t.Errorf("%v: unexpected stderr got %v want %v", name, stderr.String(), estderr)
	}
}
func TestTestSuite(t *testing.T) {
	tests := []struct {
		filename, estdout, estderr string
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

		{"enumarray", "", ""},
		{"enumarrayexplicit", "", ""},
		{"sumtypearray", "", ""},
		{"sumtypearraymatch", "stringy33", ""},
		// needs "default" keyword and variable slice indexing
		//{"printint", "67\n-568", ""},
	}

	for _, tc := range tests {
		as := parseTestCase(t, tc.filename)

		compileAndRun(t, tc.filename, as, tc.estdout, tc.estderr)
	}
}
