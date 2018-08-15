package ast

import (
	"fmt"
	//	"reflect"
)

func IsLiteral(v Value) bool {
	switch t := v.(type) {
	case IntLiteral, BoolLiteral, StringLiteral, ArrayLiteral:
		return true
	case TupleValue:
		for _, c := range t {
			if !IsLiteral(c) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (c *Context) IsCompatibleType(typ Type, v Value) error {
	switch t := typ.(type) {
	case SumType:
		// A SumType is compatible if the value is compatible with
		// at least one option.
		for _, subtype := range t {
			if err := c.IsCompatibleType(subtype, v); err == nil {
				// If it's compatible with one of the options for a sum
				// type, it's compatible
				return nil
			}
		}
		return fmt.Errorf("Value %v is not compatible with sum type %v", v, typ.TypeName())
	case ArrayType:
		// Easy case, for variables
		if v.Type().TypeName() == t.TypeName() {
			return nil
		}

		// Otherwise, an array type with a literal of the same size
		// where every member is compatible with the base type
		v2, ok := v.(ArrayLiteral)
		if !ok {
			return fmt.Errorf("%v is not compatible with %v", v, typ.TypeName())
		}
		if len(v2.Values) != int(t.Size) {
			return fmt.Errorf("%v is not compatible with %v", v.Type().TypeName(), typ.TypeName())
		}

		for i, val := range v2.Values {
			if err := c.IsCompatibleType(t.Base, val); err != nil {
				return fmt.Errorf("Array Type error at index %d: %v", i, err)
			}
		}
		return nil
	case SliceType:
		// Easy case, for variables
		if v.Type().TypeName() == t.TypeName() {
			return nil
		}

		// Otherwise, an array type with a literal of the same size
		// where every member is compatible with the base type
		v2, ok := v.(ArrayLiteral)
		if !ok {
			return fmt.Errorf("%v is not compatible with %v", v, typ.TypeName())
		}

		for i, val := range v2.Values {
			if err := c.IsCompatibleType(t.Base, val); err != nil {
				return fmt.Errorf("Array Type error at index %d: %v", i, err)
			}
		}
		return nil
	case UserType:
		if IsLiteral(v) {
			return c.IsCompatibleType(t.Typ, v)
		}
		switch t.Typ.(type) {
		case SumType:
			return c.IsCompatibleType(t.Typ, v)
		default:
			if v.Type().TypeName() != typ.TypeName() {
				return fmt.Errorf("%v is not compatible with %v", v.Type().TypeName(), typ.TypeName())
			}
		}
	case TupleType:
		tv, ok := v.(TupleValue)
		if !ok {
			return fmt.Errorf("tuples is not compatible with non-tuple value")
		}
		if len(tv) != len(t) {
			return fmt.Errorf("incorrect size for tuple value")
		}
		for i, comp := range t {
			if err := c.IsCompatibleType(comp.Type(), tv[i]); err != nil {
				return fmt.Errorf("tuple component %d: %v", i, err)
			}
		}
		return nil
	}
	t, ok := c.Types[typ.TypeName()]
	if !ok {
		panic(fmt.Sprintf("Could not find type information for %v", typ.TypeName()))
	}

	switch t2 := v.(type) {
	case BoolLiteral:
		if t.ConcreteType == TypeLiteral("bool") {
			return nil
		}
		return fmt.Errorf("Can not assign bool to %v", t.Name)
	case IntLiteral:
		if t.ConcreteType == nil {
			panic(fmt.Sprintf("No concrete type for literal: %v", t))
		}
		switch t.ConcreteType.TypeName() {
		case "int":
			return nil
		case "uint8", "byte":
			if t2 >= 0 && t2 < 256 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between 0 and 255", t2)
		case "uint16":
			if t2 >= 0 && t2 < 65536 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between 0 and 65535", t2)
		case "uint32":
			if t2 >= 0 && t2 < 4294967296 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between 0 and 4,294,967,295", t2)
		case "uint64":
			// FIXME: The upper range check doesn't work because IntLiteral
			// is an int64 type. Replace with varint?
			if t2 >= 0 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between 0 and 18,446,744,073,709,551,615", t2)
		case "int8":
			if t2 >= -128 && t2 < 128 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between -128 and 127", t2)
		case "int16":
			if t2 >= -32768 && t2 < 32768 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between -32,768 and 32,767", t2)
		case "int32":
			if t2 >= -2147483648 && t2 < 2147483648 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between -2,147,483,648 and 2,147,483,647", t2)
		case "int64":
			if t2 >= -9223372036854775808 && t2 <= 9223372036854775807 {
				return nil
			}
			return fmt.Errorf("value (%d) must be between -9,223,372,036,854,775,808 and 9,223,372,036,854,775,807", t2)
		default:
			return fmt.Errorf("Can not assign int to %v", t.Name)
		}
	case StringLiteral:
		if t.ConcreteType == TypeLiteral("string") {
			return nil
		}
		return fmt.Errorf("Can not assign string to %v", t.Name)
	case ArrayLiteral:
		if t.ConcreteType == nil {
			return fmt.Errorf("%v does not have a concrete type", t)
		}
		if t.ConcreteType.TypeName() == v.Type().TypeName() {
			return nil
		}

		// Check if each element in the literal is compatible.
		if st, ok := t.ConcreteType.(SliceType); ok {
			// Fake an array type comparison instead of a slice type.
			return c.IsCompatibleType(ArrayType{st.Base, IntLiteral(len(t2.Values))}, v)
		}
		for _, el := range t2.Values { // t2=ArrayLiteral
			if err := c.IsCompatibleType(el.Type(), el); err != nil {
				return err
			}
		}
		return nil
	default:
		if v.Type().TypeName() == t.Name {
			return nil
		}
		return fmt.Errorf("Incompatible type for non-literal")
	}
}

type StringLiteral string

func (v StringLiteral) Value() interface{} {
	return v
}

func (s StringLiteral) Node() Node {
	return s
}

func (s StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral(%v)", string(s))
}

func (s StringLiteral) TypeName() string {
	return "string"
}
func (s StringLiteral) Type() Type {
	return TypeLiteral("string")
}

func (s StringLiteral) PrettyPrint(lvl int) string {
	return fmt.Sprintf(`%v"%v"`, nTabs(lvl), string(s))
}

type IntLiteral int64

func (v IntLiteral) Value() interface{} {
	return v
}

func (s IntLiteral) Node() Node {
	return s
}

func (i IntLiteral) String() string {
	return fmt.Sprintf("IntLiteral(%d)", i)
}

func (i IntLiteral) Type() Type {
	return TypeLiteral("int")
}

func (i IntLiteral) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%d", nTabs(lvl), int64(i))
}

type BoolLiteral bool

func (v BoolLiteral) BoolValue() bool {
	return bool(v)
}

func (v BoolLiteral) Value() interface{} {
	return v
}

func (b BoolLiteral) Node() Node {
	return b
}

func (b BoolLiteral) String() string {
	if b {
		return "BoolLiteral(true)"
	}
	return "BoolLiteral(false)"
}

func (b BoolLiteral) Type() Type {
	return TypeLiteral("bool")
}

func (s BoolLiteral) PrettyPrint(lvl int) string {
	return fmt.Sprintf("%v%v", nTabs(lvl), bool(s))
}
