package codegen

import (
	"fmt"
	"io"
	"strings"

	"github.com/driusan/lang/compiler/ir"
)

// Compile takes an AST and writes the assembly that it compiles to to
// w.
func Compile(w io.Writer, f ir.Func) error {
	printPragmas(w, f)
	fmt.Fprintf(w, "TEXT %v(SB), $%d\n", f.Name, reserveStackSize(f))
	data := dataLiterals(w, f)
	cpu := Amd64{stringLiterals: data, numArgs: f.NumArgs}
	cpu.clearRegisterMapping()
	for i := range f.Body {
		// For debugging, add a comment with the IR serialization
		//fmt.Fprintf(w, "\t%s // %s", cpu.ConvertInstruction(i, f.Body), f.Body[i])
		fmt.Fprintf(w, "\t%s\n", cpu.ConvertInstruction(i, f.Body))
	}
	if len(f.Body) == 0 || f.Body[len(f.Body)-1] != (ir.RET{}) {
		fmt.Fprintf(w, "\tRET\n")
	}

	return nil
}

type PhysicalRegister string

var stringNum uint

func printPragmas(w io.Writer, f ir.Func) {
	for _, op := range f.Body {
		if op == (ir.CALL{"printf"}) {
			fmt.Fprintf(w, "#pragma lib \"libstdio.a\"\n")
			fmt.Fprintf(w, "#pragma lib \"libc.a\"\n\n")
			return
		}
	}
	fmt.Fprintf(w, "#pragma lib \"libc.a\"\n\n")
	return
}
func dataLiterals(w io.Writer, f ir.Func) map[ir.StringLiteral]PhysicalRegister {
	v := make(map[ir.StringLiteral]PhysicalRegister)
	for _, op := range f.Body {
		rs := op.Registers()
		for _, r := range rs {
			if s, ok := r.(ir.StringLiteral); ok {
				name := printDataLiteral(w, string(s))
				v[s] = name
			}
		}
	}
	return v
}

func printDataLiteral(w io.Writer, str string) PhysicalRegister {
	name := fmt.Sprintf(".string%d<>", stringNum)
	stringNum++
	str = strings.Replace(str, `\n`, "\n", -1)
	for i := 0; i < len(str); i += 8 {
		if i+8 > len(str) {
			padding := i + 8 - len(str)
			toPrint := strings.Replace(str[i:], "\n", `\n`, -1)
			fmt.Fprintf(w, `%vDATA %s+%d(SB)/8, $"%s`, "\t", name, i, toPrint)
			for j := 0; j < padding; j++ {
				fmt.Fprintf(w, `\z`)
			}
			fmt.Fprintf(w, "\"\n")
			fmt.Fprintf(w, "\tGLOBL %s+0(SB), $%d\n", name, len(str)+padding)
			return PhysicalRegister(name)
		}
		toPrint := strings.Replace(str[i:i+8], "\n", `\n`, -1)
		fmt.Fprintf(w, "\tDATA %s+%d(SB)/8, $\"%s\"\n", name, i, toPrint)
	}
	fmt.Fprintf(w, "\tGLOBL %s+0(SB), $%d\n", name, len(str))
	return PhysicalRegister(name)
}
func reserveStackSize(f ir.Func) uint {
	// FIXME: This should be MIN(0, (numArgs-1)*8) + (8*NumLocalVariables)
	// but ir.Func doesn't know NumVariables
	return 32
}
