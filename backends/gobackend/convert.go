package gobackend

import (
	"fmt"
	"io"
	"reflect"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/types"
)

func nTabs(lvl int) string {
	rv := ""
	for i := 0; i < lvl; i++ {
		rv += "\t"
	}
	return rv
}

// Convert writes node to w as Go source code
// Returns a map containing the names of imports which are
// needed to compile the go code or an error.
func Convert(w io.Writer, node ast.Node, emap types.EnumMap) (map[string]bool, error) {
	context := NewContext()
	context.EnumMap = emap
	err := convert(context, w, node, 0)
	if err != nil {
		return nil, err
	}
	return context.importsNeeded, nil
}

// convert writes node to w as Go source code, indented
// by level n
func convert(c *Context, w io.Writer, node ast.Node, lvl int) error {
	switch n := node.(type) {
	case ast.FuncDecl:
		return convertFuncDecl(c, w, n, lvl)
	case ast.TypeDefn:
		return convertTypeDefn(c, w, n, lvl)
	case ast.EnumTypeDefn:
		return convertEnumTypeDefn(c, w, n, lvl)
	case ast.EnumValue:
		return convertEnumValue(c, w, n, lvl)
	case ast.EnumOption:
		return convertEnumOption(c, w, n, lvl)
	case ast.FuncCall:
		return convertFuncCall(c, w, n, lvl)
	case ast.Assertion:
		return convertAssertion(c, w, n, lvl)
	case ast.IntLiteral:
		fmt.Fprintf(w, "%v%d", nTabs(lvl), n)
		return nil
	case ast.BoolLiteral:
		if n == false {
			fmt.Fprintf(w, "%vfalse", nTabs(lvl))
		} else {
			fmt.Fprintf(w, "%vtrue", nTabs(lvl))
		}
		return nil
	case ast.StringLiteral:
		fmt.Fprintf(w, `%v"%s"`, nTabs(lvl), string(n))
		return nil
	case ast.VarWithType:
		fmt.Fprintf(w, "%v%v", nTabs(lvl), c.GetVarName(n))
		return nil
	case ast.AdditionOperator:
		return convertOperator(c, w, n.Left, n.Right, "+", lvl)
	case ast.SubtractionOperator:
		return convertOperator(c, w, n.Left, n.Right, "-", lvl)
	case ast.MulOperator:
		return convertOperator(c, w, n.Left, n.Right, "*", lvl)
	case ast.DivOperator:
		return convertOperator(c, w, n.Left, n.Right, "/", lvl)
	case ast.ModOperator:
		return convertOperator(c, w, n.Left, n.Right, "%", lvl)
	case ast.GreaterComparison:
		return convertOperator(c, w, n.Left, n.Right, ">", lvl)
	case ast.GreaterOrEqualComparison:
		return convertOperator(c, w, n.Left, n.Right, ">=", lvl)
	case ast.EqualityComparison:
		return convertOperator(c, w, n.Left, n.Right, "==", lvl)
	case ast.NotEqualsComparison:
		return convertOperator(c, w, n.Left, n.Right, "!=", lvl)
	case ast.LessThanComparison:
		return convertOperator(c, w, n.Left, n.Right, "<", lvl)
	case ast.LessThanOrEqualComparison:
		return convertOperator(c, w, n.Left, n.Right, "<=", lvl)
	case ast.LetStmt:
		return convertLetStmt(c, w, n, lvl)
	case ast.MutStmt:
		return convertMutStmt(c, w, n, lvl)
	case ast.AssignmentOperator:
		return convertAssignmentOperator(c, w, n, lvl)
	case ast.ReturnStmt:
		return convertReturnStmt(c, w, n, lvl)
	case ast.WhileLoop:
		return convertWhileLoop(c, w, n, lvl)
	case ast.BlockStmt:
		return convertBlockStmt(c, w, n, lvl)
	case ast.IfStmt:
		return convertIfStmt(c, w, n, lvl)
	case ast.MatchStmt:
		return convertMatchStmt(c, w, n, lvl)
	default:
		return fmt.Errorf("Converting %v not implemented", reflect.TypeOf(node))
	}
}

