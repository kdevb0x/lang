package wasm

import (
	"fmt"
)

type I32Const int32

func (i I32Const) TextFormat(ctx Context) string {
	return i.String()
}

func (i I32Const) String() string {
	return fmt.Sprintf("i32.const %d", i)
}

type I32Add struct{}

func (i I32Add) TextFormat(ctx Context) string {
	return "i32.add"
}

func (i I32Add) String() string {
	return "i32.add"
}

type I32Sub struct{}

func (i I32Sub) TextFormat(ctx Context) string {
	return "i32.sub"
}

func (i I32Sub) String() string {
	return "i32.sub"
}

type I32Mul struct{}

func (i I32Mul) TextFormat(ctx Context) string {
	return "i32.mul"
}

func (i I32Mul) String() string {
	return "i32.mul"
}

type I32Div_S struct{}

func (i I32Div_S) TextFormat(ctx Context) string {
	return "i32.div_s"
}

func (i I32Div_S) String() string {
	return "i32.div_s"
}

type I32Rem_S struct{}

func (i I32Rem_S) TextFormat(ctx Context) string {
	return "i32.rem_s"
}

func (i I32Rem_S) String() string {
	return "i32.rem_s"
}

type I32GT_S struct{}

func (i I32GT_S) TextFormat(ctx Context) string {
	return "i32.gt_s"
}

func (i I32GT_S) String() string {
	return "i32.gt_s"
}

type I32GE_S struct{}

func (i I32GE_S) TextFormat(ctx Context) string {
	return "i32.ge_s"
}

func (i I32GE_S) String() string {
	return "i32.ge_s"
}

type I32GE_U struct{}

func (i I32GE_U) TextFormat(ctx Context) string {
	return "i32.ge_u"
}

func (i I32GE_U) String() string {
	return "i32.ge_u"
}

type I32EQ struct{}

func (i I32EQ) TextFormat(ctx Context) string {
	return "i32.eq"
}

func (i I32EQ) String() string {
	return "i32.eq"
}

type I32NE struct{}

func (i I32NE) TextFormat(ctx Context) string {
	return "i32.ne"
}

func (i I32NE) String() string {
	return "i32.ne"
}

type I32LT_S struct{}

func (i I32LT_S) TextFormat(ctx Context) string {
	return "i32.lt_s"
}

func (i I32LT_S) String() string {
	return "i32.lt_s"
}

type I32LE_S struct{}

func (i I32LE_S) TextFormat(ctx Context) string {
	return "i32.le_s"
}

func (i I32LE_S) String() string {
	return "i32.le_s"
}

type I32EQZ struct{}

func (i I32EQZ) TextFormat(ctx Context) string {
	return "i32.eqz"
}

func (i I32EQZ) String() string {
	return "i32.eqz"
}

type I32WrapI64 struct{}

func (i I32WrapI64) TextFormat(ctx Context) string {
	return "i32.wrap/i64"
}

func (i I32WrapI64) String() string {
	return "i32.wrap/i64"
}

type I32Store struct{}

func (i I32Store) TextFormat(ctx Context) string {
	return i.String()
}

func (i I32Store) String() string {
	return "i32.store"
}

type I32Store8 struct{}

func (i I32Store8) TextFormat(ctx Context) string {
	return i.String()
}

func (i I32Store8) String() string {
	return "i32.store8"
}

type I32Load struct{}

func (i I32Load) TextFormat(ctx Context) string {
	return i.String()
}

func (i I32Load) String() string {
	return "i32.load"
}
