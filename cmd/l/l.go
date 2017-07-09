package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/driusan/lang/compiler/codegen"
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
	if err := buildAndCopyProgram(src); err != nil {
		log.Fatal(err)
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
	defer os.RemoveAll(d)

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