func convertFuncDecl(c *Context, w io.Writer, fnc ast.FuncDecl, lvl int) error {
	if fnc.Name == "" {
		return fmt.Errorf("Invalid function")
	}
	indent := nTabs(lvl)
	fmt.Fprintf(w, "%vfunc %v", indent, fnc.Name)
	switch len(fnc.Args) {
	case 0:
		fmt.Fprintf(w, "()")
	default:
		fmt.Fprintf(w, "(")
		for i, arg := range fnc.Args {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			name := c.GetVarName(arg)
			typ, err := convertType(arg.Type())
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%v %v", name, typ)
		}
		fmt.Fprintf(w, ")")
	}
	switch len(fnc.Return) {
	case 0:
		// Nothing
	case 1:
		t, err := convertType(fnc.Return[0].Type())
		if err != nil {
			return err
		}
		fmt.Fprintf(w, " %v", t)
	default:
		return fmt.Errorf("Multiple returns not implemented")
	}
	fmt.Fprintf(w, " {\n")

	for _, stmt := range fnc.Body.Stmts {
		if err := convert(c, w, stmt, lvl+1); err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "%v}\n", indent)
	return nil
}
func convertAssertion(c *Context, w io.Writer, a ast.Assertion, lvl int) error {
	indent := nTabs(lvl)
	fmt.Fprintf(w, "%vif !(", indent)
	if err := convert(c, w, a.Predicate, 0); err != nil {
		return err
	}
	fmt.Fprintf(w, ") {\n")
	if a.Message != "" {
		fmt.Fprintf(w, "%v\tpanic(\"%v\")\n", indent, a.Message)
	} else {
		// FIXME: Add more detail about which assertion failed instead
		// of relying on Go to include a stack trace
		fmt.Fprintf(w, "%v\tpanic(\"Assertion failed\")\n", indent)
	}

	fmt.Fprintf(w, "%v}\n", indent)
	return nil
}

func convertFuncCall(c *Context, w io.Writer, fnc ast.FuncCall, lvl int) error {
	if fnc.Name == "" {
		return fmt.Errorf("Invalid function call")
	}
	indent := nTabs(lvl)

	if len(fnc.UserArgs) == 0 {
		fmt.Fprintf(w, "%v%v()", indent, fnc.Name)
	} else {
		switch fnc.Name {
		case "PrintString":
			fmt.Fprintf(w, "%vfmt.Printf(\"%vs\",\n", indent, "%")
			for _, arg := range fnc.UserArgs {
				if err := convert(c, w, arg, lvl+1); err != nil {
					return err
				}
				fmt.Fprintf(w, ",\n")
			}
			c.importsNeeded["fmt"] = true
			fmt.Fprintf(w, "%v)\n", indent)
		case "PrintInt":
			fmt.Fprintf(w, "%vfmt.Printf(\"%vd\",\n", indent, "%")
			for _, arg := range fnc.UserArgs {
				if err := convert(c, w, arg, lvl+1); err != nil {
					return err
				}
				fmt.Fprintf(w, ",\n")
			}
			c.importsNeeded["fmt"] = true
			fmt.Fprintf(w, "%v)\n", indent)
		default:
			fmt.Fprintf(w, "%v%v(\n", indent, fnc.Name)
			for _, arg := range fnc.UserArgs {
				if err := convert(c, w, arg, lvl+1); err != nil {
					return err
				}
				fmt.Fprintf(w, ",\n")
			}
			fmt.Fprintf(w, "%v)", indent)
		}
	}
	return nil
}

