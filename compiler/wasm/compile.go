package wasm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/driusan/lang/compiler/hlir"
	"github.com/driusan/lang/parser/ast"
)

func Parse(src string) (Module, error) {
	nodes, ti, callables, err := ast.Parse(src)
	if err != nil {
		return Module{}, err
	}

	ctx := NewContext(callables)
	ctx.typeinfo = ti

	// Identify required type information before code generation
	// for the functions.
	enums := make(hlir.EnumMap)
	for _, v := range nodes {
		switch v.(type) {
		case ast.EnumTypeDefn:
			_, opts, _, err := hlir.Generate(v, ti, callables, enums)
			if err != nil {
				return Module{}, err
			}
			for k, v := range opts {
				enums[k] = v
			}
		default:
			// Handled below
		}

	}

	// Generate the IR for the functions.
	for _, v := range nodes {
		switch v.(type) {
		case ast.FuncDecl:
			fnc, _, registers, err := hlir.Generate(v, ti, callables, enums)
			if err != nil {
				return Module{}, err
			}
			ctx.registerData = registers

			rfnc, err := Generate(fnc, &ctx)
			if err != nil {
				return Module{}, err
			}
			ctx.Functions = append(ctx.Functions, rfnc)
		case ast.TypeDefn, ast.EnumTypeDefn:
			// No IR for types, we've already verified them.
		default:
			panic("Unhandled AST node type for code generation")
		}
	}

	mem := ctx.memTop
	memo := Memory{}
	if mem > 0 || ctx.curFuncMaxMem > 0 {
		// convert from bytes to 64k (wasm size) pages
		mem /= (64 * 1024)
		mem += 1

		memo.Size = mem
		memo.Name = "mem"
	}
	m := Module{
		Imports: ctx.Imports,
		Memory:  memo,
		Data:    ctx.Data,
		Funcs:   ctx.Functions,
	}
	if ctx.needsGlobal {
		m.Globals = []Global{
			Global{
				Mutable:      true,
				Type:         i32,
				InitialValue: ctx.memTop,
			},
		}
	}
	return m, nil
}

