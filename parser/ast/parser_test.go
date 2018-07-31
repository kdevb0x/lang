package ast

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/driusan/lang/parser/sampleprograms"
	"github.com/driusan/lang/parser/token"
)

func compare(v1, v2 Node) bool {
	// Easy types that don't have anything preventing them from being compared
	// with ==
	switch v1.(type) {
	case StringLiteral, BoolLiteral, IntLiteral,
		Variable,
		AdditionOperator, SubtractionOperator, AssignmentOperator,
		MulOperator, DivOperator,
		NotEqualsComparison,
		GreaterOrEqualComparison, LessThanOrEqualComparison,
		TypeLiteral, nil:
		return v1 == v2
	}

	if v1a, ok := v1.(EqualityComparison); ok {
		if v2a, ok := v2.(EqualityComparison); ok {
			return compare(v1a.Left, v2a.Left) && compare(v2a.Right, v2a.Right)
		}
		return false
	}
	if v1a, ok := v1.(VarWithType); ok {
		if v2a, ok := v2.(VarWithType); ok {
			return v1a.Name == v2a.Name && v1a.Reference == v2a.Reference && compare(v1a.Typ, v2a.Typ)
		}
		return false
	}
	if v1a, ok := v1.(SliceType); ok {
		if v2a, ok := v2.(SliceType); ok {
			return compare(v1a.Base, v2a.Base)
		}
		return false
	}
	if v1a, ok := v1.(LessThanComparison); ok {
		if v2a, ok := v2.(LessThanComparison); ok {
			return compare(v1a.Left, v2a.Left) && compare(v1a.Right, v2a.Right)
		}
		return false
	}
	if v1a, ok := v1.(GreaterComparison); ok {
		if v2a, ok := v2.(GreaterComparison); ok {
			return compare(v1a.Left, v2a.Left) && compare(v1a.Right, v2a.Right)
		}
		return false
	}
	if v1a, ok := v1.(Brackets); ok {
		if v2a, ok := v2.(Brackets); ok {
			return compare(v1a.Val, v2a.Val)
		}
		return false
	}
	if v1a, ok := v1.(*IfStmt); ok {
		if v2a, ok := v2.(IfStmt); ok {
			return compare(*v1a, v2a)
		} else if v2a, ok := v2.(*IfStmt); ok {
			return v1a == v2a
		}
		return false
	}
	if v1a, ok := v1.(LetStmt); ok {
		if v2a, ok := v2.(LetStmt); ok {
			return compare(v1a.Var, v2a.Var) && compare(v1a.Val, v2a.Val)
		}
		return false
	}
	if v1a, ok := v1.(MutStmt); ok {
		if v2a, ok := v2.(MutStmt); ok {
			return compare(v1a.Var, v2a.Var) && compare(v1a.InitialValue, v2a.InitialValue)
		}
		return false
	}
	if v1a, ok := v1.(ReturnStmt); ok {
		v2a, ok := v2.(ReturnStmt)
		if !ok {
			return false
		}
		if v1a.Val == nil && v2a.Val == nil {
			return true
		}
		return compare(v1a.Val, v2a.Val)
	}

	if v1a, ok := v1.(ModOperator); ok {
		v2a, ok := v2.(ModOperator)
		if !ok {
			return false
		}
		return compare(v1a.Left, v2a.Left) && compare(v1a.Right, v2a.Right)
	}

	if v1a, ok := v1.(FuncDecl); ok {
		v2a, ok := v2.(FuncDecl)
		if !ok {
			return false
		}
		if v1a.Name != v2a.Name {
			return false
		}
		if len(v1a.Args) != len(v2a.Args) {
			return false
		}
		for i := range v1a.Args {
			if compare(v1a.Args[i], v2a.Args[i]) == false {
				return false
			}
		}
		if len(v1a.Return) != len(v2a.Return) {
			return false
		}

		for i := range v1a.Return {
			if compare(v1a.Return[i], v2a.Return[i]) == false {
				return false
			}
		}
		if len(v1a.Effects) != len(v2a.Effects) {
			return false
		}

		for i := range v1a.Effects {
			if v1a.Effects[i] != v2a.Effects[i] {
				return false
			}
		}
		if compare(v1a.Body, v2a.Body) == false {
			return false
		}
		return true
	}
	if v1a, ok := v1.(BlockStmt); ok {
		v2a, ok := v2.(BlockStmt)
		if !ok {
			return false
		}
		if len(v1a.Stmts) != len(v2a.Stmts) {
			return false
		}
		for i := range v1a.Stmts {
			if compare(v1a.Stmts[i], v2a.Stmts[i]) == false {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(*BlockStmt); ok {
		v2a, ok := v2.(*BlockStmt)
		if !ok {
			return false
		}
		if v1a == v2a {
			return true
		}
		if v1a == nil || v2a == nil {
			return false
		}
		return compare(*v1a, *v2a)
	}

	if v1a, ok := v1.(FuncCall); ok {
		v2a, ok := v2.(FuncCall)
		if !ok {
			return false
		}
		if v1a.Name != v2a.Name {
			return false
		}
		if len(v1a.UserArgs) != len(v2a.UserArgs) {
			return false
		}
		for i := range v1a.UserArgs {
			arg1 := v1a.UserArgs[i].(Node)
			arg2 := v2a.UserArgs[i].(Node)
			if compare(arg1, arg2) == false {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(ReturnStmt); ok {
		v2a, ok := v2.(ReturnStmt)
		if !ok {
			return false
		}
		return compare(v1a.Val, v2a.Val)
	}
	if v1a, ok := v1.(WhileLoop); ok {
		v2a, ok := v2.(WhileLoop)
		if !ok {
			return false
		}
		return compare(v1a.Condition, v2a.Condition) && compare(v1a.Body, v2a.Body)
	}
	if v1a, ok := v1.(IfStmt); ok {
		v2a, ok := v2.(IfStmt)
		if !ok {
			return false
		}
		return compare(v1a.Condition, v2a.Condition) && compare(v1a.Body, v2a.Body) && compare(v1a.Else, v2a.Else)
	}
	if v1a, ok := v1.(EnumTypeDefn); ok {
		v2a, ok := v2.(EnumTypeDefn)
		if !ok {
			return false
		}
		if v1a.Name != v2a.Name {
			return false
		}
		if len(v1a.Options) != len(v2a.Options) {
			return false
		}
		for i := range v1a.Options {
			if !compare(v1a.Options[i], v2a.Options[i]) {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(TypeDefn); ok {
		v2a, ok := v2.(TypeDefn)
		if !ok {
			return false
		}
		if v1a.Name != v2a.Name {
			return false
		}
		if len(v1a.Parameters) != len(v2a.Parameters) {
			return false
		}
		for i := range v1a.Parameters {
			if v1a.Parameters[i] != v2a.Parameters[i] {
				return false
			}
		}
		return compare(v1a.ConcreteType, v2a.ConcreteType)
	}
	if v1a, ok := v1.(MatchStmt); ok {
		v2a, ok := v2.(MatchStmt)
		if !ok {
			return false
		}
		if !compare(v1a.Condition, v2a.Condition) {
			return false
		}
		if len(v1a.Cases) != len(v2a.Cases) {
			return false
		}
		for i := range v1a.Cases {
			if !compare(v1a.Cases[i], v2a.Cases[i]) {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(MatchCase); ok {
		v2a, ok := v2.(MatchCase)
		if !ok {
			return false
		}
		if len(v1a.LocalVariables) != len(v2a.LocalVariables) {
			return false
		}
		for i := range v1a.LocalVariables {
			if !compare(v1a.LocalVariables[i], v2a.LocalVariables[i]) {
				return false
			}
		}
		return compare(v1a.Variable, v2a.Variable) && compare(v1a.Body, v2a.Body)
	}
	if v1a, ok := v1.(EnumOption); ok {
		v2a, ok := v2.(EnumOption)
		if !ok {
			return false
		}
		if v1a.Constructor != v2a.Constructor {
			return false
		}
		if v1a.ParentType != v2a.ParentType {
			return false
		}
		if len(v1a.Parameters) != len(v2a.Parameters) {
			return false
		}
		for i := range v1a.Parameters {
			if v1a.Parameters[i] != v2a.Parameters[i] {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(EnumValue); ok {
		v2a, ok := v2.(EnumValue)
		if !ok {
			return false
		}
		if !compare(v1a.Constructor, v2a.Constructor) {
			return false
		}
		if len(v1a.Parameters) != len(v2a.Parameters) {
			return false
		}
		for i := range v1a.Parameters {
			if !compare(v1a.Parameters[i], v2a.Parameters[i]) {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(ArrayType); ok {
		v2a, ok := v2.(ArrayType)
		if !ok {
			return false
		}
		if !compare(v1a.Base, v2a.Base) {
			return false
		}
		return v1a.Size == v2a.Size
	}
	if v1a, ok := v1.(ArrayLiteral); ok {
		v2a, ok := v2.(ArrayLiteral)
		if !ok {
			return false
		}
		if len(v1a.Values) != len(v2a.Values) {
			return false
		}
		for i := range v1a.Values {
			if !compare(v1a.Values[i], v2a.Values[i]) {
				return false
			}
		}
		return true
	}
	if v1a, ok := v1.(ArrayValue); ok {
		v2a, ok := v2.(ArrayValue)
		if !ok {
			return false
		}
		if !compare(v1a.Base, v2a.Base) {
			return false
		}
		if !compare(v1a.Index, v2a.Index) {
			return false
		}
		return true
	}
	if v1a, ok := v1.(Cast); ok {
		v2a, ok := v2.(Cast)
		if !ok {
			return false
		}
		if !compare(v1a.Val, v2a.Val) {
			return false
		}
		if !compare(v1a.Typ, v2a.Typ) {
			return false
		}
		return true
	}
	if v1a, ok := v1.(Assertion); ok {
		v2a, ok := v2.(Assertion)
		if !ok {
			return false
		}
		if !compare(v1a.Predicate, v2a.Predicate) {
			return false
		}
		return v1a.Message == v2a.Message
	}
	if v1a, ok := v1.(SumType); ok {
		v2a, ok := v2.(SumType)
		if !ok {
			return false
		}
		if len(v1a) != len(v2a) {
			return false
		}
		for i := range v1a {
			if !compare(v1a[i], v2a[i]) {
				return false
			}
		}
		return true
	}

	if v1a, ok := v1.(TupleType); ok {
		if v2a, ok := v2.(TupleType); ok {
			if len(v1a) != len(v2a) {
				return false
			}
			for i := range v1a {
				if !compare(v1a[i], v2a[i]) {
					return false
				}
			}
			return true
		}
		return false
	}
	if v1a, ok := v1.(TupleValue); ok {
		if v2a, ok := v2.(TupleValue); ok {
			if len(v1a) != len(v2a) {
				return false
			}
			for i := range v1a {
				if !compare(v1a[i], v2a[i]) {
					return false
				}
			}
			return true
		}
		return false
	}
	if v1a, ok := v1.(UserType); ok {
		if v2a, ok := v2.(UserType); ok {
			if v1a.Name != v2a.Name {
				return false
			}
			return compare(v1a.Typ, v2a.Typ)
		}
		return false
	}
	if v1a, ok := v1.(Slice); ok {
		if v2a, ok := v2.(Slice); ok {
			if v1a.Size != v2a.Size {
				return false
			}
			return compare(v1a.Base, v2a.Base)
		}
		return false
	}
	panic(fmt.Sprintf("Unimplemented type for compare %v vs %v", reflect.TypeOf(v1), reflect.TypeOf(v2)))
}

func buildAst(t *testing.T, filename string) ([]Node, TypeInformation, Callables) {
	t.Helper()
	f, err := os.Open("../../testsuite/" + filename + ".l")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tokens, err := token.Tokenize(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	ast, ti, c, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	return ast, ti, c
}

func TestParseFizzBuzz(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.Fizzbuzz))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{Name: "terminate", Typ: TypeLiteral("bool")},
						InitialValue: BoolLiteral(false),
					},
					MutStmt{
						Var:          VarWithType{Name: "i", Typ: TypeLiteral("int")},
						InitialValue: IntLiteral(1),
					},
					WhileLoop{
						Condition: NotEqualsComparison{
							Left:  VarWithType{"terminate", TypeLiteral("bool"), false},
							Right: BoolLiteral(true),
						},
						Body: BlockStmt{
							[]Node{
								IfStmt{
									Condition: EqualityComparison{
										Left: ModOperator{
											Left:  VarWithType{Variable("i"), TypeLiteral("int"), false},
											Right: IntLiteral(15),
										},
										Right: IntLiteral(0),
									},
									Body: BlockStmt{
										[]Node{
											FuncCall{
												Name: "PrintString",
												UserArgs: []Value{
													StringLiteral(`fizzbuzz`),
												},
											},
										},
									},
									Else: BlockStmt{
										[]Node{
											IfStmt{
												Condition: EqualityComparison{
													Left: ModOperator{
														Left:  VarWithType{Variable("i"), TypeLiteral("int"), false},
														Right: IntLiteral(5),
													},
													Right: IntLiteral(0),
												},
												Body: BlockStmt{
													[]Node{
														FuncCall{
															Name: "PrintString",
															UserArgs: []Value{
																StringLiteral(`buzz`),
															},
														},
													},
												},
												Else: BlockStmt{
													[]Node{
														IfStmt{
															Condition: EqualityComparison{
																Left: ModOperator{
																	Left:  VarWithType{"i", TypeLiteral("int"), false},
																	Right: IntLiteral(3),
																},
																Right: IntLiteral(0),
															},
															Body: BlockStmt{
																[]Node{
																	FuncCall{
																		Name: "PrintString",
																		UserArgs: []Value{
																			StringLiteral(`fizz`),
																		},
																	},
																},
															},
															Else: BlockStmt{
																[]Node{
																	FuncCall{
																		Name: "PrintInt",
																		UserArgs: []Value{
																			VarWithType{"i", TypeLiteral("int"), false},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"i", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
								IfStmt{
									Condition: GreaterOrEqualComparison{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(100),
									},
									Body: BlockStmt{
										[]Node{
											AssignmentOperator{
												Variable: VarWithType{"terminate", TypeLiteral("bool"), false},
												Value:    BoolLiteral(true),
											},
										},
									},
								},
							}},
					},
				},
			},
		},
	}
	if len(expected) != len(ast) {
		t.Errorf("Unexpected number of nodes returned. got %v want %v", len(ast), len(expected))
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestHelloWorld(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.HelloWorld))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`Hello, world!\n`),
						},
					},
				},
			},
		},
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("hello, world: got %v want %v", ast[i], v)
		}
	}
}

func TestEmptyMain(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.EmptyMain))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{},
		},
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("empty function: got %v want %v", ast[i], v)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.LetStatement))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n", TypeLiteral("int"), false},
						Val: IntLiteral(5),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"n", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestLetStatementShadow(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.LetStatementShadow))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n", TypeLiteral("int"), false},
						Val: IntLiteral(5),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"n", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},

					LetStmt{
						Var: VarWithType{"n", TypeLiteral("string"), false},
						Val: StringLiteral("hello"),
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							VarWithType{"n", TypeLiteral("string"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestMutStatement(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.MutAddition))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"x", TypeLiteral("int"), false},
						InitialValue: IntLiteral(3),
					},
					MutStmt{
						Var: VarWithType{"y", TypeLiteral("int"), false},
						InitialValue: AdditionOperator{
							Left:  VarWithType{"x", TypeLiteral("int"), false},
							Right: IntLiteral(1),
						},
					},
					AssignmentOperator{
						Variable: VarWithType{"x", TypeLiteral("int"), false},
						Value: AdditionOperator{
							Left: VarWithType{"x", TypeLiteral("int"), false},
							Right: AdditionOperator{
								Left:  VarWithType{"y", TypeLiteral("int"), false},
								Right: IntLiteral(1),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"x", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("mut statement: got %v want %v", ast[i], v)
		}
	}
}

func TestTwoProcs(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.TwoProcs))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name: "foo",
			Args: nil,
			Return: []VarWithType{
				VarWithType{"", TypeLiteral("int"), false},
			},
			Body: BlockStmt{
				[]Node{
					ReturnStmt{
						Val: IntLiteral(3),
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "foo",
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("TwoProcs (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestOutOfOrder(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.OutOfOrder))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "foo",
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name: "foo",
			Args: nil,
			Return: []VarWithType{
				VarWithType{Name: "", Typ: TypeLiteral("int")},
			},
			Body: BlockStmt{
				[]Node{
					ReturnStmt{
						Val: IntLiteral(3),
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("TwoProcs (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestSumToTen(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SumToTen))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name: "sum",
			Args: []VarWithType{
				{Name: "x", Typ: TypeLiteral("int")},
			},
			Return: []VarWithType{
				{Typ: TypeLiteral("int")},
			},
			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"val", TypeLiteral("int"), false},
						InitialValue: VarWithType{"x", TypeLiteral("int"), false},
					},
					MutStmt{
						Var:          VarWithType{"sum", TypeLiteral("int"), false},
						InitialValue: IntLiteral(0),
					},
					WhileLoop{
						Condition: GreaterComparison{
							Left:  VarWithType{"val", TypeLiteral("int"), false},
							Right: IntLiteral(0),
						},
						Body: BlockStmt{
							[]Node{
								AssignmentOperator{
									Variable: VarWithType{"sum", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"sum", TypeLiteral("int"), false},
										Right: VarWithType{"val", TypeLiteral("int"), false},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"val", TypeLiteral("int"), false},
									Value: SubtractionOperator{
										Left:  VarWithType{"val", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
					ReturnStmt{VarWithType{"sum", TypeLiteral("int"), false}},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "sum",
								UserArgs: []Value{
									IntLiteral(10),
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("sum to ten (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSimpleFunc(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SimpleFunc))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name: "foo",
			Args: nil,
			Return: []VarWithType{
				VarWithType{Name: "", Typ: TypeLiteral("int")},
			},
			Body: BlockStmt{
				[]Node{
					ReturnStmt{
						Val: IntLiteral(3),
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "foo",
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("SimpleFunc (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestSumToTenRecursive(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SumToTenRecursive))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name: "sum",
			Args: []VarWithType{
				{Name: "x", Typ: TypeLiteral("int")},
			},
			Return: []VarWithType{
				{Typ: TypeLiteral("int")},
			},
			Body: BlockStmt{
				[]Node{
					ReturnStmt{
						Val: FuncCall{
							Name: "partial_sum",
							UserArgs: []Value{
								IntLiteral(0),
								VarWithType{"x", TypeLiteral("int"), false},
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name: "partial_sum",
			Args: []VarWithType{
				{"partial", TypeLiteral("int"), false},
				{"x", TypeLiteral("int"), false},
			},
			Return: []VarWithType{
				{Typ: TypeLiteral("int")},
			},
			Body: BlockStmt{
				[]Node{
					IfStmt{
						Condition: EqualityComparison{
							Left:  VarWithType{"x", TypeLiteral("int"), false},
							Right: IntLiteral(0),
						},
						Body: BlockStmt{
							[]Node{
								ReturnStmt{
									Val: VarWithType{"partial", TypeLiteral("int"), false},
								},
							},
						},
					},
					ReturnStmt{
						Val: FuncCall{
							Name: "partial_sum",
							UserArgs: []Value{
								AdditionOperator{
									Left:  VarWithType{"partial", TypeLiteral("int"), false},
									Right: VarWithType{"x", TypeLiteral("int"), false},
								},
								SubtractionOperator{
									Left:  VarWithType{"x", TypeLiteral("int"), false},
									Right: IntLiteral(1),
								},
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "sum",
								UserArgs: []Value{
									IntLiteral(10),
								},
							},
						},
					},
				},
			},
		},
	}
	if len(ast) != len(expected) {
		t.Fatalf("sum to ten recursive: incorrect number of nodes in AST. got %v want %v", len(ast), len(expected))
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("sum to ten recursive(%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestSomeMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SomeMath))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{Name: "add", Typ: TypeLiteral("int")},
						Val: AdditionOperator{
							Left:  IntLiteral(1),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "sub", Typ: TypeLiteral("int")},
						Val: SubtractionOperator{
							Left:  IntLiteral(1),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "mul", Typ: TypeLiteral("int")},
						Val: MulOperator{
							Left:  IntLiteral(2),
							Right: IntLiteral(3),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "div", Typ: TypeLiteral("int")},
						Val: DivOperator{
							Left:  IntLiteral(6),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "x", Typ: TypeLiteral("int")},
						Val: AdditionOperator{
							Left: IntLiteral(1),
							Right: SubtractionOperator{
								Left: MulOperator{
									Left:  IntLiteral(2),
									Right: IntLiteral(3),
								},
								Right: DivOperator{
									Left:  IntLiteral(4),
									Right: IntLiteral(2),
								},
							},
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`Add: `),
						},
					},
					FuncCall{
						Name: "PrintInt",

						UserArgs: []Value{
							VarWithType{"add", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`Sub: `),
						},
					},
					FuncCall{
						Name: "PrintInt",

						UserArgs: []Value{
							VarWithType{"sub", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`Mul: `),
						},
					},
					FuncCall{
						Name: "PrintInt",

						UserArgs: []Value{
							VarWithType{"mul", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`Div: `),
						},
					},
					FuncCall{
						Name: "PrintInt",

						UserArgs: []Value{
							VarWithType{"div", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`Complex: `),
						},
					},
					FuncCall{
						Name: "PrintInt",

						UserArgs: []Value{
							VarWithType{"x", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",

						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestEqualComparisonMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.EqualComparison))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(3),
					},
					LetStmt{
						Var: VarWithType{"b", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					IfStmt{
						Condition: EqualityComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: EqualityComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"a", TypeLiteral("int"), false},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},

								AssignmentOperator{
									Variable: VarWithType{"a", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"a", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestNotEqualComparisonMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.NotEqualComparison))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(3),
					},
					LetStmt{
						Var: VarWithType{"b", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					IfStmt{
						Condition: NotEqualsComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: NotEqualsComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"a", TypeLiteral("int"), false},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"a", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestGreaterComparisonMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.GreaterComparison))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var: VarWithType{"b", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					IfStmt{
						Condition: GreaterComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: GreaterComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"a", TypeLiteral("int"), false},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", TypeLiteral("int"), false},
									Value: SubtractionOperator{
										Left:  VarWithType{"a", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestGreaterOrEqualComparisonMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.GreaterOrEqualComparison))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var: VarWithType{"b", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					IfStmt{
						Condition: GreaterOrEqualComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: GreaterOrEqualComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"a", TypeLiteral("int"), false},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},

								AssignmentOperator{
									Variable: VarWithType{"a", TypeLiteral("int"), false},
									Value: SubtractionOperator{
										Left:  VarWithType{"a", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestLessThanComparisonMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.LessThanComparison))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var: VarWithType{"b", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					IfStmt{
						Condition: LessThanComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: LessThanComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"a", TypeLiteral("int"), false},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"a", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestLessThanOrEqualComparisonMath(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.LessThanOrEqualComparison))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(1),
					},
					LetStmt{
						Var: VarWithType{"b", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					IfStmt{
						Condition: LessThanOrEqualComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "PrintString",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: LessThanOrEqualComparison{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: VarWithType{"b", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"a", TypeLiteral("int"), false},
									},
								},
								FuncCall{
									Name: "PrintString",
									UserArgs: []Value{
										StringLiteral(`\n`),
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"a", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestUserDefinedType(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.UserDefinedType))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		TypeDefn{
			Name:         "Foo",
			ConcreteType: TypeLiteral("int"),
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"x", UserType{TypeLiteral("int"), "Foo"}, false},
						Val: IntLiteral(4),
					},
					FuncCall{
						Name:     "PrintInt",
						UserArgs: []Value{VarWithType{"x", UserType{TypeLiteral("int"), "Foo"}, false}},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("user defined type (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestTypeInference(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.TypeInference))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:   "foo",
			Args:   []VarWithType{{"x", TypeLiteral("int"), false}},
			Return: []VarWithType{{"", TypeLiteral("int"), false}},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: VarWithType{"x", TypeLiteral("int"), false},
					},
					AssignmentOperator{
						Variable: VarWithType{"a", TypeLiteral("int"), false},
						Value: AdditionOperator{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: IntLiteral(1),
						},
					},
					LetStmt{
						Var: VarWithType{"x", TypeLiteral("int"), false},
						Val: AdditionOperator{
							Left:  VarWithType{"a", TypeLiteral("int"), false},
							Right: IntLiteral(1),
						},
					},
					IfStmt{
						Condition: GreaterComparison{
							Left:  VarWithType{"x", TypeLiteral("int"), false},
							Right: IntLiteral(3),
						},
						Body: BlockStmt{
							[]Node{
								ReturnStmt{VarWithType{"a", TypeLiteral("int"), false}},
							},
						},
					},
					ReturnStmt{IntLiteral(0)},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name:     "foo",
								UserArgs: []Value{IntLiteral(1)},
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`, `),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name:     "foo",
								UserArgs: []Value{IntLiteral(3)},
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("type inference (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestEnumType(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.EnumType))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		EnumTypeDefn{
			"Foo",
			[]EnumOption{
				EnumOption{"A", nil, UserType{TypeLiteral("int64"), "Foo"}},
				EnumOption{"B", nil, UserType{TypeLiteral("int64"), "Foo"}},
			},
			nil,
			0,
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var: VarWithType{
							"a",
							EnumTypeDefn{
								"Foo",
								[]EnumOption{
									EnumOption{"A", nil, UserType{TypeLiteral("int64"), "Foo"}},
									EnumOption{"B", nil, UserType{TypeLiteral("int64"), "Foo"}},
								},
								nil,
								0,
							},
							false,
						},
						Val: EnumValue{Constructor: EnumOption{
							"A",
							nil,
							UserType{TypeLiteral("int64"), "Foo"},
						},
						},
					},
					MatchStmt{
						Condition: VarWithType{
							"a",
							EnumTypeDefn{
								"Foo",
								[]EnumOption{
									EnumOption{"A", nil, UserType{TypeLiteral("int64"), "Foo"}},
									EnumOption{"B", nil, UserType{TypeLiteral("int64"), "Foo"}},
								},
								nil,
								0,
							},
							false},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{
									"A",
									nil,
									UserType{TypeLiteral("int64"), "Foo"},
								},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`I am A!\n`),
											},
										},
									},
								},
							},
							MatchCase{
								Variable: EnumOption{
									"B",
									nil,
									UserType{TypeLiteral("int64"), "Foo"},
								},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`I am B!\n`),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestEnumTypeInferred(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.EnumTypeInferred))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		EnumTypeDefn{
			"Foo",
			[]EnumOption{
				EnumOption{"A", nil, UserType{TypeLiteral("int64"), "Foo"}},
				EnumOption{"B", nil, UserType{TypeLiteral("int64"), "Foo"}},
			},
			nil,
			0,
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var: VarWithType{"a", UserType{TypeLiteral("int64"), "Foo"}, false},
						Val: EnumValue{Constructor: EnumOption{"B", nil, UserType{TypeLiteral("int64"), "Foo"}}},
					},
					MatchStmt{
						Condition: VarWithType{"a", UserType{TypeLiteral("int64"), "Foo"}, false},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"A", nil, UserType{TypeLiteral("int64"), "Foo"}},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`I am A!\n`),
											},
										},
									},
								},
							},
							MatchCase{
								Variable: EnumOption{"B", nil, UserType{TypeLiteral("int64"), "Foo"}},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`I am B!\n`),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestIfElseMatch(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.IfElseMatch))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var: VarWithType{"x", TypeLiteral("int"), false},
						Val: IntLiteral(3),
					},
					MatchStmt{
						Condition: BoolLiteral(true),
						Cases: []MatchCase{
							MatchCase{
								Variable: LessThanComparison{
									Left:  VarWithType{"x", TypeLiteral("int"), false},
									Right: IntLiteral(3),
								},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`x is less than 3\n`),
											},
										},
									},
								},
							},
							MatchCase{
								Variable: GreaterComparison{
									Left:  VarWithType{"x", TypeLiteral("int"), false},
									Right: IntLiteral(3),
								},

								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`x is greater than 3\n`),
											},
										},
									},
								},
							},
							MatchCase{
								Variable: LessThanComparison{
									Left:  VarWithType{"x", TypeLiteral("int"), false},
									Right: IntLiteral(4),
								},

								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`x is less than 4\n`),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestGenericEnumType(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.GenericEnumType))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		EnumTypeDefn{
			"Maybe",
			[]EnumOption{
				EnumOption{"Nothing", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
				EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
			},
			nil,
			1,
		},
		FuncDecl{
			Name: "DoSomething",
			Args: []VarWithType{{"x", TypeLiteral("int"), false}},
			Return: []VarWithType{
				{
					"",
					EnumTypeDefn{
						"Maybe",
						[]EnumOption{
							EnumOption{"Nothing", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
							EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
						},
						[]Type{
							TypeLiteral("int"),
						},
						1,
					},
					false,
				},
			},
			Body: BlockStmt{
				[]Node{
					IfStmt{
						Condition: GreaterComparison{
							Left:  VarWithType{"x", TypeLiteral("int"), false},
							Right: IntLiteral(3),
						},
						Body: BlockStmt{
							[]Node{
								ReturnStmt{EnumValue{Constructor: EnumOption{"Nothing", nil, UserType{TypeLiteral("int64"), "Maybe"}}}},
							},
						},
					},
					ReturnStmt{EnumValue{
						Constructor: EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
						Parameters:  []Value{IntLiteral(5)},
					},
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							EnumTypeDefn{
								"Maybe",
								[]EnumOption{
									EnumOption{"Nothing", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
									EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
								},
								[]Type{
									TypeLiteral("int"),
								},
								1,
							},
							false,
						},
						Val: FuncCall{
							Name: "DoSomething",
							UserArgs: []Value{
								IntLiteral(3),
							},
						},
					},
					MatchStmt{
						Condition: VarWithType{
							"x",
							EnumTypeDefn{
								"Maybe",
								[]EnumOption{
									EnumOption{"Nothing", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
									EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
								},
								[]Type{
									TypeLiteral("int"),
								},
								1,
							},
							false,
						},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"Nothing", nil, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`I am nothing!\n`),
											},
										},
									},
								},
							},
							MatchCase{
								LocalVariables: []VarWithType{
									VarWithType{"n", TypeLiteral("int"), false},
								},
								Variable: EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintInt",
											UserArgs: []Value{
												VarWithType{"n", TypeLiteral("int"), false},
											},
										},
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`\n`),
											},
										},
									},
								},
							},
						},
					},

					LetStmt{
						Var: VarWithType{
							"x",
							EnumTypeDefn{
								"Maybe",
								[]EnumOption{
									EnumOption{"Nothing", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
									EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
								},
								[]Type{
									TypeLiteral("int"),
								},
								1,
							},
							false,
						},
						Val: FuncCall{
							Name: "DoSomething",
							UserArgs: []Value{
								IntLiteral(4),
							},
						},
					},
					MatchStmt{
						Condition: VarWithType{
							"x",
							EnumTypeDefn{
								"Maybe",
								[]EnumOption{
									EnumOption{"Nothing", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
									EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
								},
								[]Type{
									TypeLiteral("int"),
								},
								1,
							},
							false,
						},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"Nothing", nil, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`I am nothing!\n`),
											},
										},
									},
								},
							},
							MatchCase{
								LocalVariables: []VarWithType{
									VarWithType{"n", TypeLiteral("int"), false},
								},
								Variable: EnumOption{"Just", []string{"a"}, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										FuncCall{

											Name: "PrintInt",
											UserArgs: []Value{
												VarWithType{"n", TypeLiteral("int"), false},
											},
										},
										FuncCall{

											Name: "PrintString",
											UserArgs: []Value{
												StringLiteral(`\n`),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestMatchParam(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.MatchParam))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		EnumTypeDefn{
			"Maybe",
			[]EnumOption{
				EnumOption{"Nothing", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
				EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
			},
			nil,
			1,
		},
		FuncDecl{
			Name: "foo",
			Args: []VarWithType{{
				"x",
				EnumTypeDefn{
					"Maybe",
					[]EnumOption{
						EnumOption{"Nothing", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
						EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
					},
					[]Type{
						TypeLiteral("int"),
					},
					1,
				},
				false,
			},
			},
			Return: []VarWithType{
				{
					"",
					TypeLiteral("int"),
					false,
				},
			},
			Body: BlockStmt{
				[]Node{
					MatchStmt{
						Condition: VarWithType{
							"x",
							EnumTypeDefn{
								"Maybe",
								[]EnumOption{
									EnumOption{"Nothing", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
									EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
								},
								[]Type{
									TypeLiteral("int"),
								},
								1,
							},
							false,
						},
						Cases: []MatchCase{
							MatchCase{
								LocalVariables: []VarWithType{
									VarWithType{"n", TypeLiteral("int"), false},
								},
								Variable: EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										ReturnStmt{VarWithType{"n", TypeLiteral("int"), false}},
									},
								},
							},
							MatchCase{
								Variable: EnumOption{"Nothing", nil, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										ReturnStmt{IntLiteral(0)},
									},
								},
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "foo",
								UserArgs: []Value{
									EnumValue{
										Constructor: EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
										Parameters:  []Value{IntLiteral(5)},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestMatchParam2(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.MatchParam2))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		EnumTypeDefn{
			"Maybe",
			[]EnumOption{
				EnumOption{"Nothing", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
				EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
			},
			nil,
			1,
		},
		FuncDecl{
			Name: "foo",
			Args: []VarWithType{
				{
					"x",
					EnumTypeDefn{
						"Maybe",
						[]EnumOption{
							EnumOption{"Nothing", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
							EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
						},
						[]Type{
							TypeLiteral("int"),
						},
						1,
					},
					false,
				},
			},
			Return: []VarWithType{
				{
					"",
					TypeLiteral("int"),
					false,
				},
			},
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`x`),
						},
					},
					MatchStmt{
						Condition: VarWithType{
							"x",
							EnumTypeDefn{
								"Maybe",
								[]EnumOption{
									EnumOption{"Nothing", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
									EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
								},
								[]Type{
									TypeLiteral("int"),
								},
								1,
							},
							false,
						},
						Cases: []MatchCase{
							MatchCase{
								LocalVariables: []VarWithType{
									VarWithType{"n", TypeLiteral("int"), false},
								},
								Variable: EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										ReturnStmt{VarWithType{"n", TypeLiteral("int"), false}},
									},
								},
							},
							MatchCase{
								Variable: EnumOption{"Nothing", nil, UserType{TypeLiteral("int64"), "Maybe"}},
								Body: BlockStmt{
									[]Node{
										ReturnStmt{IntLiteral(0)},
									},
								},
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "foo",
								UserArgs: []Value{
									EnumValue{
										Constructor: EnumOption{"Just", []string{"x"}, UserType{TypeLiteral("int64"), "Maybe"}},
										Parameters:  []Value{IntLiteral(5)},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestSimpleAlgorithm(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SimpleAlgorithm))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name: "loop",
			Args: []VarWithType{{"high", TypeLiteral("int"), false}},
			Return: []VarWithType{
				{
					"",
					TypeLiteral("int"),
					false,
				},
			},
			Body: BlockStmt{
				[]Node{

					MutStmt{
						Var:          VarWithType{"total", TypeLiteral("int"), false},
						InitialValue: IntLiteral(0),
					},
					MutStmt{
						Var:          VarWithType{"i", TypeLiteral("int"), false},
						InitialValue: IntLiteral(0),
					},
					LetStmt{
						Var: VarWithType{"high", TypeLiteral("int"), false},
						Val: MulOperator{
							Left:  VarWithType{"high", TypeLiteral("int"), false},
							Right: IntLiteral(2),
						},
					},

					AssignmentOperator{
						Variable: VarWithType{"i", TypeLiteral("int"), false},
						Value:    IntLiteral(1),
					},
					WhileLoop{
						Condition: LessThanComparison{
							Left:  VarWithType{"i", TypeLiteral("int"), false},
							Right: VarWithType{"high", TypeLiteral("int"), false},
						},
						Body: BlockStmt{
							[]Node{
								IfStmt{
									Condition: EqualityComparison{
										Left: ModOperator{
											Left:  VarWithType{Variable("i"), TypeLiteral("int"), false},
											Right: IntLiteral(2),
										},
										Right: IntLiteral(0),
									},
									Body: BlockStmt{
										[]Node{
											AssignmentOperator{
												Variable: VarWithType{"total", TypeLiteral("int"), false},
												Value: AdditionOperator{
													Left: VarWithType{Variable("total"), TypeLiteral("int"), false},
													Right: MulOperator{
														Left:  VarWithType{Variable("i"), TypeLiteral("int"), false},
														Right: IntLiteral(2),
													},
												},
											},
										},
									},
								},

								AssignmentOperator{
									Variable: VarWithType{"i", TypeLiteral("int"), false},
									Value: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
					ReturnStmt{VarWithType{"total", TypeLiteral("int"), false}},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name:     "loop",
								UserArgs: []Value{IntLiteral(10)},
							},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestConcreteTypeInt64(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.ConcreteTypeInt64))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"x", TypeLiteral("int64"), false},
						Val: IntLiteral(-4),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{Variable("x"), TypeLiteral("int64"), false},
						},
					},
				},
			},
		},
	}

	if len(expected) != len(ast) {
		t.Fatalf("Unexpected AST: got %v want %v", ast, expected)
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("enum test (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestSimpleArray(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SimpleArray))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							ArrayType{
								Base: TypeLiteral("int"),
								Size: IntLiteral(5),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(5),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSimpleArrayInference(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SimpleArrayInference))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							ArrayType{
								Base: TypeLiteral("int"),
								Size: IntLiteral(5),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(5),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestArrayMutation(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.ArrayMutation))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var: VarWithType{"n",
							ArrayType{
								Base: TypeLiteral("int"),
								Size: IntLiteral(5),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(5),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					AssignmentOperator{
						Variable: ArrayValue{
							Base: VarWithType{"n",
								ArrayType{
									Base: TypeLiteral("int"),
									Size: IntLiteral(5),
								},
								false,
							},
							Index: IntLiteral(3),
						},
						Value: IntLiteral(2),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(5),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(5),
									},
									false,
								},
								Index: IntLiteral(2),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestReferenceVariable(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.ReferenceVariable))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name: "changer",
			Args: []VarWithType{
				VarWithType{"x", TypeLiteral("int"), true},
				VarWithType{"y", TypeLiteral("int"), false},
			},
			Return: []VarWithType{
				VarWithType{Name: "", Typ: TypeLiteral("int")},
			},
			Effects: []Effect{"mutate"},
			Body: BlockStmt{
				[]Node{
					AssignmentOperator{
						Variable: VarWithType{"x", TypeLiteral("int"), true},
						Value:    IntLiteral(4),
					},
					ReturnStmt{
						AdditionOperator{
							Left:  VarWithType{"x", TypeLiteral("int"), true},
							Right: VarWithType{"y", TypeLiteral("int"), false},
						},
					},
				},
			},
		},

		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"var", TypeLiteral("int"), false},
						InitialValue: IntLiteral(3),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"var", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					LetStmt{
						Var: VarWithType{"sum", TypeLiteral("int"), false},
						Val: FuncCall{
							Name: "changer",
							UserArgs: []Value{
								VarWithType{"var", TypeLiteral("int"), false},
								IntLiteral(3),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"var", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"sum", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSimpleSlice(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SimpleSlice))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSimpleSliceInference(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SimpleSliceInference))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					LetStmt{
						Var: VarWithType{"n2",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						Val: VarWithType{
							"n",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n2",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSliceMutation(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SliceMutation))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var: VarWithType{"n",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					AssignmentOperator{
						Variable: ArrayValue{
							Base: VarWithType{"n",
								SliceType{
									Base: TypeLiteral("int"),
								},
								false,
							},
							Index: IntLiteral(3),
						},
						Value: IntLiteral(2),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(3),
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`\n`),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"n",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(2),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSliceParam(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SliceParam))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"b",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(44),
								IntLiteral(55),
								IntLiteral(88),
							},
						},
					},
					FuncCall{
						Name: "PrintASlice",
						UserArgs: []Value{
							VarWithType{
								"b",
								SliceType{
									Base: TypeLiteral("byte"),
								},
								false,
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name: "PrintASlice",
			Args: []VarWithType{
				{Name: "A", Typ: SliceType{Base: TypeLiteral("byte")}},
			},
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintByteSlice",
						UserArgs: []Value{
							VarWithType{
								"A",
								SliceType{
									Base: TypeLiteral("byte"),
								},
								false,
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestReadSyscall(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.ReadSyscall))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO", "Filesystem"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"fd",
							TypeLiteral("uint64"),
							false,
						},
						Val: FuncCall{
							Name: "Open",
							UserArgs: []Value{
								StringLiteral("foo.txt"),
							},
						},
					},
					MutStmt{
						Var: VarWithType{
							"dta",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(0),
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					LetStmt{
						Var: VarWithType{
							"n",
							TypeLiteral("uint64"),
							false,
						},
						Val: FuncCall{
							Name: "Read",
							UserArgs: []Value{
								VarWithType{
									"fd",
									TypeLiteral("uint64"),
									false,
								},
								VarWithType{
									"dta",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
							},
						},
					},
					FuncCall{
						Name: "PrintByteSlice",
						UserArgs: []Value{
							VarWithType{
								"dta",
								SliceType{
									Base: TypeLiteral("byte"),
								},
								false,
							},
						},
					},
					FuncCall{
						Name: "Close",
						UserArgs: []Value{

							VarWithType{
								"fd",
								TypeLiteral("uint64"),
								false,
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestSliceLength(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SliceLength))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SliceType{TypeLiteral("string")},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								StringLiteral("3"),
								StringLiteral("foo"),
								StringLiteral("hello"),
								StringLiteral("world"),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "len",
								UserArgs: []Value{
									VarWithType{
										"x",
										SliceType{TypeLiteral("string")},
										false,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestIndexAssignment(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.IndexAssignment))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SliceType{TypeLiteral("int")},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					MutStmt{
						Var: VarWithType{
							"n",
							TypeLiteral("int"),
							false,
						},
						InitialValue: ArrayValue{
							Base: VarWithType{"x",
								SliceType{
									Base: TypeLiteral("int"),
								},
								false,
							},
							Index: IntLiteral(1),
						},
					},
					LetStmt{
						Var: VarWithType{
							"n2",
							TypeLiteral("int"),
							false,
						},
						Val: ArrayValue{
							Base: VarWithType{"x",
								SliceType{
									Base: TypeLiteral("int"),
								},
								false,
							},
							Index: IntLiteral(2),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{
								"n",
								TypeLiteral("int"),
								false,
							},
						},
					},
					FuncCall{
						Name:     "PrintString",
						UserArgs: []Value{StringLiteral(`\n`)},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{
								"n2",
								TypeLiteral("int"),
								false,
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}

}

func TestIndexedAddition(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.IndexedAddition))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SliceType{TypeLiteral("int")},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					MutStmt{
						Var: VarWithType{
							"n",
							TypeLiteral("int"),
							false,
						},
						InitialValue: ArrayValue{
							Base: VarWithType{"x",
								SliceType{
									Base: TypeLiteral("int"),
								},
								false,
							},
							Index: IntLiteral(1),
						},
					},
					AssignmentOperator{
						Variable: VarWithType{"n", TypeLiteral("int"), false},
						Value: AdditionOperator{
							Left: VarWithType{"n", TypeLiteral("int"), false},
							Right: ArrayValue{
								Base: VarWithType{"x",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(2),
							},
						},
					},
					LetStmt{
						Var: VarWithType{
							"n2",
							TypeLiteral("int"),
							false,
						},
						Val: AdditionOperator{
							Left: ArrayValue{
								Base: VarWithType{"x",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(2),
							},
							Right: ArrayValue{
								Base: VarWithType{"x",
									SliceType{
										Base: TypeLiteral("int"),
									},
									false,
								},
								Index: IntLiteral(0),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{
								"n",
								TypeLiteral("int"),
								false,
							},
						},
					},
					FuncCall{
						Name:     "PrintString",
						UserArgs: []Value{StringLiteral(`\n`)},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{
								"n2",
								TypeLiteral("int"),
								false,
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("let statement (%d): got %v want %v", i, ast[i], v)
		}
	}
}

func TestPrecedence(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.Precedence))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatalf("%v: %v", err, tokens)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							TypeLiteral("int"),
							false,
						},
						Val: MulOperator{
							Left: Brackets{
								AdditionOperator{
									Left:  IntLiteral(1),
									Right: IntLiteral(2),
								},
							},
							Right: Brackets{
								SubtractionOperator{
									Left:  IntLiteral(3),
									Right: IntLiteral(4),
								},
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"x", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestLetCondition(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.LetCondition))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"i",
							TypeLiteral("int"),
							false,
						},
						Val: IntLiteral(0),
					},
					IfStmt{
						Condition: EqualityComparison{
							Left: Brackets{
								LetStmt{
									Var: VarWithType{"i", TypeLiteral("int"), false},
									Val: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
							Right: IntLiteral(1),
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"i", TypeLiteral("int"), false},
									},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										IntLiteral(-1),
									},
								},
							},
						},
					},
					IfStmt{
						Condition: NotEqualsComparison{
							Left: Brackets{
								LetStmt{
									Var: VarWithType{"i", TypeLiteral("int"), false},
									Val: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
							Right: IntLiteral(1),
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"i", TypeLiteral("int"), false},
									},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										IntLiteral(-1),
									},
								},
							},
						},
					},
					WhileLoop{
						Condition: LessThanComparison{
							Left: Brackets{
								LetStmt{
									Var: VarWithType{"i", TypeLiteral("int"), false},
									Val: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},

							Right: IntLiteral(3),
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "PrintInt",
									UserArgs: []Value{
										VarWithType{"i", TypeLiteral("int"), false},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestUnbufferedCat2(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.UnbufferedCat2))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name: "main",
			Args: []VarWithType{
				VarWithType{
					"args",
					SliceType{TypeLiteral("string")},
					false,
				},
			},
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var: VarWithType{"buf",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(0),
							},
						},
					},
					LetStmt{
						Var: VarWithType{"i",
							TypeLiteral("int"),
							false,
						},
						Val: IntLiteral(0),
					},
					WhileLoop{
						Condition: LessThanComparison{
							Left: Brackets{
								LetStmt{
									Var: VarWithType{"i", TypeLiteral("int"), false},
									Val: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
							Right: FuncCall{
								Name: "len",
								UserArgs: []Value{
									VarWithType{
										"args",
										SliceType{
											Base: TypeLiteral("string"),
										},
										false,
									},
								},
							},
						},
						Body: BlockStmt{
							[]Node{
								LetStmt{
									Var: VarWithType{"file", TypeLiteral("uint64"), false},
									Val: FuncCall{
										Name: "Open",
										UserArgs: []Value{
											ArrayValue{
												Base: VarWithType{
													"args",
													SliceType{
														Base: TypeLiteral("string"),
													},
													false,
												},
												Index: VarWithType{"i", TypeLiteral("int"), false},
											},
										},
									},
								},
								WhileLoop{
									Condition: GreaterComparison{
										Left: Brackets{
											LetStmt{
												Var: VarWithType{"n", TypeLiteral("uint64"), false},
												Val: FuncCall{
													Name: "Read",
													UserArgs: []Value{
														VarWithType{
															"file",
															TypeLiteral("uint64"),
															false,
														},
														VarWithType{
															"buf",
															SliceType{
																Base: TypeLiteral("byte"),
															},
															false,
														},
													},
												},
											},
										},
										Right: IntLiteral(0),
									},
									Body: BlockStmt{
										[]Node{
											FuncCall{
												Name: "PrintByteSlice",
												UserArgs: []Value{
													VarWithType{
														"buf",
														SliceType{
															Base: TypeLiteral("byte"),
														},
														false,
													},
												},
											},
										},
									},
								},
								FuncCall{
									Name: "Close",
									UserArgs: []Value{
										VarWithType{
											"file",
											TypeLiteral("uint64"),
											false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestMethodSyntax(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.MethodSyntax))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"foo",
							TypeLiteral("int"),
							false,
						},
						Val: IntLiteral(3),
					},
					LetStmt{
						Var: VarWithType{"y",
							TypeLiteral("int"),
							false,
						},
						Val: FuncCall{
							Name: "add",
							UserArgs: []Value{
								FuncCall{
									Name: "add3",
									UserArgs: []Value{
										VarWithType{"foo",
											TypeLiteral("int"),
											false,
										},
									},
								},
								IntLiteral(4),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{
								"y",
								TypeLiteral("int"),
								false,
							},
						},
					},
				},
			},
		},
	}
	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestUnbufferedCat3(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.UnbufferedCat3))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name: "main",
			Args: []VarWithType{
				VarWithType{
					"args",
					SliceType{TypeLiteral("string")},
					false,
				},
			},
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var: VarWithType{"buf",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(0),
							},
						},
					},
					LetStmt{
						Var: VarWithType{"i",
							TypeLiteral("int"),
							false,
						},
						Val: IntLiteral(0),
					},
					WhileLoop{
						Condition: LessThanComparison{
							Left: Brackets{
								LetStmt{
									Var: VarWithType{"i", TypeLiteral("int"), false},
									Val: AdditionOperator{
										Left:  VarWithType{"i", TypeLiteral("int"), false},
										Right: IntLiteral(1),
									},
								},
							},
							Right: FuncCall{
								Name: "len",
								UserArgs: []Value{
									VarWithType{
										"args",
										SliceType{
											Base: TypeLiteral("string"),
										},
										false,
									},
								},
							},
						},
						Body: BlockStmt{
							[]Node{
								LetStmt{
									Var: VarWithType{"file", TypeLiteral("uint64"), false},
									Val: FuncCall{
										Name: "Open",
										UserArgs: []Value{
											ArrayValue{
												Base: VarWithType{
													"args",
													SliceType{
														Base: TypeLiteral("string"),
													},
													false,
												},
												Index: VarWithType{"i", TypeLiteral("int"), false},
											},
										},
									},
								},
								WhileLoop{
									Condition: GreaterComparison{
										Left: Brackets{
											LetStmt{
												Var: VarWithType{"n", TypeLiteral("uint64"), false},
												Val: FuncCall{
													Name: "Read",
													UserArgs: []Value{
														VarWithType{
															"file",
															TypeLiteral("uint64"),
															false,
														},
														VarWithType{
															"buf",
															SliceType{
																Base: TypeLiteral("byte"),
															},
															false,
														},
													},
												},
											},
										},
										Right: IntLiteral(0),
									},
									Body: BlockStmt{
										[]Node{
											FuncCall{
												Name: "PrintByteSlice",
												UserArgs: []Value{
													VarWithType{
														"buf",
														SliceType{
															Base: TypeLiteral("byte"),
														},
														false,
													},
												},
											},
										},
									},
								},
								FuncCall{
									Name: "Close",
									UserArgs: []Value{
										VarWithType{
											"file",
											TypeLiteral("uint64"),
											false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestAssignmentToVariableIndex(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.AssignmentToVariableIndex))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{

				[]Node{
					MutStmt{
						Var: VarWithType{"x",
							ArrayType{
								Base: TypeLiteral("int"),
								Size: IntLiteral(4),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					LetStmt{
						Var: VarWithType{"y", TypeLiteral("int"), false},
						Val: ArrayValue{
							Base: VarWithType{"x",
								ArrayType{
									Base: TypeLiteral("int"),
									Size: IntLiteral(4),
								},
								false,
							},
							Index: IntLiteral(0),
						},
					},
					AssignmentOperator{
						Variable: ArrayValue{
							Base: VarWithType{"x",
								ArrayType{
									Base: TypeLiteral("int"),
									Size: IntLiteral(4),
								},
								false,
							},
							Index: VarWithType{"y", TypeLiteral("int"), false},
						},
						Value: IntLiteral(6),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"x",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(4),
									},
									false,
								},
								Index: VarWithType{"y", TypeLiteral("int"), false},
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"x",
									ArrayType{
										Base: TypeLiteral("int"),
										Size: IntLiteral(4),
									},
									false,
								},
								Index: AdditionOperator{
									Left:  VarWithType{"y", TypeLiteral("int"), false},
									Right: IntLiteral(1),
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestAssignmentToSliceVariableIndex(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.AssignmentToSliceVariableIndex))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{

				[]Node{
					MutStmt{
						Var: VarWithType{"x",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					LetStmt{
						Var: VarWithType{"y", TypeLiteral("byte"), false},
						Val: ArrayValue{
							Base: VarWithType{"x",
								SliceType{
									Base: TypeLiteral("byte"),
								},
								false,
							},
							Index: IntLiteral(0),
						},
					},
					AssignmentOperator{
						Variable: ArrayValue{
							Base: VarWithType{"x",
								SliceType{
									Base: TypeLiteral("byte"),
								},
								false,
							},
							Index: VarWithType{"y", TypeLiteral("byte"), false},
						},
						Value: IntLiteral(6),
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"x",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Index: VarWithType{"y", TypeLiteral("byte"), false},
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"x",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Index: AdditionOperator{
									Left:  VarWithType{"y", TypeLiteral("byte"), false},
									Right: IntLiteral(1),
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestCastBuiltin(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.CastBuiltin))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: []Effect{"IO"},

			Body: BlockStmt{

				[]Node{
					LetStmt{
						Var: VarWithType{
							"foo",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(70),
								IntLiteral(111),
								IntLiteral(111),
							},
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							Cast{
								Val: VarWithType{
									"foo",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Typ: TypeLiteral("string"),
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestAssertionFail(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.AssertionFail))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					Assertion{
						Predicate: BoolLiteral(false),
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestAssertionPass(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.AssertionPass))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					Assertion{
						Predicate: BoolLiteral(true),
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestAssertionFailWithMessage(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.AssertionFailWithMessage))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					Assertion{
						Predicate: BoolLiteral(false),
						Message:   "This always fails",
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestAssertionPassWithMessage(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.AssertionPassWithMessage))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					Assertion{
						Predicate: BoolLiteral(true),
						Message:   "You should never see this",
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestSumTypeDefn(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SumTypeDefn))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		TypeDefn{
			Name: "Foo",
			ConcreteType: SumType{
				TypeLiteral("int"),
				TypeLiteral("string"),
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestSumTypeFuncCall(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SumTypeFuncCall))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name: "foo",
			Args: []VarWithType{
				VarWithType{Name: "x", Typ: SumType{TypeLiteral("int"), TypeLiteral("string")}},
			},
			Return:  nil,
			Effects: []Effect{"IO"},
			Body: BlockStmt{
				[]Node{
					MatchStmt{
						Condition: VarWithType{Name: "x", Typ: SumType{TypeLiteral("int"), TypeLiteral("string")}},
						Cases: []MatchCase{
							MatchCase{
								Variable: VarWithType{Name: "x", Typ: TypeLiteral("int")},
								Body: BlockStmt{
									[]Node{

										FuncCall{
											Name: "PrintInt",
											UserArgs: []Value{
												VarWithType{"x", TypeLiteral("int"), false},
											},
										},
									},
								},
							},
							MatchCase{
								Variable: VarWithType{Name: "x", Typ: TypeLiteral("string")},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												VarWithType{"x", TypeLiteral("string"), false},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "foo",
						UserArgs: []Value{
							StringLiteral("bar"),
						},
					},
					FuncCall{
						Name: "foo",
						UserArgs: []Value{
							IntLiteral(3),
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("got %v want %v", ast[i], v)
		}
	}
}

func TestSumTypeFuncReturn(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.SumTypeFuncReturn))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name: "foo",
			Args: []VarWithType{
				VarWithType{Name: "x", Typ: TypeLiteral("bool")},
			},
			Return: []VarWithType{
				VarWithType{Name: "", Typ: SumType{TypeLiteral("int"), TypeLiteral("string")}},
			},
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					IfStmt{
						Condition: VarWithType{"x", TypeLiteral("bool"), false},
						Body: BlockStmt{
							[]Node{
								ReturnStmt{
									Val: IntLiteral(3),
								},
							},
						},
					},
					ReturnStmt{
						Val: StringLiteral("not3"),
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SumType{
								TypeLiteral("int"),
								TypeLiteral("string"),
							},
							false,
						},
						Val: FuncCall{
							Name:     "foo",
							UserArgs: []Value{BoolLiteral(false)},
						},
					},
					MatchStmt{
						Condition: VarWithType{Name: "x", Typ: SumType{TypeLiteral("int"), TypeLiteral("string")}},

						Cases: []MatchCase{
							MatchCase{
								Variable: VarWithType{Name: "x", Typ: TypeLiteral("int")},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintInt",
											UserArgs: []Value{
												VarWithType{"x", TypeLiteral("int"), false},
											},
										},
									},
								},
							},
							MatchCase{
								Variable: VarWithType{Name: "x", Typ: TypeLiteral("string")},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												VarWithType{"x", TypeLiteral("string"), false},
											},
										},
									},
								},
							},
						},
					},
					LetStmt{
						Var: VarWithType{
							"x",
							SumType{
								TypeLiteral("int"),
								TypeLiteral("string"),
							},
							false,
						},
						Val: FuncCall{
							Name:     "foo",
							UserArgs: []Value{BoolLiteral(true)},
						},
					},
					MatchStmt{
						Condition: VarWithType{Name: "x", Typ: SumType{TypeLiteral("int"), TypeLiteral("string")}},

						Cases: []MatchCase{
							MatchCase{
								Variable: VarWithType{Name: "x", Typ: TypeLiteral("int")},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintInt",
											UserArgs: []Value{
												VarWithType{"x", TypeLiteral("int"), false},
											},
										},
									},
								},
							},
							MatchCase{
								Variable: VarWithType{Name: "x", Typ: TypeLiteral("string")},
								Body: BlockStmt{
									[]Node{
										FuncCall{
											Name: "PrintString",
											UserArgs: []Value{
												VarWithType{"x", TypeLiteral("string"), false},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v", i, ast[i], v)
		}
	}
}

func TestLineComment(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.LineComment))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							TypeLiteral("int"),
							false,
						},
						Val: IntLiteral(3),
					},

					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"x", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v", i, ast[i], v)
		}
	}
}

func TestBlockComment(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.BlockComment))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							TypeLiteral("int"),
							false,
						},
						Val: IntLiteral(3),
					},

					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							VarWithType{"x", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v", i, ast[i], v)
		}
	}
}

func TestProductTypeValue(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.ProductTypeValue))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							TupleType{
								VarWithType{"x", TypeLiteral("int"), false},
								VarWithType{"y", TypeLiteral("bool"), false},
							},
							false,
						},
						Val: TupleValue{
							IntLiteral(3),
							BoolLiteral(false),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							// FIXME: Decide what this should really be.
							// Special case of VarWithType or a new node type?
							VarWithType{"x.x", TypeLiteral("int"), false},
						},
					},
					FuncCall{
						Name:     "PrintString",
						UserArgs: []Value{StringLiteral(`\n`)},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							// FIXME: Decide what this should really be.
							// Special case of VarWithType or a new node type?
							VarWithType{"x.y", TypeLiteral("bool"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v", i, ast[i], v)
		}
	}
}

func TestUserProductTypeValue(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.UserProductTypeValue))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		TypeDefn{
			Name: "Foo",
			ConcreteType: TupleType{
				VarWithType{"x", TypeLiteral("int"), false},
				VarWithType{"y", TypeLiteral("string"), false},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							UserType{
								TupleType{
									VarWithType{"x", TypeLiteral("int"), false},
									VarWithType{"y", TypeLiteral("string"), false},
								},
								"Foo",
							},
							false,
						},
						Val: TupleValue{
							IntLiteral(3),
							StringLiteral(`hello\n`),
						},
					},
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							VarWithType{"x.y", TypeLiteral("string"), false},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							// FIXME: Decide what this should really be.
							// Special case of VarWithType or a new node type?
							VarWithType{"x.x", TypeLiteral("int"), false},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v", i, ast[i], v)
		}
	}
}

func TestUserSumTypeDefn(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.UserSumTypeDefn))
	if err != nil {
		t.Fatal(err)
	}
	ast, _, _, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		EnumTypeDefn{
			"Keyword",
			[]EnumOption{
				EnumOption{"While", nil, UserType{TypeLiteral("int64"), "Keyword"}},
				EnumOption{"Mutable", nil, UserType{TypeLiteral("int64"), "Keyword"}},
			},
			nil,
			0,
		},
		TypeDefn{
			Name: "Token",
			ConcreteType: SumType{
				EnumTypeDefn{
					"Keyword",
					[]EnumOption{
						EnumOption{"While", nil, UserType{TypeLiteral("int64"), "Keyword"}},
						EnumOption{"Mutable", nil, UserType{TypeLiteral("int64"), "Keyword"}},
					},
					nil,
					0,
				},
				TypeLiteral("string"),
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v (%v, %v)", i, ast[i], v, reflect.TypeOf(ast[i]), reflect.TypeOf(v))
		}
	}
}

func TestSliceFromArray(t *testing.T) {
	ast, _, _ := buildAst(t, "slicefromarray")

	expected := []Node{
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"x",
							ArrayType{
								Base: TypeLiteral("byte"),
								Size: IntLiteral(5),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					LetStmt{
						Var: VarWithType{"y",
							SliceType{
								Base: TypeLiteral("byte"),
							},
							false,
						},
						Val: Slice{
							Base: ArrayValue{
								Base: VarWithType{"x",
									ArrayType{
										Base: TypeLiteral("byte"),
										Size: IntLiteral(5),
									},
									false,
								},
								Index: IntLiteral(2),
							},
							Size: IntLiteral(2),
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"y",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Index: IntLiteral(0),
							},
						},
					},
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"y",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Index: IntLiteral(1),
							},
						},
					},

					Assertion{
						Predicate: EqualityComparison{
							Left: FuncCall{
								Name: "len",

								UserArgs: []Value{
									VarWithType{"y",
										SliceType{
											Base: TypeLiteral("byte"),
										},
										false,
									},
								},
							},
							Right: IntLiteral(2),
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ArrayValue{
								Base: VarWithType{"y",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Index: IntLiteral(0),
							},
							Right: IntLiteral(3),
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ArrayValue{
								Base: VarWithType{"y",
									SliceType{
										Base: TypeLiteral("byte"),
									},
									false,
								},
								Index: IntLiteral(1),
							},
							Right: IntLiteral(4),
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v (%v, %v)", i, ast[i], v, reflect.TypeOf(ast[i]), reflect.TypeOf(v))
		}
	}
}