func convertVar(c *Context, w io.Writer, vname ast.VarWithType, val ast.Value, assignop string, lvl int) error {
	indent := nTabs(lvl)
	fmt.Fprintf(w, "%v", indent)
	if err := convert(c, w, vname, 0); err != nil {
		return err
	}
	// We don't want new variables to be created for variables on the assignment side
	c.newVar = false
	fmt.Fprintf(w, " %s ", assignop)
	if err := convert(c, w, val, 0); err != nil {
		return err
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func convertLetStmt(c *Context, w io.Writer, let ast.LetStmt, lvl int) error {
	c.newVar = true
	defer func() {
		c.newVar = false
	}()

	if let.Var.Name == "_" {
		return convertVar(c, w, let.Var, let.Val, "=", lvl)
	} else {
		return convertVar(c, w, let.Var, let.Val, ":=", lvl)
	}
}

func convertMutStmt(c *Context, w io.Writer, mut ast.MutStmt, lvl int) error {
	return convertVar(c, w, mut.Var, mut.InitialValue, ":=", lvl)
}

func convertAssignmentOperator(c *Context, w io.Writer, ass ast.AssignmentOperator, lvl int) error {
	switch v := ass.Variable.(type) {
	case ast.VarWithType:
		return convertVar(c, w, v, ass.Value, "=", lvl)
	default:
		return fmt.Errorf("Assignment to %v not implemented", reflect.TypeOf(ass.Variable))
	}
}

func convertOperator(c *Context, w io.Writer, left, right ast.Value, opsymbol string, lvl int) error {
	if err := convert(c, w, left, lvl); err != nil {
		return err
	}
	fmt.Fprintf(w, " %v ", opsymbol)
	if err := convert(c, w, right, 0); err != nil {
		return err
	}
	return nil
}

func convertReturnStmt(c *Context, w io.Writer, ret ast.ReturnStmt, lvl int) error {
	fmt.Fprintf(w, "%vreturn ", nTabs(lvl))
	if ret.Val != nil {
		if err := convert(c, w, ret.Val, 0); err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "\n")
	return nil
}

func convertWhileLoop(c *Context, w io.Writer, whl ast.WhileLoop, lvl int) error {
	// FIXME: Deal with shadowing properly and restore variables that were
	// shadowed to their old name after the loop is done.
	indent := nTabs(lvl)
	// We don't use the Go conditions but explicitly check with an if and break
	// so that we can support while (let x = ...) later
	fmt.Fprintf(w, "%vfor {\n", indent)
	fmt.Fprintf(w, "%v\tif !(", indent)
	if err := convert(c, w, whl.Condition, 0); err != nil {
		return err
	}
	fmt.Fprintf(w, ") {\n")
	fmt.Fprintf(w, "%v\tbreak;\n}\n", indent)
	if err := convert(c, w, whl.Body, lvl+1); err != nil {
		return err
	}

	fmt.Fprintf(w, "%v}\n", indent)
	return nil
}

func convertBlockStmt(c *Context, w io.Writer, block ast.BlockStmt, lvl int) error {
	for _, stmt := range block.Stmts {
		if err := convert(c, w, stmt, lvl); err != nil {
			return err
		}
	}
	return nil
}

func convertIfStmt(c *Context, w io.Writer, ifstmt ast.IfStmt, lvl int) error {
	// FIXME: Deal with shadowing properly and restore variables that were
	// shadowed to their old name after the if is done.
	// FIXME: Support let embedded in conditional
	indent := nTabs(lvl)
	fmt.Fprintf(w, "%vif (", indent)
	if err := convert(c, w, ifstmt.Condition, 0); err != nil {
		return err
	}
	fmt.Fprintf(w, ") {\n")
	if err := convert(c, w, ifstmt.Body, lvl+1); err != nil {
		return err
	}

	if len(ifstmt.Else.Stmts) > 0 {
		fmt.Fprintf(w, "%v} else {\n", indent)
		if err := convert(c, w, ifstmt.Else, lvl+1); err != nil {
			return err
		}
		fmt.Fprintf(w, "%v}\n", indent)
	} else {
		fmt.Fprintf(w, "%v}\n", indent)
	}
	return nil
}

func convertMatchStmt(c *Context, w io.Writer, matchstmt ast.MatchStmt, lvl int) error {
	indent := nTabs(lvl)
	switch {
	case matchstmt.IsEnumMatch():
		fmt.Fprintf(w, "%vswitch ", indent)
		if err := convert(c, w, matchstmt.Condition, 0); err != nil {
			return err
		}
		fmt.Fprintf(w, " {\n")
		for _, cse := range matchstmt.Cases {
			fmt.Fprintf(w, "%vcase ", indent)
			if err := convert(c, w, cse.Variable, lvl); err != nil {
				return err
			}
			fmt.Fprintf(w, ":\n")
			if err := convert(c, w, cse.Body, lvl+1); err != nil {
				return err
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "%v}\n", indent)
		return nil
	case matchstmt.IsSumTypeDestructure():
		return fmt.Errorf("Sum type destructuring not implemented")
	default:
		return fmt.Errorf("Match switch not implemented")
	}
}
