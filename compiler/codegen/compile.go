package codegen

import (
	"fmt"
	"io"
	"strings"

	"github.com/driusan/lang/compiler/mlir"
)

var debug bool = false

// Compile takes an AST and writes the assembly that it compiles to to
// w.
func Compile(w io.Writer, f mlir.Func) error {
	fmt.Fprintf(w, "TEXT %v(SB), 4+16, $%v\n", f.Name, reserveStackSize(f))
	data := dataLiterals(w, f)
	cpu := Amd64{stringLiterals: data, numArgs: f.NumArgs, lvOffsets: make(map[uint]uint), sliceBase: make(map[mlir.Register]bool)}
	cpu.clearRegisterMapping()
	// calculate the offsets of every local value
	offset := uint(f.NumArgs * 8)
	for _, op := range f.Body {
		regs := op.Registers()
		for _, r := range regs {
			switch lv := r.(type) {
			case mlir.LocalValue:
				_, ok := cpu.lvOffsets[lv.Id]
				if !ok {
					cpu.lvOffsets[lv.Id] = offset
					offset += uint(lv.Size())
				}
			}
		}

	}
	for i := range f.Body {
		// For debugging, add a comment with the IR serialization
		if debug {
			fmt.Fprintf(w, "\t%s // %s", cpu.ConvertInstruction(i, f.Body), f.Body[i])
		} else {
			fmt.Fprintf(w, "\t%s\n", cpu.ConvertInstruction(i, f.Body))
		}
	}
	if len(f.Body) == 0 || f.Body[len(f.Body)-1] != (mlir.RET{}) {
		fmt.Fprintf(w, "\tRET\n")
	}

	return nil
}

type PhysicalRegister string

func (pr PhysicalRegister) IsRealRegister() bool {
	switch string(pr) {
	case "AX", "BX", "CX", "DX", "SI", "DI", "BP", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15", "SP":
		return true
	default:
		return false
	}
}

var stringNum uint

func dataLiterals(w io.Writer, f mlir.Func) map[mlir.StringLiteral]PhysicalRegister {
	v := make(map[mlir.StringLiteral]PhysicalRegister)
	for _, op := range f.Body {
		rs := op.Registers()
		for _, r := range rs {
			if s, ok := r.(mlir.StringLiteral); ok {
				name := printDataLiteral(w, string(s))
				v[s] = name
			}
		}
	}
	return v
}

func printDataLiteral(w io.Writer, str string) PhysicalRegister {
	name := fmt.Sprintf("string%d<>", stringNum)

	// strings have the format struct{len int64, cstr *byte}.
	// Add the length before doing anything..

	stringNum++
	str = strings.Replace(str, `\n`, "\n", -1)
	fmt.Fprintf(w, "\tDATA %s+0(SB)/8, $%d\n", name, len(str))

	// Ensure that the string is nil terminated, in case it escapes to a
	// C function.
	if last := str[len(str)-1]; last != 0 {
		str += "\000"
	}

	for i := 0; i < len(str); i += 8 {
		if i+8 > len(str) {
			padding := i + 8 - len(str)
			toPrint := strings.Replace(str[i:], "\n", `\n`, -1)
			toPrint = strings.Replace(toPrint, "\000", `\000`, -1)
			fmt.Fprintf(w, `%vDATA %s+%d(SB)/8, $"%s`, "\t", name, i+8, toPrint)
			for j := 0; j < padding; j++ {
				fmt.Fprintf(w, `\000`)
			}
			fmt.Fprintf(w, "\"\n")
			fmt.Fprintf(w, "\tGLOBL %s+0(SB), 8+16, $%d\n", name, len(str)+padding+8)
			return PhysicalRegister(name)
		}
		toPrint := strings.Replace(str[i:i+8], "\n", `\n`, -1)
		fmt.Fprintf(w, "\tDATA %s+%d(SB)/8, $\"%s\"\n", name, i+8, toPrint)
	}
	fmt.Fprintf(w, "\tGLOBL %s+0(SB), 8+16, $%d\n", name, len(str)+8)
	return PhysicalRegister(name)
}

func reserveStackSize(f mlir.Func) string {
	if f.NumLocals == 0 && f.NumArgs == 0 {
		return fmt.Sprintf("%v", (f.LargestFuncCall+1)*8)
	} else if f.NumLocals == 0 {
		return fmt.Sprintf("%v-%d", (f.LargestFuncCall+1)*8, (f.NumArgs * 8))
	} else if f.NumArgs == 0 {
		return fmt.Sprintf("%d", (f.NumLocals*8)+((f.LargestFuncCall+1)*8))
	}
	return fmt.Sprintf("%d-%d", (f.NumLocals*8)+(f.LargestFuncCall*8)+8, f.NumArgs*8)
}