func TestArrayArg(t *testing.T) {
	ast, _, _ := buildAst(t, "arrayparam")

	expected := []Node{
		FuncDecl{
			Name: "foo",
			Args: []VarWithType{
				VarWithType{"x", ArrayType{Base: TypeLiteral("byte"), Size: 5}, false},
			},

			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							ArrayValue{
								Base: VarWithType{"x",
									ArrayType{
										Base: TypeLiteral("byte"),
										Size: 5,
									},
									false,
								},
								Index: IntLiteral(0),
							},
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ModOperator{
								Left: ArrayValue{
									Base: VarWithType{"x",
										ArrayType{
											Base: TypeLiteral("byte"),
											Size: 5,
										},
										false,
									},
									Index: IntLiteral(0),
								},
								Right: IntLiteral(5),
							},
							Right: IntLiteral(1),
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ModOperator{
								Left: ArrayValue{
									Base: VarWithType{"x",
										ArrayType{
											Base: TypeLiteral("byte"),
											Size: 5,
										},
										false,
									},
									Index: IntLiteral(1),
								},
								Right: IntLiteral(5),
							},
							Right: IntLiteral(2),
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ModOperator{
								Left: ArrayValue{
									Base: VarWithType{"x",
										ArrayType{
											Base: TypeLiteral("byte"),
											Size: 5,
										},
										false,
									},
									Index: IntLiteral(2),
								},
								Right: IntLiteral(5),
							},
							Right: IntLiteral(3),
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ModOperator{
								Left: ArrayValue{
									Base: VarWithType{"x",
										ArrayType{
											Base: TypeLiteral("byte"),
											Size: 5,
										},
										false,
									},
									Index: IntLiteral(3),
								},
								Right: IntLiteral(5),
							},
							Right: IntLiteral(4),
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: ModOperator{
								Left: ArrayValue{
									Base: VarWithType{"x",
										ArrayType{
											Base: TypeLiteral("byte"),
											Size: 5,
										},
										false,
									},
									Index: IntLiteral(4),
								},
								Right: IntLiteral(5),
							},
							Right: IntLiteral(5),
						},
					},
				},
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							ArrayType{
								Base: TypeLiteral("byte"),
								Size: IntLiteral(5),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								IntLiteral(1),
								IntLiteral(2),
								IntLiteral(3),
								IntLiteral(4),
								IntLiteral(5),
							},
						},
					},
					MutStmt{
						Var: VarWithType{"n2",
							ArrayType{
								Base: TypeLiteral("byte"),
								Size: IntLiteral(5),
							},
							false,
						},
						InitialValue: ArrayLiteral{
							Values: []Value{
								IntLiteral(6),
								IntLiteral(7),
								IntLiteral(8),
								IntLiteral(9),
								IntLiteral(10),
							},
						},
					},
					FuncCall{
						Name: "foo",
						UserArgs: []Value{
							VarWithType{"n",
								ArrayType{
									Base: TypeLiteral("byte"),
									Size: IntLiteral(5),
								},
								false,
							},
						},
					},
					FuncCall{
						Name: "foo",
						UserArgs: []Value{
							VarWithType{"n2",
								ArrayType{
									Base: TypeLiteral("byte"),
									Size: IntLiteral(5),
								},
								false,
							},
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v (%v, %v)", i, ast[i], v, reflect.TypeOf(ast[i]), reflect.TypeOf(v))
		}
	}
}