func Generate(hlfnc hlir.Func, ctx *Context) (Func, error) {
	ret := Func{
		Name:      hlfnc.Name,
		Signature: ctx.GetSignature(hlfnc.Name),
	}

	ctx.curFuncMemVariables = make(map[hlir.Register]uint)
	ctx.curFuncLocalVariables = make(map[hlir.Register]uint)

	ctx.curFuncNumArgs = ctx.GetNumArgs(hlfnc.Name)
	ctx.curFuncNumReturns = ctx.GetNumReturns(hlfnc.Name)

	// Figure out which variables need to be stored in memory
	for _, op := range hlfnc.Body {
		regs := op.Registers()
		for idx, reg := range regs {
			switch v := reg.(type) {
			case hlir.Offset:
				if _, ok := v.Offset.(hlir.IntLiteral); !ok {
					switch base := v.Container.Typ.(type) {
					case ast.ArrayType:
						for i := ast.IntLiteral(0); i < base.Size; i++ {
							lv := v.Base.(hlir.LocalValue)
							lv += hlir.LocalValue(i)
							ctx.addMemoryVar(lv)
						}
					case ast.SliceType:
						switch b := v.Base.(type) {
						case hlir.LocalValue:
							toconvert, ok := ctx.registerData[b-1]
							if !ok {
								panic("Could not get size of local slice")
							}
							for i := uint(b - 1); i <= toconvert.SliceSize; i++ {
								lv := v.Base.(hlir.LocalValue)
								lv += hlir.LocalValue(i - 1)
								ctx.addMemoryVar(lv)
							}
						}
					default:
						panic(fmt.Sprintf("Unhandle offset base type: %v for base: %v", reflect.TypeOf(base), reg))
					}
				}
			case hlir.Pointer:
				info := ctx.registerData[v]
				switch info.Creator.Typ.(type) {
				case ast.SliceType:
					size := regs[idx-1]
					toconvert, ok := ctx.registerData[size]
					if !ok {
						panic("Could not determine slice size")
					}
					for i := uint(1); i <= toconvert.SliceSize; i++ {
						lv := size.(hlir.LocalValue)
						lv += hlir.LocalValue(i)
						ctx.addMemoryVar(lv)
					}
				case ast.TypeLiteral:
					ctx.addMemoryVar(v.Register.(hlir.LocalValue))
					// FIXME: Implement reference handling
				default:
					panic(fmt.Sprintf("Unhandled pointer type: %v", reflect.TypeOf(info.Creator.Typ)))
				}
			case hlir.LastFuncCallRetVal:

				if v.RetNum == 1 {
					// This is a multireturn function, and we've encountered the second return, so
					// add both the first and second arguments as memory variables
					ctx.curFuncMemVariables[v] = 4
					v.RetNum = 0
					ctx.curFuncMemVariables[v] = 0
				} else if v.RetNum > 1 {
					// For anything after the second, assume 0 and 1 are already there, so just
					// add this one.
					ctx.curFuncMemVariables[v] = 4 * v.RetNum
				}
			}
		}
	}

	for i := uint(0); i < hlfnc.NumLocals; i++ {
		var t VarType
		lv := hlir.LocalValue(i)
		typeinfo := ctx.registerData[lv]
		if typeinfo.TypeInfo.Size > 4 {
			t = i64
		} else {
			t = i32
		}

		_, ok := ctx.curFuncMemVariables[lv]
		if !ok {
			// Not a memory variable, so add it to the signature and add a map to the WASM
			// local index.
			if _, ok := ctx.curFuncLocalVariables[lv]; ok {
				continue
			}
			ctx.curFuncLocalVariables[lv] = uint(len(ctx.curFuncLocalVariables)) + ctx.curFuncNumArgs
			ret.Signature = append(ret.Signature, Variable{t, Local, fmt.Sprintf("LV%d", i)})
		}
	}

	for _, opi := range hlfnc.Body {
		ops, err := evaluateOp(opi, ctx)
		if err != nil {
			return Func{}, err
		}
		ret.Body = append(ret.Body, ops...)

	}
	if ctx.curFuncNumReturns > 0 && hlfnc.Body[len(hlfnc.Body)-1] != (hlir.RET{}) {
		ret.Body = append(ret.Body, Unreachable{})
	}
	return ret, nil
}

