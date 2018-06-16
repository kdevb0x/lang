package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/driusan/lang/compiler/codegen"
	"github.com/driusan/lang/compiler/hlir/vm"
)

var debug bool

func main() {
	flag.BoolVar(&debug, "debug", false, "do not delete temporary files and print extra information to stderr")
	flag.Parse()
	// For now, jut assume the command is building a program in the
	// current directory.
	if debug {
		log.Println("Reading files in current directory")
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	inctest := len(args) > 0 && args[0] == "test"

	// Combine all the .l files into a MultiReader for BuildProgram
	var srcFiles []io.Reader
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".l" {
			continue
		}
		if !inctest && strings.HasSuffix(f.Name(), "_test.l") {
			continue
		}

		if debug {
			log.Println("Using", f.Name())
		}

		fi, err := os.Open(f.Name())
		if err != nil {
			log.Println(err)
			return
		}
		defer fi.Close()

		srcFiles = append(srcFiles, fi)
	}
	if len(srcFiles) == 0 {
		fmt.Fprintln(os.Stderr, "No source files available in current directory.")
		os.Exit(1)
	}
	src := io.MultiReader(srcFiles...)

	if len(args) > 0 {
		switch args[0] {
		case "test":
			if err := getVMAndRunTests(src); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		default:
			// And build the program.
			if err := buildAndCopyProgram(src); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		// And build the program.
		if err := buildAndCopyProgram(src); err != nil {
			log.Fatal(err)
		}
	}
}

// Builds a program in /tmp and copies the result to the current directory.
func buildAndCopyProgram(src io.Reader) error {
	// FIXME: BuildProgram should probably be in some other package,
	// so that it can be used by both the compiler tests and the
	// command line client.
	d, err := ioutil.TempDir("", "langbuild")
	if err != nil {
		return err
	}
	if debug {
		log.Println("Using temporary directory", d, "(WARNING: will not automatically delete in debug mode)")
	}
	if !debug {
		defer os.RemoveAll(d)
	}
	exe, err := codegen.BuildProgram(d, src)
	if err != nil {
		return err
	}
	if exe == "" {
		return fmt.Errorf("No executable built.")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name := path.Base(cwd)
	if name == "." || name == "" || name == "/" {
		log.Fatal("Could not determine appropriate executable name.")
	}
	return copyFile(d+"/"+exe, "./"+name)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return nil
}

func getVMAndRunTests(src io.Reader) error {
	machine, err := vm.ParseFromReader(src)
	if err != nil {
		return err
	}
	var fail, run uint
	for fname := range machine.Funcs {
		if strings.HasPrefix(fname, "Test") {
			// Make a clone of the VM to ensure that there's
			// no interactions between tests
			m2 := machine.Clone()
			_, _, err := vm.RunWithSideEffects(fname, m2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "--- FAIL %v:\n", fname)
				fmt.Fprintf(os.Stderr, "\t%v\n", err)
				fail++
			}
			run++
		}
	}
	if run > 0 {
		if fail == 0 {
			return nil
		}
		return fmt.Errorf("FAIL")
	} else {
		return fmt.Errorf("No tests run.")
	}
}
