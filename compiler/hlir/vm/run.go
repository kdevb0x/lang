package vm

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/ast"
)

var debug = false

func Write(fd int, nbytes int, str *byte)

// Run the HLIR function in a virtual machine, allowing all side-effects. This is primarily
// used for testing or scripting.
func RunWithSideEffects(f string, vm *Context) (stdout, stderr io.Reader, err error) {
	return run(vm.Funcs[f], vm, nil)
}

// Run the HLIR function in a virtual machine, but only allow the allowed side-effects. This is
// primarily used for compile time evaluation.
func RunWithLimitedEffects(f string, vm *Context, allowed []ast.Effect) (stdout, stderr io.Reader, err error) {
	return run(vm.Funcs[f], vm, allowed)
}

func run(f hlir.Func, ctx *Context, allowedEffects []ast.Effect) (stdout, stderr io.Reader, err error) {
	for _, op := range f.Body {
		stop, err := runOp(op, ctx, allowedEffects)
		if err != nil {
			return strings.NewReader(ctx.stdout.String()), strings.NewReader(ctx.stderr.String()), err
		}
		if stop {
			break
		}
	}
	return strings.NewReader(ctx.stdout.String()), strings.NewReader(ctx.stderr.String()), nil
}

func runOp(op hlir.Opcode, ctx *Context, allowed []ast.Effect) (stop bool, err error) {
	switch o := op.(type) {
	case hlir.CALL:
		newctx := ctx.Clone()
		newctx.localValues = make(map[hlir.LocalValue]interface{})
		newctx.funcRetVal = make(map[hlir.FuncRetVal]interface{})
		newctx.lastFuncCallRetVal = make(map[hlir.LastFuncCallRetVal]interface{})
		newctx.tempValue = make(map[hlir.TempValue]interface{})
		newctx.funcArg = make(map[hlir.FuncArg]interface{})

		rd := ctx.RegisterData[string(o.FName)]

		switch string(o.FName) {
		case "Write":
			fd := evalRegister(o.Args[0], ctx)
			l := evalRegister(o.Args[1], ctx)
			s := evalRegister(o.Args[2], ctx)

			switch fd.(int) {
			case 1:
				if s2, ok := s.(string); ok {
					fmt.Fprintf(ctx.stdout, "%s", s2)
				} else {
					base, nctx := dereferencePointer(o.Args[2], ctx)
				outer:
					for {
						switch o := base.(type) {
						case hlir.LocalValue:
							break outer
						case hlir.Offset:
							base, nctx = resolveOffset(o, nctx)
							break outer
						case hlir.SliceBasePointer:
							break outer
						default:
							base, nctx = dereferencePointer(base, nctx)
						}

					}
					ptr, ok := base.(hlir.SliceBasePointer)
					if ok {
						// It might not be a SliceBasePointer if it's
						// a string
						if o, ok := ptr.Register.(hlir.Offset); ok {
							base, nctx = resolveOffset(o, nctx)
						} else {
							base = ptr.Register
						}
					}
					for i := 0; i < l.(int); i++ {
						ch := evalRegister(base, nctx)
						if s, ok := ch.(string); ok {
							// Hack because Strings and Byte slices aren't represented
							// the same way in the VM, even though they should be
							fmt.Fprintf(ctx.stdout, "%s", s)
							break
						}
						fmt.Fprintf(ctx.stdout, "%c", evalRegister(base, nctx))
						switch b := base.(type) {
						case hlir.LocalValue:
							base = b + 1
						default:
							panic(fmt.Sprintf("Unhandled register type in Write of byte slice %v", reflect.TypeOf(b)))
						}
					}
				}
			case 2:
				if s2, ok := s.(string); ok {
					ctx.writeStderr(s2)
				} else {
					base, nctx := dereferencePointer(o.Args[2], ctx)

				outererr:
					for {
						switch o := base.(type) {
						case hlir.LocalValue:
							break outererr
						case hlir.Offset:
							base, nctx = resolveOffset(o, nctx)
							break outererr
						case hlir.SliceBasePointer:
							break outererr
						default:
							base, nctx = dereferencePointer(base, nctx)
						}

					}
					ptr, ok := base.(hlir.SliceBasePointer)
					if ok {
						// It might not be a SliceBasePointer if it's
						// a string
						if o, ok := ptr.Register.(hlir.Offset); ok {
							base, nctx = resolveOffset(o, nctx)
						} else {
							base = ptr.Register
						}
					}
					for i := 0; i < l.(int); i++ {
						ch := evalRegister(base, nctx)
						if s, ok := ch.(string); ok {
							// Hack because Strings and Byte slices aren't represented
							// the same way in the VM, even though they should be
							ctx.writeStderr(fmt.Sprintf("%s", s))
							break
						}
						ctx.writeStderr(fmt.Sprintf("%c", evalRegister(base, nctx)))
						switch b := base.(type) {
						case hlir.LocalValue:
							base = b + 1
						default:
							panic(fmt.Sprintf("Unhandled register type in Write of byte slice %v", reflect.TypeOf(b)))
						}
					}
				}
			default:
				if s2, ok := s.(string); ok {
					Write(fd.(int), l.(int), &([]byte(s2)[0]))
				} else {
					panic("Unhandled write of non-string")
				}
			}
		case "PrintString":
			// Special case for some things in stdlib.
			// FIXME: This should be more robust.
			if len(o.Args) == 1 {
				// It's a string literal
				fmt.Fprintf(ctx.stdout, "%v", evalRegister(o.Args[0], ctx))
			} else if len(o.Args) == 2 {
				// It's a len, localvalue pair
				l := evalRegister(o.Args[0], ctx)
				s := evalRegister(o.Args[1], ctx)
				if s2, ok := s.(string); ok {
					fmt.Fprintf(ctx.stdout, "%v", s2)
				} else {
					base, nctx := dereferencePointer(o.Args[1], ctx)

					for i := 0; i < l.(int); i++ {
						ch := evalRegister(base, nctx)
						if s, ok := ch.(string); ok {
							// Hack because Strings and Byte slices aren't represented
							// the same way in the VM, even though they should be
							fmt.Fprintf(ctx.stdout, "%s", s)
							break
						}
						if reg := evalRegister(base, nctx); reg != nil {
							fmt.Fprintf(ctx.stdout, "%c", reg)
						}
						switch b := base.(type) {
						case hlir.LocalValue:
							base = b + 1
						case hlir.FuncArg:
						default:
							panic("Unhandled register type in PrintByteSlice")
						}
					}
				}
			} else {
				panic("Unhandled PrintString")
			}
		case "PrintByteSlice":
			// FIXME: This should be way more robust and also more efficient.
			if len(o.Args) == 2 {
				// It's a len, localvalue pair
				l := evalRegister(o.Args[0], ctx)
				base, nctx := dereferencePointer(o.Args[1], ctx)

				for i := 0; i < l.(int); i++ {
					ch := evalRegister(base, nctx)
					if s, ok := ch.(string); ok {
						// Hack because Strings and Byte slices aren't represented
						// the same way in the VM, even though they should be
						fmt.Fprintf(ctx.stdout, "%s", s)
						break
					}
					fmt.Fprintf(ctx.stdout, "%c", evalRegister(base, nctx))
					switch b := base.(type) {
					case hlir.LocalValue:
						base = b + 1
					default:
						panic("Unhandled register type in PrintByteSlice")
					}
				}
			} else {
				panic("Unhandled byte slice")
			}
		case "len":
			ctx.lastFuncCallRetVal[hlir.LastFuncCallRetVal{0, 0}] = evalRegister(o.Args[0], ctx)
		case "PrintInt":
			fmt.Fprintf(ctx.stdout, "%v", evalRegister(o.Args[0], ctx))
		case "Create", "Open":
			var name string
			if len(o.Args) == 1 {
				// It's a string literal
				s := evalRegister(o.Args[0], ctx)
				name = s.(string)
			} else if len(o.Args) == 2 {
				// It's a len, localvalue pair
				l := evalRegister(o.Args[0], ctx)
				s := evalRegister(o.Args[1], ctx)
				if s2, ok := s.(string); ok {
					name = s2
				} else {
					base, nctx := dereferencePointer(o.Args[1], ctx)

					for i := 0; i < l.(int); i++ {
						ch := evalRegister(base, nctx)
						if s, ok := ch.(string); ok {
							// Hack because Strings and Byte slices aren't represented
							// the same way in the VM, even though they should be
							name = s
							break
						}
						name += fmt.Sprintf("%c", evalRegister(base, nctx))
						switch b := base.(type) {
						case hlir.LocalValue:
							base = b + 1
						default:
							panic("Unhandled register type in PrintByteSlice")
						}
					}
				}
			}
			var syscall func(string) (*os.File, error)
			if string(o.FName) == "Create" {
				syscall = os.Create
			} else if string(o.FName) == "Open" {
				syscall = os.Open
			}

			f, err := syscall(name)
			if err != nil {
				return true, err
			}
			ctx.lastFuncCallRetVal[hlir.LastFuncCallRetVal{0, 0}] = int(f.Fd())
		case "Read":
			fd := evalRegister(o.Args[0], ctx)
			l := evalRegister(o.Args[1], ctx)
			base, nctx := dereferencePointer(o.Args[2], ctx)
			f := os.NewFile(uintptr(fd.(int)), "unknown")
			bytes := make([]byte, l.(int), l.(int))
			n, err := f.Read(bytes)
			if err != nil && err != io.EOF {
				return true, err
			}
			for i := 0; i < l.(int); i++ {
				reg := base.(hlir.LocalValue) + hlir.LocalValue(i)
				nctx.SetRegister(reg, bytes[i])
			}
			ctx.lastFuncCallRetVal[hlir.LastFuncCallRetVal{0, 0}] = int(n)
		case "Close":
			fd := evalRegister(o.Args[0], ctx)
			f := os.NewFile(uintptr(fd.(int)), "unknown")
			f.Close()
		default:
			for i, r := range o.Args {
				farg := hlir.FuncArg{uint(i), false}
				if _, ok := rd[farg]; !ok {
					farg.Reference = true
					if _, okref := rd[farg]; !okref {
						panic(fmt.Sprintf("Function Argument %d does not have register data for %v", i, o.FName))
					}
				}
				if !farg.Reference {
					switch r := r.(type) {
					case hlir.Pointer:
						pointer := Pointer{r, ctx}
						newctx.pointers[hlir.Pointer{farg}] = pointer
					case hlir.SliceBasePointer:
						// Convert it to a normal pointer when passing
						// as an argument so that it dereferences properly.
						pointer := Pointer{hlir.Pointer{r.Register}, ctx}
						newctx.pointers[hlir.Pointer{farg}] = pointer
					default:
						newctx.funcArg[farg] = evalRegister(r, ctx)
					}
				} else {
					pointer := Pointer{r, ctx}
					newctx.pointers[hlir.Pointer{farg}] = pointer
				}
			}

			fnc, ok := ctx.Funcs[string(o.FName)]
			if !ok {
				return true, fmt.Errorf("Call to undefined function %s", o.FName)
			}
			if _, _, err := run(fnc, newctx, allowed); err != nil {
				return true, err
			}

			// Convert from funcRetVal in callee to lastFuncCallRetVal in caller
			//ctx.lastFuncCallRetVal = make(map[hlir.LastFuncCallRetVal]interface{})
			if o.TailCall {
				ctx.funcRetVal = newctx.funcRetVal
			} else {
				for k, v := range newctx.funcRetVal {
					ctx.lastFuncCallRetVal[hlir.LastFuncCallRetVal{0, uint(k)}] = v
				}
			}
		}
	case hlir.MOV:
		if fa, ok := o.Dst.(hlir.FuncArg); ok && fa.Reference {
			// If moving into a variable that was declared as a reference,
			// mutate in the source context
			ctxptr := ctx.pointers[hlir.Pointer{fa}]
			ptr := ctxptr.r.(hlir.Pointer)
			ctxptr.ctx.SetRegister(ptr.Register, evalRegister(o.Src, ctx))
		} else if ptr, ok := o.Src.(hlir.Pointer); ok {
			// If moving a pointer into a variable, make that variable a
			// pointer so that it can be used as a slice referece
			dptr := hlir.Pointer{o.Dst}

			switch ptr.Register.(type) {
			case hlir.Offset:
				deref, derefctx := resolveOffset(ptr.Register.(hlir.Offset), ctx)
				ctx.pointers[dptr] = Pointer{deref, derefctx}
			default:
				panic("Unhandled case for pointer MOV")
			}
		} else {
			ctx.SetRegister(o.Dst, evalRegister(o.Src, ctx))
		}
	case hlir.RET:
		return true, nil
	case hlir.ADD:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int)+b.(int)); err != nil {
			return true, err
		}
	case hlir.SUB:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int)-b.(int)); err != nil {
			return true, err
		}
	case hlir.MUL:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int)*b.(int)); err != nil {
			return true, err
		}
	case hlir.DIV:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int)/b.(int)); err != nil {
			return true, err
		}
	case hlir.MOD:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int)%b.(int)); err != nil {
			return true, err
		}
	case hlir.LT:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int) < b.(int)); err != nil {
			return true, err
		}
	case hlir.GT:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int) > b.(int)); err != nil {
			return true, err
		}
	case hlir.LTE:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int) <= b.(int)); err != nil {
			return true, err
		}
	case hlir.GEQ:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		// FIXME: Handle non-int
		if err := ctx.SetRegister(o.Dst, a.(int) >= b.(int)); err != nil {
			return true, err
		}
	case hlir.EQ:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		if err := ctx.SetRegister(o.Dst, a == b); err != nil {
			return true, err
		}
	case hlir.NEQ:
		a := evalRegister(o.Left, ctx)
		b := evalRegister(o.Right, ctx)
		if err := ctx.SetRegister(o.Dst, a != b); err != nil {
			return true, err
		}
	case hlir.LOOP:
		for _, in := range o.Initializer {
			stop, err := runOp(in, ctx, allowed)
			if err != nil || stop {
				return stop, err
			}
		}
		for evalCondition(o.Condition, ctx, allowed) {
			for _, in := range o.Body {
				if stop, err := runOp(in, ctx, allowed); err != nil || stop {
					return stop, err
				}
			}
		}
	case hlir.IF:
		for _, in := range o.Initializer {
			if stop, err := runOp(in, ctx, allowed); err != nil || stop {
				return stop, err
			}
		}
		if evalCondition(o.Condition, ctx, allowed) {
			for _, in := range o.Body {
				if stop, err := runOp(in, ctx, allowed); err != nil || stop {
					return stop, err
				}
			}
		} else {
			for _, in := range o.ElseBody {
				if stop, err := runOp(in, ctx, allowed); err != nil || stop {
					return stop, err
				}
			}
		}
	case hlir.JumpTable:
	jumptable:
		for _, cse := range o {
			for _, in := range cse.Initializer {
				if stop, err := runOp(in, ctx, allowed); err != nil || stop {
					return stop, err
				}
			}

			if evalCondition(cse.Condition, ctx, allowed) {
				for _, in := range cse.Body {
					if stop, err := runOp(in, ctx, allowed); err != nil || stop {
						return stop, err
					}
				}
				break jumptable
			}
		}
	case hlir.ASSERT:
		if !evalCondition(o.Predicate, ctx, []ast.Effect{}) {
			err := assertionError{string(o.Message), o.Node}
			ctx.writeStderr(err.Error())
			return true, err
		}
	default:
		panic(fmt.Sprintf("Unrecognized op: %v", reflect.TypeOf(op).Name()))
	}
	return false, nil
}