func TestEnumArray(t *testing.T) {
	ast, _, _ := buildAst(t, "enumarray")

	expected := []Node{
		EnumTypeDefn{
			"Light",
			[]EnumOption{
				EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}},
				EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}},
				EnumOption{"Amber", nil, UserType{TypeLiteral("int64"), "Light"}},
			},
			nil,
			0,
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"foo",
							ArrayType{
								Base: EnumTypeDefn{
									"Light",
									[]EnumOption{
										EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}},
										EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}},
										EnumOption{"Amber", nil, UserType{TypeLiteral("int64"), "Light"}},
									},
									nil,
									0,
								},
								Size: IntLiteral(2),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								EnumValue{Constructor: EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}}},
								EnumValue{Constructor: EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}}},
							},
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: FuncCall{
								Name: "len",
								UserArgs: []Value{
									VarWithType{
										"foo",
										ArrayType{
											Base: EnumTypeDefn{
												"Light",
												[]EnumOption{
													EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}},
													EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}},
													EnumOption{"Amber", nil, UserType{TypeLiteral("int64"), "Light"}},
												},
												nil,
												0,
											},
											Size: IntLiteral(2),
										},
										false,
									},
								},
							},
							Right: IntLiteral(2),
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v (%v, %v)", i, ast[i], v, reflect.TypeOf(ast[i]), reflect.TypeOf(v))
		}
	}
}

