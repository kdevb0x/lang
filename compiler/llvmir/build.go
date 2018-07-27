package llvmir

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/llir/llvm/ir"
)

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

// Compile a program and return the dir that was used as a temporary dir.
// It's the caller's responsibility to clean up dir.
func Compile(src io.Reader, noopt bool) (string, error) {
	dir, err := ioutil.TempDir("", "langbuild_")
	if err != nil {
		return dir, err
	}
	f, err := os.Create(dir + "/usercode.ll")
	if err != nil {
		return dir, err
	}
	defer f.Close()

	runtime, err := os.Create(dir + "/runtime.ll")
	if err != nil {
		return dir, err
	}
	defer runtime.Close()

	m, err := Generate(src)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(f, "%v", m)

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
		subq $$16, %rdi`+ // FIXME: This should be aligned, this is just so that movsb doesn't
		// overwrite the return address on the stack.
		`
		cld
		rep movsb
		subq %rax, %rdi
		movb $$0, (%rdi, %rax)
movq $$438, %rdx
movq $$`+CREATE_CONST+`, %rsi
movq $$`+SYS_OPEN+`, %rax
syscall`, "\n", `\0A`, -1)+`", "=A,{si},{cx},~{dirflag},~{fpsr},~{flags}"(i8* %base, i64 %size)
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

	cmd := exec.Command("llvm-link", "-o="+dir+"/code.bc", dir+"/runtime.ll", dir+"/usercode.ll")
	if err := cmd.Start(); err != nil {
		return dir, err
	}
	if err := cmd.Wait(); err != nil {
		return dir, err
	}

	if !noopt {
		// This is the ideal, but some of the options cause the test
		// suite to crash, so we need to figure out why before enabling
		// them.
		cmd = exec.Command("opt", "-o", dir+"/code.bc", "-internalize", "-internalize-public-api-list=_start", "-inline", "-ipconstprop", "-globaldce", "-constprop", "-sccp", "-dce", "-sroa", "-jump-threading", "-mem2reg", dir+"/code.bc")
		// This is what is run:
		// Missing: -inline, -mem2reg
		cmd = exec.Command("opt", "-o", dir+"/code.bc", "-internalize", "-internalize-public-api-list=_start", "-ipconstprop", "-globaldce", "-constprop", "-sccp", "-dce", "-sroa", "-jump-threading", dir+"/code.bc")
		if err := cmd.Start(); err != nil {
			return dir, err
		}
		if err := cmd.Wait(); err != nil {
			return dir, err
		}
	}

	cmd = exec.Command("llc", "-filetype=obj", dir+"/code.bc")
	if err := cmd.Start(); err != nil {
		return dir, err
	}
	if err := cmd.Wait(); err != nil {
		return dir, err
	}
	cmd = exec.Command("ld", "-o", dir+"/main", dir+"/code.o")
	if err := cmd.Start(); err != nil {
		return dir, err
	}
	if err := cmd.Wait(); err != nil {
		return dir, err
	}
	return dir, nil
}
