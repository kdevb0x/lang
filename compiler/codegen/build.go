package codegen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/driusan/lang/compiler/irgen"
	"github.com/driusan/lang/parser/ast"

	_ "github.com/driusan/noruntime/runtime"
)

func BuildProgram(name, p string) (exe, dir string, err error) {
	// FIXME: This should be a param so the remove can be deferred
	d, err := ioutil.TempDir("", "langtest"+name)
	if err != nil {
		return "", "", err
	}
	//	defer os.RemoveAll(d)

	// FIXME: This should be a library, not hardcoded.
	// FIXME: Make other architecture entrypoints..
	f, err := os.Create(d + "/_main.s")
	if err != nil {
		return "", d, err
	}
	defer f.Close()
	fmt.Fprintf(f, entrypoint+"\n")
	fmt.Fprintf(f, exits+"\n")
	fmt.Fprintf(f, printstring+"\n")
	fmt.Fprintf(f, printint+"\n")

	f, err = os.Create(d + "/main.s")
	if err != nil {
		return "", d, err
	}
	defer f.Close()

	prog, ti, err := ast.Parse(p)
	if err != nil {
		return "", d, err
	}

	enums := make(irgen.EnumMap)
	for _, v := range prog {
		switch v.(type) {
		case ast.SumTypeDefn:
			_, opts, err := irgen.GenerateIR(v, ti, enums)
			if err != nil {
				return "", d, err
			}
			for k, v := range opts {
				enums[k] = v
			}
		default:
			// Handled below
		}

	}

	for _, v := range prog {
		switch v.(type) {

		case ast.FuncDecl, ast.ProcDecl:
			fnc, _, err := irgen.GenerateIR(v, ti, enums)
			if err != nil {
				return "", d, err
			}
			if err := Compile(f, fnc); err != nil {
				return "", d, err
			}
		case ast.TypeDefn, ast.SumTypeDefn:
			// No IR for types, we've already verified them.
		default:
			panic("Unhandled AST node type for code generation")
		}
	}

	// FIXME: Make this more robust, or at least move it to a helper. It
	// will only work on Plan 9 right now.
	cmd := exec.Command("go", "tool", "asm", "-o", d+"/main.o", d+"/main.s")
	_, err = cmd.Output()
	if err != nil {
		return "", d, err
	}
	cmd = exec.Command("go", "tool", "asm", "-o", d+"/_main.o", d+"/_main.s")
	_, err = cmd.Output()
	if err != nil {
		return "", d, err
	}

	cmd = exec.Command("go", "tool", "pack", "c", d+"/main.a", d+"/_main.o", d+"/main.o")
	_, err = cmd.Output()
	if err != nil {
		return "", d, err
	}

	if p := os.Getenv("LPATH"); p == "" {
		// Avoid the Go runtime "main" symbol by building a fake runtime
		cmd := exec.Command("go", "build", "-o", d+"/runtime.a", "github.com/driusan/noruntime/runtime")
		_, err = cmd.Output()
		if err != nil {
			return "", d, err
		}

		cmd = exec.Command("go", "tool", "link", "-E", "_main", "-g", "-L", d, "-w", "-o", d+"/main", d+"/main.a")
		_, err = cmd.Output()
		if err != nil {
			return "", d, err
		}
	} else {
		// There should already be a fake runtime in LPATH/lib/
		cmd = exec.Command("go", "tool", "link", "-E", "_main", "-g", "-L", p+"/lib/", "-w", "-o", d+"/main", d+"/main.a")
		_, err = cmd.Output()
		if err != nil {
			return "", d, err
		}
	}
	return "main", d, err
}