func TestEnumArrayExplicit(t *testing.T) {
	ast, _, _ := buildAst(t, "enumarrayexplicit")

	expected := []Node{
		EnumTypeDefn{
			"Light",
			[]EnumOption{
				EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}},
				EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}},
				EnumOption{"Amber", nil, UserType{TypeLiteral("int64"), "Light"}},
			},
			nil,
			0,
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"foo",
							ArrayType{
								Base: EnumTypeDefn{
									"Light",
									[]EnumOption{
										EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}},
										EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}},
										EnumOption{"Amber", nil, UserType{TypeLiteral("int64"), "Light"}},
									},
									nil,
									0,
								},
								Size: IntLiteral(2),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								EnumValue{Constructor: EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}}},
								EnumValue{Constructor: EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}}},
							},
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: FuncCall{
								Name: "len",
								UserArgs: []Value{
									VarWithType{
										"foo",
										ArrayType{
											Base: EnumTypeDefn{
												"Light",
												[]EnumOption{
													EnumOption{"Red", nil, UserType{TypeLiteral("int64"), "Light"}},
													EnumOption{"Green", nil, UserType{TypeLiteral("int64"), "Light"}},
													EnumOption{"Amber", nil, UserType{TypeLiteral("int64"), "Light"}},
												},
												nil,
												0,
											},
											Size: IntLiteral(2),
										},
										false,
									},
								},
							},
							Right: IntLiteral(2),
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v (%v, %v)", i, ast[i], v, reflect.TypeOf(ast[i]), reflect.TypeOf(v))
		}
	}
}

