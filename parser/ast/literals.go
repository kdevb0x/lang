package ast

import (
	"fmt"
)

func IsLiteral(v Value) bool {
	switch v.(type) {
	case IntLiteral, BoolLiteral, StringLiteral:
		return true
	default:
		return false
	}
}

func IsCompatibleType(t TypeDefn, v Value) error {
	switch t2 := v.(type) {
	case BoolLiteral:
		if t.ConcreteType == "bool" {
			return nil
		}
		return fmt.Errorf("Can not assign bool to %v", t.Name)
	case IntLiteral:
		switch t.ConcreteType {
		case "int":
			return nil
		case "uint8":
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
			//if t2 >= 0 && t2 <  {
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
		if t.ConcreteType == "string" {
			return nil
		}
		return fmt.Errorf("Can not assign string to %v", t.Name)
	default:
		panic("Unhandled literal type in IsCompatibleType")
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

func (s StringLiteral) Type() Type {
	return "string"
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
	return "int"
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
	return "bool"
}
