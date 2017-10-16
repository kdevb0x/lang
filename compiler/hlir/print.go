package hlir

import (
	"fmt"
)

func PrettyPrint(level uint, ops []Opcode) string {
	ret := ""
	for _, op := range ops {
		for i := uint(0); i < level; i++ {
			ret += fmt.Sprintf("\t")
		}
		switch op.(type) {
		case RET, CALL, MOV, ADD, SUB, DIV, MUL, MOD:
			ret += fmt.Sprintf("%v", op)
		default:
			panic("Unhandled op in PrettyPrint")
		}
		ret += "\n"
	}
	return ret
}