func TestSumTypeArray(t *testing.T) {
	ast, _, _ := buildAst(t, "sumtypearray")

	expected := []Node{
		TypeDefn{
			Name: "TestType",
			ConcreteType: SumType{
				TypeLiteral("string"),
				TypeLiteral("int"),
			},
		},
		FuncDecl{
			Name:    "main",
			Args:    nil,
			Return:  nil,
			Effects: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"foo",
							ArrayType{
								Base: UserType{
									SumType{
										TypeLiteral("string"),
										TypeLiteral("int"),
									},
									"TestType",
								},
								Size: IntLiteral(2),
							},
							false,
						},
						Val: ArrayLiteral{
							Values: []Value{
								StringLiteral("string"),
								IntLiteral(33),
							},
						},
					},
					Assertion{
						Predicate: EqualityComparison{
							Left: FuncCall{
								Name: "len",
								UserArgs: []Value{

									VarWithType{
										"foo",
										ArrayType{
											Base: UserType{
												SumType{
													TypeLiteral("string"),
													TypeLiteral("int"),
												},
												"TestType",
											},
											Size: IntLiteral(2),
										},
										false,
									},
								},
							},
							Right: IntLiteral(2),
						},
					},
				},
			},
		},
	}

	for i, v := range expected {
		if !compare(ast[i], v) {
			t.Errorf("Node %d: got %v want %v (%v, %v)", i, ast[i], v, reflect.TypeOf(ast[i]), reflect.TypeOf(v))
		}
	}
}
