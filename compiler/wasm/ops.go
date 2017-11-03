package wasm

import (
	"fmt"
)

type Instruction interface {
	TextFormat(c Context) string
	// This will come later..
	// BinaryFormat(c Context) []byte
}

type Call struct {
	FuncName string
}

func (c Call) TextFormat(ctx Context) string {
	return fmt.Sprintf("call $%v", c.FuncName)
}

func (c Call) String() string {
	return fmt.Sprintf("call $%v", c.FuncName)
}

type Return struct{}

func (r Return) TextFormat(ctx Context) string {
	return r.String()
}

func (r Return) String() string {
	return fmt.Sprintf("return")
}

type GetGlobal int

func (gg GetGlobal) TextFormat(ctx Context) string {
	// FIXME: This should look up the index in ctx
	return gg.String()
}

func (gg GetGlobal) String() string {
	return fmt.Sprintf("get_global %d", gg)
}

type GetLocal int

func (gg GetLocal) TextFormat(ctx Context) string {
	// FIXME: This should look up the index in ctx
	return gg.String()
}

func (gg GetLocal) String() string {
	return fmt.Sprintf("get_local %d", gg)
}

type SetLocal int

func (gg SetLocal) TextFormat(ctx Context) string {
	// FIXME: This should look up the index in ctx
	return gg.String()
}

func (gg SetLocal) String() string {
	return fmt.Sprintf("set_local %d", gg)
}

type Drop struct{}

func (d Drop) TextFormat(ctx Context) string {
	return d.String()
}

func (d Drop) String() string {
	return "drop"
}

type Unreachable struct{}

func (d Unreachable) TextFormat(ctx Context) string {
	return d.String()
}

func (d Unreachable) String() string {
	return "unreachable"
}
