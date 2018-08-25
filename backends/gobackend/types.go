package gobackend

import (
	"fmt"
	"io"
	"reflect"

	"github.com/driusan/lang/parser/ast"
)

func convertType(t ast.Type) (string, error) {
	switch t.(type) {
	case ast.TypeLiteral:
		return t.TypeName(), nil
	default:
		return "", fmt.Errorf("Convert type %v not implemented", reflect.TypeOf(t))
	}
}
func convertTypeDefn(c *Context, w io.Writer, t ast.TypeDefn, lvl int) error {
	indent := nTabs(lvl)
	ty, err := convertType(t.ConcreteType)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%vtype %v %v\n", indent, t.Name, ty)
	return nil
}

func convertEnumTypeDefn(c *Context, w io.Writer, t ast.EnumTypeDefn, lvl int) error {
	indent := nTabs(lvl)
	if len(t.Parameters) != 0 {
		return fmt.Errorf("Parameterized enumerations not implemented")
	}

	fmt.Fprintf(w, "%vtype %v int\n", indent, t.Name)
	fmt.Fprintf(w, "%vconst (\n", indent)
	for i, v := range t.Options {
		fmt.Fprintf(w, "%v\t%v = %v(%d)\n", indent, v.Constructor, t.Name, i)
	}
	fmt.Fprintf(w, "%v)\n", indent)
	return nil
}

func convertEnumValue(c *Context, w io.Writer, ev ast.EnumValue, lvl int) error {
	indent := nTabs(lvl)
	if len(ev.Parameters) != 0 {
		return fmt.Errorf("Parameterized enumerations not implemented")
	}

	id := c.EnumMap.GetIndex(ev.Constructor)
	fmt.Fprintf(w, "%v%v(%d)", indent, ev.Type().TypeName(), id)

	return nil
}

func convertEnumOption(c *Context, w io.Writer, ev ast.EnumOption, lvl int) error {
	indent := nTabs(lvl)
	if len(ev.Parameters) != 0 {
		return fmt.Errorf("Parameterized enumerations not implemented")
	}

	id := c.EnumMap.GetIndex(ev)
	fmt.Fprintf(w, "%v%v(%d)", indent, ev.Type().TypeName(), id)

	return nil
}

func getConcreteType(t ast.Type) ast.Type {
	switch ty := t.(type) {
	case ast.UserType:
		return getConcreteType(ty.Typ)
	default:
		return ty
	}
}