func evalRegister(r hlir.Register, ctx *Context) interface{} {
	switch reg := r.(type) {
	case hlir.StringLiteral:
		return strings.Replace(string(reg), `\n`, "\n", -1)
	case hlir.IntLiteral:
		return int(reg)
	case hlir.LocalValue:
		v, ok := ctx.localValues[reg]
		if !ok {
			panic(fmt.Sprintf("No variable named %v", reg))
		}
		return v
	case hlir.LastFuncCallRetVal:
		// hlir numbers the calls so that it can attach different type
		// information to it. The VM doesn't care, it just wants the last
		// call.
		reg.CallNum = 0
		val, ok := ctx.lastFuncCallRetVal[reg]
		if !ok {
			panic(fmt.Sprintf("Unknown function return value %v (Known: %v)", reg, ctx.lastFuncCallRetVal))
		}
		return val
	case hlir.TempValue:
		v, ok := ctx.tempValue[reg]
		if !ok {
			panic(fmt.Sprintf("Unknown temp value: %v", reg))
		}
		return v
	case hlir.FuncArg:
		if reg.Reference {
			// If it's a reference parameter, we need to de-reference it
			// to get its values
			ctxptr := ctx.pointers[hlir.Pointer{reg}]
			ptr := ctxptr.r.(hlir.Pointer)
			return evalRegister(ptr.Register, ctxptr.ctx)
		}
		v, ok := ctx.funcArg[reg]
		if !ok {
			panic(fmt.Sprintf("Unknown FuncArg: %v", reg))
		}
		return v
	case hlir.Offset:
		lv, nctx := resolveOffset(reg, ctx)
		return evalRegister(lv, nctx)
	case hlir.Pointer:
		p, ok := ctx.pointers[reg]
		if !ok {
			return evalRegister(reg.Register, ctx)
		}
		return evalRegister(p.r, p.ctx)
	case hlir.SliceBasePointer:
		return evalRegister(reg.Register, ctx)
	default:
		panic(fmt.Sprintf("Unhandled type: %v", reflect.TypeOf(r).Name()))
	}
}

