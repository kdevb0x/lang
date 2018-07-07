# lang

This is just me playing around with creating a toy language, which
I really have no business doing.

The language is (crudely) specified in the LANG.md file. The
`parser/sampleprograms` directory has a number of simple programs
that act as tests, and aren't necessarily idiomatic programs, as
they're mostly intended to catch regressions in the compiler.

The compiler will currently only work for the AMD64 architecture,
and requires the Go toolchain to be installed (which isn't a high
barrier, because the compiler itself is written in Go.) It currently
supports Plan 9, DragonFlyBSD, Linux, and MacOS X.

The code itself isn't terribly well written, because my initial
plan was to quickly bootstrap a self-hosted compiler, but in
retrospect that was probably too ambitious, so is a long way off
and for now the reference compiler will remain in Go and needs some
refactoring.

The package `github.com/driusan/noruntime/runtime` is required in
order to link a binary that doesn't include the overhead of the Go
runtime.

To install, run: 

```
go get github.com/driusan/noruntime/... 
go get github.com/driusan/lang/...
```

Which will install a compiler named "l" into your `$GOBIN` directory.
To invoke it, run "l" with no arguments. It'll concatenate all files
with the `.l` extension in the current directory and compile them
into a binary named after the directory.

The compiler is very buggy. If (when) you encounter any crashes,
or it compiles something that should be valid as per the language
spec but crashes, please create a GitHub issue with a sample program.
Ideally, your bug report should be small/short enough that it can
be included as a regression in the `parser/sampleprograms` directory.
(For now, the priority is to get all valid programs to compile
before getting all invalid programs to be rejected.)

## Roadmap

### Pre 0.1.0 (Status: "Sort of works, nothing significant has been written")

- [ ] Write all non-syscall standard library functions in native code, not assembly
	- [x] PrintString
	- [x] PrintByteSlice
	- [ ] PrintInt (almost done, needs variable slice indexing)
	- [ ] convert from Go strings to real files
- [ ] Write (native) autoformatter
	- [ ] port token package
	- [ ] new parse tree (before ast) package (keeps comments, whitespace)
		- [ ] packages/imports

### Pre 0.2.0 (Status: "Sort of works, but writing a few packages might have shaken out the bugs")

- [ ] Port test sub-command to native code
	- [ ] ast package
	- [ ] hlir
	- [ ] hlir/vm

### Pre 0.3.0 (Status: "Probably mostly works")

- [ ] Port compiler (while still invoking go for linking)
	- [ ] mlir package
	- [ ] codegen package
		- [ ] fork/exec support

### Pre 0.4.0 (Status: bootstrapped!)

- [ ] Fully bootstrapped
	- [ ] Native code generation (don't depend on Go assembler/linker, output binary directly)
	- [ ] other?

### Pre 0.5.0 (may be split up) (Status: Useable?)

- [ ] Optimizations 
	- [ ] Arithmetric simplification for constants
	- [ ] Comparison simplification for constants
	- [ ] Constant propagation
	- [ ] Eliminate unused instructions
		- [ ] Unused Dst for MOV
		- [ ] Unused Dst for ADD, SUB, DIV, MUL, and MOD
		- [ ] Unused Dst for EQ, NEQ, GEQ, GT, LT, LTE
	- [ ] Eliminate unused blocks
		- [ ] if false { ..1 } else { ..2} => ..2	 
		- [ ] if true  { ..1 } else { ..2} => ..1
		- [ ] while false { ..1 } => eliminate
	- [ ] Compile time evaluation of pure functions with constant arguments
	- [ ] Inlining
		- [ ] Small functions
		- [ ] Pure functions with at least 1 constant argument

### Unscheduled/When needed

- [ ] heap variables
- [ ] interfaces/polymorphism
- [ ] "l test -compile" (run tests by compiling a binary, not running in a VM)
- [ ] "l test -static" (static analysis tests)
- [ ] "l test -fuzz [-compile] [funcname]" (run tests with random arguments that meet preconditions and ensure no assertion failures)
- [ ] Improve assertion support
	- [x] make assert work with compiled code, not just interpreted
	- [ ] include line number/location in assert error message (post 0.3.0)
	- [ ] include way to change behaviour of assertions at build time (ignore/warn/die) (post 0.3.0)
- [ ] Refactor the bootstrap compiler to be good code? (Or just live with it until 0.3.0?)
- [ ] Type-based function overloading 
