package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/driusan/lang/compiler/codegen"
	"github.com/driusan/lang/compiler/irgen"
	"github.com/driusan/lang/parser/ast"
)

func parseFile(file string) ([]ast.Node, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ast.Parse(string(f))
}

func build() error {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	d, err := ioutil.TempDir("", "lbuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(d)
	var objfiles []string
	for _, f := range files {
		if n := f.Name(); strings.HasSuffix(n, ".l") {
			ast, err := parseFile(f.Name())
			if err != nil {
				return err
			}

			dst := strings.TrimSuffix(n, ".l") + ".s"
			f, err := os.Create(d + "/" + dst)
			if err != nil {
				return err
			}
			defer f.Close()

			for _, v := range ast {
				fnc, err := irgen.GenerateIR(v)
				if err != nil {
					log.Fatal(err)
				}
				if err := codegen.Compile(f, fnc); err != nil {
					return err
				}
			}

			odst := strings.TrimSuffix(n, ".l") + ".o"

			println(d+"/"+odst, dst)
			cmd := exec.Command("6a", "-o", d+"/"+odst, d+"/"+dst)
			_, err = cmd.Output()
			if err != nil {
				return err
			}
			objfiles = append(objfiles, d+"/"+odst)

		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	exename := filepath.Base(cwd)
	args := append([]string{"-o", exename}, objfiles...)
	cmd := exec.Command("6l", args...)
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	switch args[0] {
	case "build":
		if err := build(); err != nil {
			log.Fatal(err)
		}

	default:
		flag.Usage()
	}
}