func evaluateOp(opi hlir.Opcode, ctx *Context) ([]Instruction, error) {
	switch op := opi.(type) {
	case hlir.CALL:
		ops := []Instruction{}
		funcs := ctx.callables[string(op.FName)]
		if len(funcs) != 1 {
			panic("Multiple dispatch not implemented")
		}
		calleeargs := funcs[0].GetArgs()
		stackRet := false
		ret := funcs[0].ReturnTuple()
		switch {
		case len(ret) == 1:
			words := strings.Fields(ret[0].Type().TypeName())
			if len(words) > 1 {
				stackRet = true
			}
		case len(ret) > 1:
			stackRet = true
		}

		idx := 0
		for _, v := range calleeargs {
			switch v.Typ.(type) {
			case ast.SliceType:
				size := op.Args[idx]
				base := op.Args[idx+1]
				if s := getValue(size, ctx); s != nil {
					ops = append(ops, s...)
				}
				switch basearray := base.(type) {
				case hlir.Pointer:
					switch base := basearray.Register.(type) {
					case hlir.LocalValue:
						globaloffset, memvar := ctx.curFuncMemVariables[base]
						if !memvar {
							panic(fmt.Sprintf("Could not find variable memory location for %v", base))
						}
						ops = append(ops, GetGlobal(0))
						if globaloffset > 0 {
							ops = append(ops, I32Const(globaloffset), I32Add{})
						}
					default:
						panic(fmt.Sprintf("Unhandled slice type base: %v", reflect.TypeOf(basearray.Register)))
					}
				case hlir.FuncArg:
					ops = append(ops, GetLocal(basearray.Id))
				default:
					panic(fmt.Sprintf("Unhandled type for slice: %v", reflect.TypeOf(basearray)))
				}
				idx += 2
			default:
				words := strings.Fields(string(v.Type().TypeName()))
				for _, word := range words {
					typeinfo, ok := ctx.typeinfo[word]
					if !ok {
						panic(fmt.Sprintf("Could not get type info for argument: %v, %v", word, words))
					}

					// Determine the size context based on the expected argument type
					// for IntLiterals
					if typeinfo.Size > 4 {
						ctx.size64 = true
					} else {
						ctx.size64 = false
					}

					reg := op.Args[idx]
					if op := getValueForCall(reg, ctx); op != nil {
						ops = append(ops, op...)
					}
					switch reg.(type) {
					case hlir.LocalValue:
						// Determine the size of the last variable used in case PrintInt
						// needs to wrap
						reg := hlir.LocalValue(idx)
						info := ctx.registerData[reg]
						if info.TypeInfo.Size > 4 {
							ctx.size64 = true
						} else {
							ctx.size64 = false
						}
					}
					idx++
				}
			}
		}

		if op.FName == "PrintInt" && ctx.size64 {
			ops = append(ops, I32WrapI64{})
		}

		ops = append(ops, Call{string(op.FName)})
		// Automatically import PrintString..
		// (This is a temporary hack until there's proper namespaces.)
		if op.FName == "PrintString" || op.FName == "PrintInt" {
			ctx.AddImport(
				"stdlib",
				string(op.FName),
				Signature{
					Variable{i32, Param, ""},
				},
			)
		} else if op.FName == "PrintByteSlice" {
			ctx.AddImport(
				"stdlib",
				string(op.FName),
				Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
				},
			)
		} else if op.FName == "len" {
			ctx.AddImport(
				"stdlib",
				string(op.FName),
				Signature{
					Variable{i32, Param, ""},
					Variable{i32, Param, ""},
					Variable{i32, Result, ""},
				},
			)
		}

		// If the function returned values in memory instead of on the stack, handle them appropriately..
		if stackRet {
			ops = append(ops, Drop{})
			ctx.lastCallRetNeedsMem = true
		} else {
			ctx.lastCallRetNeedsMem = false
		}
		return ops, nil
	case hlir.MOV:
		ops := []Instruction{}
		switch d := op.Dst.(type) {
		case hlir.FuncRetVal:
			typeinfo := ctx.registerData[d]
			if typeinfo.TypeInfo.Size > 4 {
				ctx.size64 = true
			} else {
				ctx.size64 = false
			}
			if ctx.curFuncRetNeedsMem {
				// If the function is being returned in memory instead of the return
				// stack (things that require multireturn), add a prelude to save the
				// value in memory.
				ops = append(ops, GetGlobal(0))
				if d >= 1 {
					ops = append(ops, I32Const(d*4), I32Add{})
				}
			}
			if op := getValue(op.Src, ctx); op != nil {
				ops = append(ops, op...)
			}
			if ctx.curFuncRetNeedsMem {
				ops = append(ops, storeOp(op.Dst, ctx))
			}
		case hlir.TempValue:
			// Do nothing, it's already on the stack
		case hlir.LocalValue:
			typeinfo := ctx.registerData[d]
			if typeinfo.TypeInfo.Size > 4 {
				ctx.size64 = true
			} else {
				ctx.size64 = false
			}
			globaloffset, memvar := ctx.curFuncMemVariables[d]
			val := getValue(op.Src, ctx)
			if val != nil && !memvar {
				ops = append(ops, val...)
			}
			if typeinfo.Name == "__[0]" {
				// FIXME: I have no idea where the extra _[0] comes from
				ops = append(ops, Drop{})
			} else {
				if memvar {
					ops = append(ops, GetGlobal(0))
					if globaloffset != 0 {
						ops = append(ops, I32Const(globaloffset), I32Add{})
					}
					ops = append(ops, val...)
					ops = append(ops, storeOp(op.Dst, ctx))
				} else {
					ops = append(ops, SetLocal(ctx.LocalIndex(d)))
				}
			}
		case hlir.FuncArg:
			ops = append(ops, getValueForReferenceVariableSave(op.Dst, ctx)...)
			ops = append(ops, getValue(op.Src, ctx)...)
			ops = append(ops, storeOp(op.Dst, ctx))
			// FIXME: Implement
		case hlir.Offset:
			lv := d.Base.(hlir.LocalValue)
			if addr, ok := ctx.curFuncMemVariables[lv]; ok {
				ops = append(ops, GetGlobal(0))
				if addr > 0 {
					ops = append(ops, I32Const(addr), I32Add{})
				}
				scale := d.Scale
				if scale == 0 {
					scale = 4
				}
				switch o := d.Offset.(type) {
				case hlir.IntLiteral:
					if o > 0 {
						ops = append(ops, getValue(d.Offset, ctx)...)
						if scale != 1 {
							ops = append(ops, I32Const(scale), I32Mul{})
						}
					}
				default:
					ops = append(ops, getValue(d.Offset, ctx)...)
					if scale != 1 {
						ops = append(ops, I32Const(scale), I32Mul{})
					}
					ops = append(ops, I32Add{})
				}
				ops = append(ops, getValue(op.Src, ctx)...)
				ops = append(ops, storeOp(d.Base, ctx))
			} else {
				switch o := d.Offset.(type) {
				case hlir.IntLiteral:
					lv += hlir.LocalValue(o)
				default:
					panic("Unhandled offset type in WASM")
				}
				ops = append(ops, getValue(op.Src, ctx)...)
				ops = append(ops, SetLocal(ctx.LocalIndex(lv)))
			}
		default:
			panic(fmt.Sprintf("Unhandled dst register type in op %v: %v", op, reflect.TypeOf(d)))
		}
		return ops, nil
	case hlir.ADD:
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		typeinfo := ctx.registerData[op.Left]
		if typeinfo.TypeInfo.Size > 4 {
			ops = append(ops, I64Add{})
		} else {
			ops = append(ops, I32Add{})
		}
		switch op.Dst.(type) {
		case hlir.TempValue:
			// Do nothing, leave it on the stack..
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.SUB:
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		// FIXME: Be smarter about i32 vs i64
		ops = append(ops, I32Sub{})
		switch op.Dst.(type) {
		case hlir.TempValue:
			// Do nothing, leave it on the stack..
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.MUL:
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		// FIXME: Be smarter about i32 vs i64
		ops = append(ops, I32Mul{})
		switch op.Dst.(type) {
		case hlir.TempValue:
			// Do nothing, leave it on the stack..
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.DIV:
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		// FIXME: Be smarter about i32 vs i64
		ops = append(ops, I32Div_S{})
		switch op.Dst.(type) {
		case hlir.TempValue:
			// Do nothing, leave it on the stack..
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.MOD:
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		// FIXME: Be smarter about i32 vs i64
		ops = append(ops, I32Rem_S{})
		switch op.Dst.(type) {
		case hlir.TempValue:
			// Do nothing, leave it on the stack..
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.IF:
		ops := []Instruction{}
		for _, cond := range op.Condition.Body {
			condops, err := evaluateOp(cond, ctx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, condops...)
		}
		switch op.Condition.Register.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled Condition type: %v", reflect.TypeOf(op.Condition.Register)))
		}
		ops = append(ops, If{})
		for _, b := range op.Body {
			bops, err := evaluateOp(b, ctx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, bops...)
		}
		if op.ElseBody != nil {
			ops = append(ops, Else{})
			for _, b := range op.ElseBody {
				bops, err := evaluateOp(b, ctx)
				if err != nil {
					return nil, err
				}
				ops = append(ops, bops...)
			}
		}
		ops = append(ops, End{})
		return ops, nil
	case hlir.LOOP:
		ops := []Instruction{Block{}}
		for _, op := range op.Initializer {
			initop, err := evaluateOp(op, ctx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, initop...)
		}

		ops = append(ops, Loop{})
		for _, cond := range op.Condition.Body {
			condops, err := evaluateOp(cond, ctx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, condops...)
		}

		switch op.Condition.Register.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled Condition type: %v", reflect.TypeOf(op.Condition.Register)))
		}
		// Negate the condition, because br_if breaks out if the condition is true.
		ops = append(ops, I32EQZ{})
		ops = append(ops, BrIf(1))
		for _, b := range op.Body {
			bops, err := evaluateOp(b, ctx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, bops...)
		}
		ops = append(ops, Br(0))
		ops = append(ops, End{}) // Loop
		ops = append(ops, End{}) // Block
		return ops, nil
	case hlir.JumpTable:
		ops := []Instruction{}

		for i, condition := range op {
			for _, cond := range condition.Condition.Body {
				condops, err := evaluateOp(cond, ctx)
				if err != nil {
					return nil, err
				}
				ops = append(ops, condops...)
			}
			ops = append(ops, If{})
			for _, op := range condition.Body {
				bops, err := evaluateOp(op, ctx)
				if err != nil {
					return nil, err
				}
				ops = append(ops, bops...)
			}
			if i != len(op)-1 {
				ops = append(ops, Else{})
			}
		}
		for range op {
			ops = append(ops, End{})
		}
		return ops, nil

	case hlir.GT:
		// FIXME: This shouldn't assume signed i32
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		ops = append(ops, I32GT_S{})
		switch op.Dst.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.GEQ:
		// FIXME: This shouldn't assume signed i32
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		typeinfo := ctx.registerData[op.Left]
		if typeinfo.TypeInfo.Size > 4 {
			if typeinfo.TypeInfo.Signed {
				ops = append(ops, I64GE_S{})
			} else {
				ops = append(ops, I64GE_U{})
			}
		} else {
			if typeinfo.TypeInfo.Signed {
				ops = append(ops, I32GE_S{})
			} else {
				ops = append(ops, I32GE_U{})
			}
		}
		switch op.Dst.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.EQ:
		// FIXME: This shouldn't assume signed i32
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		ops = append(ops, I32EQ{})
		switch op.Dst.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.NEQ:
		// FIXME: This shouldn't assume signed i32
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		ops = append(ops, I32NE{})
		switch op.Dst.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.LTE:
		// FIXME: This shouldn't assume signed i32
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		ops = append(ops, I32LE_S{})
		switch op.Dst.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.LT:
		// FIXME: This shouldn't assume signed i32
		ops := []Instruction{}
		if op := getValue(op.Left, ctx); op != nil {
			ops = append(ops, op...)
		}
		if op := getValue(op.Right, ctx); op != nil {
			ops = append(ops, op...)
		}

		ops = append(ops, I32LT_S{})
		switch op.Dst.(type) {
		case hlir.TempValue:
		default:
			panic(fmt.Sprintf("Unhandled dst operand %v", reflect.TypeOf(op.Dst)))
		}
		return ops, nil
	case hlir.RET:
		if ctx.curFuncRetNeedsMem {
			return []Instruction{GetGlobal(0), Return{}}, nil
		}
		return []Instruction{Return{}}, nil
	default:
		panic(fmt.Sprintf("Unhandled op type while generating function: %v, %v", reflect.TypeOf(opi), opi))
	}
}

func getValue(reg hlir.Register, ctx *Context) []Instruction {
	switch v := reg.(type) {
	case hlir.StringLiteral:
		i := ctx.GetLiteral(string(v))
		return []Instruction{I32Const(i)}
	case hlir.IntLiteral:
		if ctx.size64 {
			return []Instruction{I64Const(v)}
		}
		return []Instruction{I32Const(v)}
	case hlir.FuncArg:
		if v.Reference {
			return []Instruction{GetLocal(v.Id), loadOp(v, ctx)}
		}
		return []Instruction{GetLocal(v.Id)}
	case hlir.LocalValue:
		memoffset, memvar := ctx.curFuncMemVariables[v]
		if !memvar {
			return []Instruction{GetLocal(ctx.LocalIndex(v))}
		}

		ops := []Instruction{
			GetGlobal(0),
		}
		if memoffset > 0 {
			ops = append(ops, I32Const(memoffset), I32Add{})
		}
		ops = append(ops, loadOp(v, ctx))
		return ops
	case hlir.FuncRetVal:
		if ctx.lastCallRetNeedsMem {
			var ops []Instruction

			ops = append(ops, GetGlobal(0))
			if v > 0 {
				ops = append(ops, I32Const(4*v), I32Add{})
			}
			ops = append(ops, loadOp(v, ctx))

			return ops
		}
		return nil
	case hlir.TempValue:
		return nil
	case hlir.Offset:
		var off hlir.IntLiteral
		var ok bool
		memoffset, memvar := ctx.curFuncMemVariables[v.Base]
		if memvar {
			goto useMem
		}
		switch base := v.Base.(type) {
		case hlir.LocalValue:
			off, ok = v.Offset.(hlir.IntLiteral)
			if !ok {
				goto useMem
			}
			return []Instruction{GetLocal(uint(base) + uint(off) + ctx.curFuncNumArgs)}
		case hlir.FuncArg:
			scale := v.Scale
			if v.Scale == 0 {
				scale = 4
			}

			switch off := v.Offset.(type) {
			case hlir.IntLiteral:
				if off == 0 {
					return []Instruction{GetLocal(base.Id), loadOp(v, ctx)}
				}
				return []Instruction{GetLocal(base.Id), I32Const(uint(off) * uint(scale)), I32Add{}, loadOp(v, ctx)}
			default:
				ops := getValue(v.Offset, ctx)
				if scale != 1 {
					ops = append(ops, I32Const(scale), I32Mul{})
				}
				ops = append(ops, GetLocal(base.Id), I32Add{}, loadOp(v, ctx))
				return ops
			}
		default:
			panic("Unhandled type of offset")
		}
	useMem:
		var ops []Instruction
		// Calculate index
		if v.Offset != hlir.IntLiteral(0) {
			ops = append(ops, getValue(v.Offset, ctx)...)
			// Scale the index to the appropriate memory location.
			if v.Scale == hlir.IntLiteral(0) {
				ops = append(ops, I32Const(4))
				ops = append(ops, I32Mul{})
			} else if v.Scale == hlir.IntLiteral(1) {
				// there's no scaling multiplication for scale 1.
			} else {
				ops = append(ops, getValue(v.Scale, ctx)...)
				ops = append(ops, I32Mul{})
			}
		}

		// Add it to the base register
		ops = append(ops, GetGlobal(0))
		if memoffset != 0 {
			ops = append(ops, I32Const(memoffset), I32Add{})
		}
		if v.Offset != hlir.IntLiteral(0) {
			ops = append(ops, I32Add{})
		}
		ops = append(ops, loadOp(v.Base, ctx))
		return ops
	case hlir.Pointer:
		ctx.needsGlobal = true
		return getValue(v.Register, ctx)
	case hlir.LastFuncCallRetVal:
		memoffset, memvar := ctx.curFuncMemVariables[v]
		if memvar {
			ret := []Instruction{GetGlobal(0)}
			if memoffset != 0 {
				ret = append(ret, I32Const(memoffset), I32Add{})
			}
			ret = append(ret, loadOp(v, ctx))
			return ret
		}
		return nil
	default:
		panic(fmt.Sprintf("Unhandled register type: %v", reflect.TypeOf(v)))
	}
}

func getValueForReferenceVariableSave(reg hlir.Register, ctx *Context) []Instruction {
	switch a := reg.(type) {
	case hlir.FuncArg:
		return []Instruction{GetLocal(a.Id)}
	default:
		return getValue(reg, ctx)
	}
}

func getValueForCall(reg hlir.Register, ctx *Context) []Instruction {
	switch a := reg.(type) {
	case hlir.FuncArg:
		return []Instruction{GetLocal(a.Id)}
	case hlir.Pointer:
		v := a.Register.(hlir.LocalValue)
		memoffset, memvar := ctx.curFuncMemVariables[v]
		if !memvar {
			// FIXME: Moving some locals to memory might have screwed up the index of locals after
			// the ones stored in memory.
			return []Instruction{GetLocal(ctx.LocalIndex(v))}
		}

		ops := []Instruction{
			GetGlobal(0),
		}
		if memoffset > 0 {
			ops = append(ops, I32Const(memoffset), I32Add{})
		}
		return ops
	default:
		return getValue(reg, ctx)
	}
}

func loadOp(src hlir.Register, ctx *Context) Instruction {
	data, ok := ctx.registerData[src]
	if !ok {
		return I32Load{}
	}
	info := data.TypeInfo
	if info.Size == 1 {
		if info.Signed {
			return I32Load8S{}
		} else {
			return I32Load8U{}
		}
	}
	return I32Load{}
}
func storeOp(dst hlir.Register, ctx *Context) Instruction {
	data, ok := ctx.registerData[dst]
	if !ok {
		panic(fmt.Sprintf("Could not find register data for %v", dst))
	}
	info := data.TypeInfo
	if info.Size == 1 {
		return I32Store8{}
	}
	return I32Store{}
}
