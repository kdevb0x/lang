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
