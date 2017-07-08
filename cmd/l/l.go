package main

import (
	//	"fmt"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/driusan/lang/compiler/codegen"
	"github.com/driusan/lang/compiler/irgen"
	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/token"
)

func main() {
	// For now, jut assume the command is building a program in the
	// current directory.
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// Combine all the .l files into a MultiReader for BuildProgram
	var srcFiles []io.Reader
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".l" {
			continue
		}
		fi, err := os.Open(f.Name())
		if err != nil {
			log.Println(err)
			return
		}
		defer fi.Close()

		srcFiles = append(srcFiles, fi)
	}
	src := io.MultiReader(srcFiles...)

	// And build the program.
	if err := BuildProgram(src); err != nil {
		log.Fatal(err)
	}

}

func BuildProgram(src io.Reader) error {
	// FIXME: BuildProgram should probably be in some other package,
	// so that it can be used by both the compiler tests and the
	// command line client.
	d, err := ioutil.TempDir("", "langbuild")
	if err != nil {
		return err
	}
	//	defer os.RemoveAll(d)

	f, err := os.Create(d + "/main.s")
	if err != nil {
		return err
	}
	defer f.Close()

	// Tokenize needs a RuneReader, so wrap the reader around a bufio
	tokens, err := token.Tokenize(bufio.NewReader(src))
	if err != nil {
		return err
	}

	// Generate the AST
	prog, ti, err := ast.Construct(tokens)
	if err != nil {
		return err
	}

	// Generate the type information before the functions.
	enums := make(irgen.EnumMap)
	for _, v := range prog {
		switch v.(type) {
		case ast.SumTypeDefn:
			_, opts, err := irgen.GenerateIR(v, ti, enums)
			if err != nil {
				return err
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
			fnc, _, err := irgen.GenerateIR(v, ti, enums)
			if err != nil {
				return err
			}
			if err := codegen.Compile(f, fnc); err != nil {
				return err
			}
		case ast.TypeDefn:
			// No IR for types, we've already verified them.
		default:
			panic("Unhandled AST node type for code generation")
		}
	}

	// Assemble and link the file.
	// FIXME: Make this more robust, or at least move it to a helper. It
	// will only work on Plan 9 right now.
	// FIXME: The program name shouldn't be hardcoded as "main"
	cmd := exec.Command("6a", "-o", d+"/main.6", d+"/main.s")
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	cmd = exec.Command("6l", "-o", "./main", d+"/main.6")
	_, err = cmd.Output()
	return err
}
