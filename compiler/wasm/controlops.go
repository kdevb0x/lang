package wasm

import (
	"fmt"
)

type Loop struct{}

func (i Loop) TextFormat(ctx Context) string {
	return "loop"
}

func (i Loop) String() string {
	return "loop"
}

type Block struct{}

func (i Block) TextFormat(ctx Context) string {
	return "block"
}

func (i Block) String() string {
	return "block"
}

type If struct{}

func (i If) TextFormat(ctx Context) string {
	return "if"
}

func (i If) String() string {
	return "if"
}

type End struct{}

func (i End) TextFormat(ctx Context) string {
	return "end"
}

func (i End) String() string {
	return "end"
}

type Else struct{}

func (i Else) TextFormat(ctx Context) string {
	return "else"
}

func (i Else) String() string {
	return "else"
}

type Br uint

func (i Br) TextFormat(ctx Context) string {
	return i.String()
}

func (i Br) String() string {
	return fmt.Sprintf("br %d", i)
}

type BrIf uint

func (i BrIf) TextFormat(ctx Context) string {
	return i.String()
}

func (i BrIf) String() string {
	return fmt.Sprintf("br_if %d", i)
}
