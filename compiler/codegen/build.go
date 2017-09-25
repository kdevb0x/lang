package codegen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/driusan/lang/compiler/irgen"
	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/token"
)

// Builds a program. Directory d is used as the workspace, to build in,
// and the source code for the program comes from src.
//
// Returns the name of the executable created in d or an error
func BuildProgram(d string, src io.Reader) (string, error) {
	// FIXME: This should be a library, not hardcoded string consts.
	// FIXME: Make other architecture entrypoints..
	f, err := os.Create(d + "/_main.s")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(f, entrypoint+"\n")
	fmt.Fprintf(f, exits+"\n")
	fmt.Fprintf(f, write+"\n")
	fmt.Fprintf(f, read+"\n")
	fmt.Fprintf(f, open+"\n")
	fmt.Fprintf(f, closestr+"\n")
	fmt.Fprintf(f, createf+"\n", O_WRONLY|O_CREAT)
	fmt.Fprintf(f, printstring+"\n")
	fmt.Fprintf(f, printbyteslice+"\n")
	fmt.Fprintf(f, printint+"\n")
	f.Close()

	f, err = os.Create(d + "/main.s")
	if err != nil {
		return "", err
	}
	defer f.Close()

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
	enums := make(irgen.EnumMap)
	for _, v := range prog {
		switch v.(type) {
		case ast.SumTypeDefn:
			_, opts, err := irgen.GenerateIR(v, ti, c, enums)
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
		case ast.FuncDecl, ast.ProcDecl:
			fnc, _, err := irgen.GenerateIR(v, ti, c, enums)
			if err != nil {
				return "", err
			}
			if err := Compile(f, fnc); err != nil {
				return "", err
			}
		case ast.TypeDefn, ast.SumTypeDefn:
			// No IR for types, we've already verified them.
		default:
			panic("Unhandled AST node type for code generation")
		}
	}

	// FIXME: Make this more robust and/or not depend on the Go toolchain.
	cmd := exec.Command("go", "tool", "asm", "-o", d+"/main.o", d+"/main.s")
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}
	cmd = exec.Command("go", "tool", "asm", "-o", d+"/_main.o", d+"/_main.s")
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	cmd = exec.Command("go", "tool", "pack", "c", d+"/main.a", d+"/_main.o", d+"/main.o")
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	if p := os.Getenv("LPATH"); p == "" {
		// Avoid the Go runtime "main" symbol by building a fake runtime
		cmd := exec.Command("go", "build", "-o", d+"/runtime.a", "github.com/driusan/noruntime/runtime")
		_, err = cmd.Output()
		if err != nil {
			return "", err
		}

		cmd = exec.Command("go", "tool", "link", "-E", "_main", "-g", "-L", d, "-w", "-o", d+"/main", d+"/main.a")
		_, err = cmd.Output()
		if err != nil {
			return "", err
		}
	} else {
		// There should already be a fake runtime in LPATH/lib/
		cmd = exec.Command("go", "tool", "link", "-E", "_main", "-g", "-L", p+"/lib/", "-w", "-o", d+"/main", d+"/main.a")
		_, err = cmd.Output()
		if err != nil {
			return "", err
		}
	}
	return "main", nil
}