func evalCondition(cond hlir.Condition, ctx *Context, allowed []ast.Effect) bool {
	for _, inst := range cond.Body {
		if stop, err := runOp(inst, ctx, allowed); err != nil || stop {
			panic(err)
		}
	}
	r := evalRegister(cond.Register, ctx)
	if r == true || r == false {
		return r.(bool)
	}
	return r != 0
}

func resolveOffset(o hlir.Offset, ctx *Context) (hlir.Register, *Context) {
	switch b := o.Base.(type) {
	case hlir.LocalValue:
		bd := b
		if p, ok := ctx.pointers[hlir.Pointer{b}]; ok {
			bd = p.r.(hlir.LocalValue)
		}
		offset := evalRegister(o.Offset, ctx)
		if o.Scale == 16 {
			bd += hlir.LocalValue((offset.(int) * 2) + 1)
		} else {
			bd += hlir.LocalValue(offset.(int))
		}
		return bd, ctx
	case hlir.FuncArg:
		base, nctx := dereferencePointer(b, ctx)
	resolveouter:
		for {
			switch o := base.(type) {
			case hlir.LocalValue:
				break resolveouter
			case hlir.FuncArg:
				base, nctx = dereferencePointer(base, nctx)
			case hlir.SliceBasePointer:
				base = o.Register
			default:
				base, nctx = dereferencePointer(base, nctx)
			}

		}
		bd := base.(hlir.LocalValue)
		offset := evalRegister(o.Offset, ctx)
		if o.Scale == 16 {
			bd += hlir.LocalValue((offset.(int) * 2) + 1)
		} else {
			bd += hlir.LocalValue(offset.(int))
		}
		return bd, nctx
	default:
		panic(fmt.Sprintf("Unhandled base type for Offset register: %v", reflect.TypeOf(o.Base)))
	}
}

func dereferencePointer(r hlir.Register, ctx *Context) (hlir.Register, *Context) {
	switch reg := r.(type) {
	case hlir.Pointer:
		p, ok := ctx.pointers[reg]
		if ok {
			return p.r, p.ctx
		}
		return reg.Register, ctx
	case hlir.FuncArg:
		if v, ok := ctx.pointers[hlir.Pointer{reg}]; ok {
			if r, ok := v.r.(hlir.Pointer); ok {
				return dereferencePointer(r, v.ctx)
			}
			return v.r, v.ctx
		}
		panic("Could not dereference FuncArg")
	case hlir.SliceBasePointer:
		return r, ctx
	default:
		if v, ok := ctx.pointers[hlir.Pointer{reg}]; ok {
			return v.r, v.ctx
		}
		return r, ctx
	}
}

type assertionError struct {
	Message string
	src     ast.Node
}

func (a assertionError) Error() string {
	if a.Message == "" {
		return fmt.Sprintf("assertion %v failed", a.src.PrettyPrint(0))
	}
	return fmt.Sprintf("assertion %v failed: %s", a.src.PrettyPrint(0), string(a.Message))
}
