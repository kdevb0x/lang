package wasm

import (
	"fmt"
)

type Import struct {
	Namespace, Function string
	Func                Func
}

func (i Import) TextFormat(c Context) string {
	return fmt.Sprintf(`(import "%v" "%v" %v)`, i.Namespace, i.Function, i.Func.TextFormat(c))

}

type Memory struct {
	Name string
	Size int
}

type Global struct {
	Mutable      bool
	Type         VarType
	InitialValue int
}

type Module struct {
	Imports []Import
	Memory  Memory
	Data    []DataSection
	Globals []Global
	Funcs   []Func
}

func (m Module) String() string {
	ret := "(module\n"
	context := NewContext(nil)
	context.Data = m.Data
	context.Functions = m.Funcs
	context.Imports = m.Imports

	for _, impt := range m.Imports {
		ret += impt.TextFormat(context) + "\n"
	}
	if m.Memory.Size != 0 {
		// FIXME: Memory should declare TextFormat
		ret += fmt.Sprintf("(memory (export \"%v\") %v)\n", m.Memory.Name, m.Memory.Size)
	}
	for _, data := range m.Data {
		ret += fmt.Sprintf(`(data (i32.const %d) "%v")`, data.Length, data.Content)
		ret += "\n"
	}
	for _, global := range m.Globals {
		ret += fmt.Sprintf("(global %v (i32.const %v))\n", global.Type, global.InitialValue)
	}

	for _, fnc := range m.Funcs {
		ret += fnc.TextFormat(context)
	}

	ret += ")"
	return ret

}
