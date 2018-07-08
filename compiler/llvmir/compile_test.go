package llvmir

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/llir/llvm/ir"
)

func parseFile(t *testing.T, testcase string) *ir.Module {
	t.Helper()
	f, err := os.Open("../../testsuite/" + testcase + ".l")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m, err := Generate(f)
	if err != nil {
		t.Fatal(err)
	}

	return m
}

func startSymbol(m *ir.Module) string {
	for _, f := range m.Funcs {
		if f.GetName() == "main" {
			switch len(f.Params()) {
			case 0:
				return `declare void @main()

	define void @_start() {
		call void @main()
		call void asm sideeffect "movq $$0, %rdi\0Amovq $$` + SYS_EXIT + `, %rax\0Asyscall\0A", ""()

		unreachable
	}`
			case 1:
				return `declare void @main({ { i8*, i64}*, i64} %args )

	define i64 @cstrlen(i8* %str) {
		%i = alloca i64
		store i64 0, i64* %i
		br label %loop
		loop:
		%ival = load i64, i64* %i
		%1 = getelementptr i8, i8* %str, i64 %ival
		%chr = load i8, i8* %1
		%zero = icmp eq i8 0, %chr
		br i1 %zero, label %ret, label %inc
		inc:
		%2 = add i64 1, %ival
		store i64 %2, i64* %i
		br label %loop
		ret:

		ret i64 %ival
	}

	define void @_start() {
		%argc = call i64 asm sideeffect "movq 0(%rdi), %rax", "=A"()
		%argv = call i8** asm sideeffect "movq %rdi, %rax\0Aaddq $$8, %rax", "=A"()

		%i = alloca i64
		store i64 %argc, i64* %i
		br label %convarg
		convarg:
		%ival = load i64, i64* %i
		%dec = sub i64 %ival, 1
		%1 = alloca {i8*, i64}
		%2 = load {i8*, i64}, {i8*, i64}* %1
		%3 = getelementptr i8*, i8** %argv, i64 %dec
		%4 = load i8*, i8** %3
		%5 = insertvalue {i8*, i64} %2, i8* %4, 0
		%6 = call i64 @cstrlen(i8* %4)
		%7 = insertvalue {i8*, i64} %5, i64 %6, 1
		store { i8*, i64 } %7, {i8*, i64}* %1

		store i64 %dec, i64* %i
		%zero = icmp eq i64 0, %dec
		br i1 %zero, label %run, label %convarg
		run:
		%withargv =insertvalue { { i8*, i64 }* , i64 } { {i8*, i64}* undef, i64 0}, {i8*, i64}* %1, 0
		%args = insertvalue { { i8*, i64 }* , i64 } %withargv, i64 %argc, 1
		call void @main({ { i8*, i64 }* , i64 } %args)

		; inline exit syscall
		call void asm sideeffect "movq $$0, %rdi\0Amovq $$` + SYS_EXIT + `, %rax\0Asyscall\0A", ""()
		unreachable
	}`
			default:
				panic("Main must be either empty or a string slice")
			}
		}
	}
	// No main symbol, must be a library.
	return ""
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

// Compile a program and return the dir that was used as a temporary dir.
// It's the caller's responsibility to clean up dir.
func compile(t *testing.T, name string) string {
	m := parseFile(t, name)
	dir, err := ioutil.TempDir("", "langtest"+name+"_")
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(dir + "/usercode.ll")
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := os.Create(dir + "/runtime.ll")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Fprintf(f, "%v", m)

	// FIXME: The builtins that have been re-writing in stdlib should be used instead of these.
	fmt.Fprintf(runtime, "%v", startSymbol(m)+`

	define void @Exit(i64 %code) {
		call void asm sideeffect "movq $$`+SYS_EXIT+`, %rax\0Asyscall\0A", "{di}"(i64 %code)
		unreachable
	}
	define void @PrintString({i8*, i64} %str) {
		tail call void @Write(i64 1, {i8*, i64} %str)
		ret void
	}

	define void @Write(i64 %fd, {i8*, i64} %buf) {
		%base = extractvalue {i8*, i64} %buf, 0
		%size = extractvalue {i8*, i64} %buf, 1
		tail call void asm sideeffect "movq $$`+SYS_WRITE+`, %rax\0Asyscall\0A", "{di},{dx},{si},~{dirflag},~{fpsr},~{flags}"(i64 %fd, i64 %size, i8* %base)
		ret void
	}

	define i64 @Read(i64 %fd, {i8*, i64} %buf) {
		%base = extractvalue {i8*, i64} %buf, 0
		%size = extractvalue {i8*, i64} %buf, 1
		%1 = call i32 asm sideeffect "movq $$`+SYS_READ+`, %rax\0Asyscall\0A", "=A,{di},{dx},{si},~{dirflag},~{fpsr},~{flags}"(i64 %fd, i64 %size, i8* %base)
		%2 = sext i32 %1 to i64
		ret i64 %2
	}


	define i64 @Create({i8*, i64} %buf) {
		%base = extractvalue {i8*, i64} %buf, 0
		%size = extractvalue {i8*, i64} %buf, 1
		%1 = call i32 asm sideeffect "`+strings.Replace(`movq %rcx, %rax
		movq %rsp,%rdi
		subq %rcx, %rdi
		cld
		rep movsb
		subq %rax, %rdi
movq $$438, %rdx
movq $$`+CREATE_CONST+`, %rsi
movq $$`+SYS_OPEN+`, %rax
syscall`, "\n", `\0A`, -1)+`", "=A,{si},{cx}"(i8* %base, i64 %size) 
		%2 = sext i32 %1 to i64
		ret i64 %2
	}

	define i64 @Open ({i8*, i64} %buf) {
		%base = extractvalue {i8*, i64} %buf, 0
		%size = extractvalue {i8*, i64} %buf, 1
		%1 = call i32 asm sideeffect "`+strings.Replace(`movq %rcx, %rax
		movq %rsp,%rdi
		subq %rcx, %rdi
		subq $$16, %rdi`+ // FIXME: This should be aligned, this is just so that movsb doesn't
		// overwrite the return address on the stack.
		`
		cld
		rep movsb
		subq %rax, %rdi
		movb $$0, (%rdi, %rax)
movq $$0, %rsi
movq $$0, %rdx
movq $$`+SYS_OPEN+`, %rax
syscall`, "\n", `\0A`, -1)+`", "=A,{si},{cx},~{dirflag},~{fpsr},~{flags}"(i8* %base, i64 %size)
		%2 = sext i32 %1 to i64
		ret i64 %2
	}


	define void @Close(i64 %fd) {
		tail call void asm sideeffect "movq $$`+SYS_CLOSE+`, %rax\0Asyscall\0A", "{di},~{dirflag},~{fpsr},~{flags}"(i64 %fd)
		ret void
	}
	
	define void @PrintInt(i64 %n) {
		%1 = alloca [20 x i8]
		%reversed = alloca [20 x i8]
		%rem = alloca i64
		%len = alloca i8
		%i = alloca i8
		store i64 %n, i64* %rem
		store i8 0, i8* %len
		store i8 0, i8* %i

		%zero = icmp eq i64 0, %n
		br i1 %zero, label %printzero, label %printnonzero
		printzero:
		%rfirst = getelementptr [20 x i8], [20 x i8]* %reversed, i8 0, i8 0
		store i8 48, i8* %rfirst
		store i8 1, i8* %len
		br label %print
		printnonzero:
		%isneg = icmp slt i64 %n, 0
		br i1 %isneg, label %neg, label %pos
		neg:
		%abs = mul i64 -1, %n
		store i64 %abs, i64* %rem
		br label %adddigit
		pos:
		store i64 %n, i64* %rem
		br label %adddigit
		adddigit:
		%idx = load i8, i8* %len
		%el = getelementptr [20 x i8], [20 x i8]* %1, i8 0, i8 %idx

		%x = load i64, i64* %rem
		%digit = srem i64 %x, 10
		%trunc = trunc i64 %digit to i8
		%dchar = add i8 48, %trunc
		store i8 %dchar, i8* %el

		%newlen = add i8 1, %idx
		store i8 %newlen, i8* %len

		%newval = sdiv i64 %x, 10
		store i64 %newval, i64* %rem

		%done = icmp eq i64 0, %newval
		br i1 %done, label %checksign, label %adddigit
		checksign:
		br i1 %isneg, label %addsign, label %reverse
		addsign:
		%newlen2 = add i8 1, %idx
		store i8 %newlen2, i8* %len
		%sel = getelementptr [20 x i8], [20 x i8]* %1, i8 0, i8 %newlen2
		store i8 45, i8* %sel
		%newlen3 = add i8 1, %newlen2
		store i8 %newlen3, i8* %len
		br label %reverse
		reverse:
		%finallen = load i8, i8* %len
		%iv = load i8, i8* %i

		%ridx = getelementptr [20 x i8], [20 x i8]* %reversed, i8 0, i8 %iv
		%2 = sub i8 %finallen, %iv
		%negidx = sub i8 %2, 1

		%rdigit = getelementptr [20 x i8], [20 x i8]* %1, i8 0, i8 %negidx
		%dgt = load i8, i8* %rdigit

		store i8 %dgt, i8* %ridx

		%addi = add i8 1, %iv
		store i8 %addi, i8* %i

		%donerev = icmp eq i8 %addi, %finallen
		br i1 %donerev, label %print, label %reverse
		print:
		%lenval = load i8, i8* %len
		%zexlen = zext i8 %lenval to i64
		%first = getelementptr [20 x i8], [20 x i8]* %reversed, i8 0, i8 0
		%ptr = insertvalue {i8*, i64} {i8* undef, i64 undef}, i8* %first, 0
		%ptrwlen = insertvalue {i8*, i64} %ptr, i64 %zexlen, 1
		call void @Write(i64 1, {i8*, i64} %ptrwlen)
		ret void
	}


	`)

	cmd := exec.Command("llc", "-filetype=obj", dir+"/runtime.ll")
	if err := cmd.Start(); err != nil {
		t.Errorf("%v: llc: %v", name, err)
		return dir
	}
	if err := cmd.Wait(); err != nil {
		t.Errorf("%v: llc: %v", name, err)
		return dir
	}
	cmd = exec.Command("llc", "-filetype=obj", dir+"/usercode.ll")
	if err := cmd.Start(); err != nil {
		t.Errorf("%v: llc: %v", name, err)
		return dir
	}
	if err := cmd.Wait(); err != nil {
		t.Errorf("%v: llc: %v", name, err)
		return dir
	}
	cmd = exec.Command("llc", "-filetype=obj", dir+"/usercode.ll")
	if err := cmd.Start(); err != nil {
		t.Errorf("%v: llc: %v", name, err)
		return dir
	}
	if err := cmd.Wait(); err != nil {
		t.Errorf("%v: llc: %v", name, err)
		return dir
	}
	cmd = exec.Command("ld", "-o", dir+"/main", dir+"/usercode.o", dir+"/runtime.o")
	if err := cmd.Start(); err != nil {
		t.Errorf("%v: ld: %v", name, err)
		return dir
	}
	if err := cmd.Wait(); err != nil {
		t.Errorf("%v: ld: %v", name, err)
		return dir
	}
	return dir
}

func compileAndRun(t *testing.T, name string, estdout, estderr string) {
	t.Helper()
	dir := compile(t, name)
	defer os.RemoveAll(dir)
	runWithArgs(t, name, dir, nil, estdout, estderr)
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
		{"swap", "", ""},
		{"reverse", "", ""},
		{"digitsinto", "", ""},
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

	dir := compile(t, "createsyscall")
	defer os.RemoveAll(dir)

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

	dir := compile(t, "readsyscall")
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

	dir := compile(t, "argc")
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

	dir := compile(t, "arglens")
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

	dir := compile(t, "echo")
	defer os.RemoveAll(dir)

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

		dir := compile(t, c)
		defer os.RemoveAll(dir)
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
