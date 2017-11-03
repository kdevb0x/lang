package wasm

import (
	"fmt"
)

type I64Const int64

func (i I64Const) TextFormat(ctx Context) string {
	return i.String()
}

func (i I64Const) String() string {
	return fmt.Sprintf("i64.const %d", i)
}

type I64Add struct{}

func (i I64Add) TextFormat(ctx Context) string {
	return "i64.add"
}

func (i I64Add) String() string {
	return "i64.add"
}

type I64Sub struct{}

func (i I64Sub) TextFormat(ctx Context) string {
	return "i64.sub"
}

func (i I64Sub) String() string {
	return "i64.sub"
}

type I64Mul struct{}

func (i I64Mul) TextFormat(ctx Context) string {
	return "i64.mul"
}

func (i I64Mul) String() string {
	return "i64.mul"
}

type I64Div_S struct{}

func (i I64Div_S) TextFormat(ctx Context) string {
	return "i64.div_s"
}

func (i I64Div_S) String() string {
	return "i64.div_s"
}

type I64Rem_S struct{}

func (i I64Rem_S) TextFormat(ctx Context) string {
	return "i64.rem_s"
}

func (i I64Rem_S) String() string {
	return "i64.rem_s"
}

type I64GT_S struct{}

func (i I64GT_S) TextFormat(ctx Context) string {
	return "i64.gt_s"
}

func (i I64GT_S) String() string {
	return "i64.gt_s"
}

type I64GE_S struct{}

func (i I64GE_S) TextFormat(ctx Context) string {
	return "i64.ge_u"
}

func (i I64GE_S) String() string {
	return "i64.ge_u"
}

type I64GE_U struct{}

func (i I64GE_U) TextFormat(ctx Context) string {
	return "i64.ge_u"
}

func (i I64GE_U) String() string {
	return "i64.ge_u"
}

type I64EQ struct{}

func (i I64EQ) TextFormat(ctx Context) string {
	return "i64.eq"
}

func (i I64EQ) String() string {
	return "i64.eq"
}

type I64NE struct{}

func (i I64NE) TextFormat(ctx Context) string {
	return "i64.ne"
}

func (i I64NE) String() string {
	return "i64.ne"
}

type I64LT_S struct{}

func (i I64LT_S) TextFormat(ctx Context) string {
	return "i64.lt_s"
}

func (i I64LT_S) String() string {
	return "i64.lt_s"
}

type I64LE_S struct{}

func (i I64LE_S) TextFormat(ctx Context) string {
	return "i64.le_s"
}

func (i I64LE_S) String() string {
	return "i64.le_s"
}

type I64EQZ struct{}

func (i I64EQZ) TextFormat(ctx Context) string {
	return "i64.eqz"
}

func (i I64EQZ) String() string {
	return "i64.eqz"
}
