package wasm

import (
	"fmt"
)

type VarType byte

const (
	i32 = VarType(iota)
	i64
	f32
	f64
)

func (v VarType) String() string {
	switch v {
	case i32:
		return "i32"
	case i64:
		return "i64"
	case f32:
		return "f32"
	case f64:
		return "f64"
	default:
		panic("Unknown VarType")
	}
}

type VarKind byte

const (
	Param = VarKind(iota)
	Result
	Local
)

func (v VarKind) String() string {
	switch v {
	case Param:
		return "param"
	case Result:
		return "result"
	case Local:
		return "local"
	default:
		panic("Unknown type of parameter")
	}
}

type Variable struct {
	VarType
	VarKind
	Name string
}

func (v Variable) TextFormat(c Context) string {
	if v.Name == "" {
		return fmt.Sprintf("(%v %v)", v.VarKind, v.VarType)
	}
	return fmt.Sprintf("(%v $%v %v)", v.VarKind, v.Name, v.VarType)

}

type Signature []Variable

func (s Signature) TextFormat(c Context) string {
	var ret string

	for _, v := range s {
		ret += v.TextFormat(c) + " "
	}
	return ret
}

type Func struct {
	Name      string
	Signature Signature
	Body      []Instruction
}

func (f Func) TextFormat(c Context) string {
	ret := "(func"
	switch f.Name {
	case "":
	case "PrintInt", "PrintString", "PrintByteSlice", "len":
		// For now, just export everything that wasn't automaticaly imported..
		ret += fmt.Sprintf(` $%v`, f.Name)
	default:
		ret += fmt.Sprintf(` $%v (export "%v")`, f.Name, f.Name)
	}
	if f.Signature != nil {
		ret += " "
		ret += fmt.Sprintf("%v", f.Signature.TextFormat(c))
	}
	ret += "\n"
	for _, op := range f.Body {
		ret += "\t" + op.TextFormat(c) + "\n"
	}
	ret += ")\n"
	return ret
}
