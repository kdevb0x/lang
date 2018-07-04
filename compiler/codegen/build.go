package codegen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/driusan/lang/stdlib"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/compiler/mlir"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/token"
)

// builds a function and appends the assembly to dst.
//
// Used for building helpers from package stdlib and appending them to _main.s
func buildFunc(dst io.Writer, src string) error {
	nodes, ti, c, err := ast.Parse(src)
	if err != nil {
		return err
	}
	ir, _, err := mlir.Generate(nodes[0], ti, c, nil)
	if err != nil {
		return err
	}
	return Compile(dst, ir)
}

// Builds a program. Directory d is used as the workspace, to build in,
// and the source code for the program comes from src.
//
// Returns the name of the executable created in d or an error
func BuildProgram(d string, src io.Reader) (string, error) {
	mlir.Debug = false
	// FIXME: This should be a library, not hardcoded string consts.
	// FIXME: Make other architecture entrypoints..
	stdf, err := os.Create(d + "/main.s")
	if err != nil {
		return "", err
	}
	defer stdf.Close()
	//fmt.Fprintf(stdf, entrypoint+"\n")
	fmt.Fprintf(stdf, exits+"\n")
	fmt.Fprintf(stdf, write+"\n")
	fmt.Fprintf(stdf, read+"\n")
	fmt.Fprintf(stdf, open+"\n")
	fmt.Fprintf(stdf, closestr+"\n")
	fmt.Fprintf(stdf, createf+"\n", O_WRONLY|O_CREAT)
	fmt.Fprintf(stdf, printint+"\n")
	fmt.Fprintf(stdf, slicelen+"\n")

	if err := buildFunc(stdf, stdlib.PrintByteSlice); err != nil {
		return "", err
	}
	if err := buildFunc(stdf, stdlib.PrintString); err != nil {
		return "", err
	}

	/*f, err := os.Create(d + "/main.s")
	if err != nil {
		return "", err
	}
	defer f.Close()
	*/
	f := stdf
	maingo, err := os.Create(d + "/main.go")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(maingo, `package main
import (
	"os"
	"unsafe"
)

type lSlice struct{
	size uint64
	base uintptr
}

type lString struct{
	size uint64
	base *byte
}

func lmain(args lSlice)

func main() {
	largs := make([]lString, len(os.Args))
	for i := range os.Args {
		largs[i] = lString{
			uint64(len([]byte(os.Args[i]))),
			&([]byte(os.Args[i])[0]),
		}
	}
	if len(largs) > 0 {
		lmain(lSlice{
			uint64(len(largs)),
			uintptr(unsafe.Pointer(&largs[0])),
		})
	} else {
		lmain(lSlice{0, 0})
	}
	os.Exit(0)
}
`)
	maingo.Close()

	// Tokenize needs a Runereader, so wrap the reader around a bufio
	tokens, err := token.Tokenize(bufio.NewReader(src))
	if err != nil {
		return "", err
	}
	prog, ti, c, err := ast.Construct(tokens)
	if err != nil {
		return "", err
	}

	// Identify required type information before code generation
	// for the functions.
	enums := make(hlir.EnumMap)
	for _, v := range prog {
		switch v.(type) {
		case ast.EnumTypeDefn:
			_, opts, err := mlir.Generate(v, ti, c, enums)
			if err != nil {
				return "", err
			}
			for k, v := range opts {
				enums[k] = v
			}
		default:
			// Handled below
		}

	}

	// Generate the IR for the functions.
	for _, v := range prog {
		switch v.(type) {
		case ast.FuncDecl:
			fnc, _, err := mlir.Generate(v, ti, c, enums)
			if err != nil {
				return "", err
			}
			if err := Compile(f, fnc); err != nil {
				return "", err
			}
		case ast.TypeDefn, ast.EnumTypeDefn:
			// No IR for types, we've already verified them.
		default:
			panic("Unhandled AST node type for code generation")
		}
	}

	defer func(d string, err error) {
		if err := os.Chdir(d); err != nil {
			panic(err)
		}
	}(os.Getwd())
	if err := os.Chdir(d); err != nil {
		panic(err)
	}
	// FIXME: Make this more robust and/or not depend on the Go toolchain.
	cmd := exec.Command("go", "build", "-o", d+"/main")
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}
	return "main", nil
}
