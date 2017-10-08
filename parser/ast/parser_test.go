package ast

import (
	"fmt"
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
		Variable, VarWithType,
		AdditionOperator, SubtractionOperator, AssignmentOperator,
		MulOperator, DivOperator,
		EqualityComparison, NotEqualsComparison, GreaterComparison,
		GreaterOrEqualComparison, LessThanComparison, LessThanOrEqualComparison,
		TypeLiteral:
		return v1 == v2
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
			return compare(v1a.Var, v2a.Var) && compare(v1a.Value, v2a.Value)
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
		return compare(v1a.Val, v2a.Val)
	}

	if v1a, ok := v1.(ModOperator); ok {
		v2a, ok := v2.(ModOperator)
		if !ok {
			return false
		}
		return compare(v1a.Left, v2a.Left) && compare(v1a.Right, v2a.Right)
	}

	if v1a, ok := v1.(ProcDecl); ok {
		v2a, ok := v2.(ProcDecl)
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
		if compare(v1a.Body, v2a.Body) == false {
			return false
		}
		return true
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
	if v1a, ok := v1.(SumTypeDefn); ok {
		v2a, ok := v2.(SumTypeDefn)
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
		return v1a.ConcreteType == v2a.ConcreteType
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

	panic(fmt.Sprintf("Unimplemented type for compare %v vs %v", reflect.TypeOf(v1), reflect.TypeOf(v2)))
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var:   VarWithType{"n", TypeLiteral("int"), false},
						Value: IntLiteral(5),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var:   VarWithType{"n", TypeLiteral("int"), false},
						Value: IntLiteral(5),
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
						Var:   VarWithType{"n", TypeLiteral("string"), false},
						Value: StringLiteral("hello"),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
		ProcDecl{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,
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
		ProcDecl{
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
		ProcDecl{
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
		ProcDecl{
			Name: "main",
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,
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
		ProcDecl{
			Name: "main",
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{Name: "add", Typ: TypeLiteral("int")},
						Value: AdditionOperator{
							Left:  IntLiteral(1),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "sub", Typ: TypeLiteral("int")},
						Value: SubtractionOperator{
							Left:  IntLiteral(1),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "mul", Typ: TypeLiteral("int")},
						Value: MulOperator{
							Left:  IntLiteral(2),
							Right: IntLiteral(3),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "div", Typ: TypeLiteral("int")},
						Value: DivOperator{
							Left:  IntLiteral(6),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "x", Typ: TypeLiteral("int")},
						Value: AdditionOperator{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(3),
					},
					LetStmt{
						Var:   VarWithType{"b", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(3),
					},
					LetStmt{
						Var:   VarWithType{"b", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var:   VarWithType{"b", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var:   VarWithType{"b", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var:   VarWithType{"b", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{"a", TypeLiteral("int"), false},
						InitialValue: IntLiteral(1),
					},
					LetStmt{
						Var:   VarWithType{"b", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
			Name:         TypeLiteral("Foo"),
			ConcreteType: TypeLiteral("int"),
		},
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var:   VarWithType{"x", TypeLiteral("Foo"), false},
						Value: IntLiteral(4),
					},
					FuncCall{
						Name:     "PrintInt",
						UserArgs: []Value{VarWithType{"x", TypeLiteral("Foo"), false}},
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
						Value: AdditionOperator{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
		SumTypeDefn{
			TypeLiteral("Foo"),
			[]EnumOption{
				EnumOption{"A", nil, TypeLiteral("Foo")},
				EnumOption{"B", nil, TypeLiteral("Foo")},
			},
			nil,
		},
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var:   VarWithType{"a", TypeLiteral("Foo"), false},
						Value: EnumValue{Constructor: EnumOption{"A", nil, TypeLiteral("Foo")}},
					},
					MatchStmt{
						Condition: VarWithType{"a", TypeLiteral("Foo"), false},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"A", nil, TypeLiteral("Foo")},
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
								Variable: EnumOption{"B", nil, TypeLiteral("Foo")},
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
		SumTypeDefn{
			TypeLiteral("Foo"),
			[]EnumOption{
				EnumOption{"A", nil, TypeLiteral("Foo")},
				EnumOption{"B", nil, TypeLiteral("Foo")},
			},
			nil,
		},
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var:   VarWithType{"a", TypeLiteral("Foo"), false},
						Value: EnumValue{Constructor: EnumOption{"B", nil, TypeLiteral("Foo")}},
					},
					MatchStmt{
						Condition: VarWithType{"a", TypeLiteral("Foo"), false},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"A", nil, TypeLiteral("Foo")},
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
								Variable: EnumOption{"B", nil, TypeLiteral("Foo")},
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var:   VarWithType{"x", TypeLiteral("int"), false},
						Value: IntLiteral(3),
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
		SumTypeDefn{
			TypeLiteral("Maybe"),
			[]EnumOption{
				EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
				EnumOption{"Just", []Type{TypeLiteral("a")}, TypeLiteral("Maybe")},
			},
			nil,
		},
		FuncDecl{
			Name: "DoSomething",
			Args: []VarWithType{{"x", TypeLiteral("int"), false}},
			Return: []VarWithType{
				{
					"",
					TypeLiteral("Maybe int"),
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
								ReturnStmt{EnumValue{Constructor: EnumOption{"Nothing", nil, TypeLiteral("Maybe")}}},
							},
						},
					},
					ReturnStmt{EnumValue{
						Constructor: EnumOption{"Just", []Type{TypeLiteral("a")}, TypeLiteral("Maybe")},
						Parameters:  []Value{IntLiteral(5)},
					},
					},
				},
			},
		},
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{

					LetStmt{
						Var: VarWithType{"x", TypeLiteral("Maybe int"), false},
						Value: FuncCall{
							Name: "DoSomething",
							UserArgs: []Value{
								IntLiteral(3),
							},
						},
					},
					MatchStmt{
						Condition: VarWithType{"x", TypeLiteral("Maybe int"), false},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
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
								Variable: EnumOption{"Just", []Type{TypeLiteral("a")}, TypeLiteral("Maybe")},
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
						Var: VarWithType{"x", TypeLiteral("Maybe int"), false},
						Value: FuncCall{
							Name: "DoSomething",
							UserArgs: []Value{
								IntLiteral(4),
							},
						},
					},
					MatchStmt{
						Condition: VarWithType{"x", TypeLiteral("Maybe int"), false},
						Cases: []MatchCase{
							MatchCase{
								Variable: EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
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
								Variable: EnumOption{"Just", []Type{TypeLiteral("a")}, TypeLiteral("Maybe")},
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
		SumTypeDefn{
			TypeLiteral("Maybe"),
			[]EnumOption{
				EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
				EnumOption{"Just", []Type{TypeLiteral("x")}, TypeLiteral("Maybe")},
			},
			nil,
		},
		FuncDecl{
			Name: "foo",
			Args: []VarWithType{{"x", TypeLiteral("Maybe int"), false}},
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
						Condition: VarWithType{"x", TypeLiteral("Maybe int"), false},
						Cases: []MatchCase{
							MatchCase{
								LocalVariables: []VarWithType{
									VarWithType{"n", TypeLiteral("int"), false},
								},
								Variable: EnumOption{"Just", []Type{TypeLiteral("x")}, TypeLiteral("Maybe")},
								Body: BlockStmt{
									[]Node{
										ReturnStmt{VarWithType{"n", TypeLiteral("int"), false}},
									},
								},
							},
							MatchCase{
								Variable: EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{

							FuncCall{
								Name: "foo",

								UserArgs: []Value{
									EnumValue{
										Constructor: EnumOption{"Just", []Type{TypeLiteral("x")}, TypeLiteral("Maybe")},
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
		SumTypeDefn{
			TypeLiteral("Maybe"),
			[]EnumOption{
				EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
				EnumOption{"Just", []Type{TypeLiteral("x")}, TypeLiteral("Maybe")},
			},
			nil,
		},
		ProcDecl{
			Name: "foo",
			Args: []VarWithType{{"x", TypeLiteral("Maybe int"), false}},
			Return: []VarWithType{
				{
					"",
					TypeLiteral("int"),
					false,
				},
			},
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintString",
						UserArgs: []Value{
							StringLiteral(`x`),
						},
					},
					MatchStmt{
						Condition: VarWithType{"x", TypeLiteral("Maybe int"), false},
						Cases: []MatchCase{
							MatchCase{
								LocalVariables: []VarWithType{
									VarWithType{"n", TypeLiteral("int"), false},
								},
								Variable: EnumOption{"Just", []Type{TypeLiteral("x")}, TypeLiteral("Maybe")},
								Body: BlockStmt{
									[]Node{
										ReturnStmt{VarWithType{"n", TypeLiteral("int"), false}},
									},
								},
							},
							MatchCase{
								Variable: EnumOption{"Nothing", nil, TypeLiteral("Maybe")},
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "PrintInt",
						UserArgs: []Value{
							FuncCall{
								Name: "foo",

								UserArgs: []Value{
									EnumValue{
										Constructor: EnumOption{"Just", []Type{TypeLiteral("x")}, TypeLiteral("Maybe")},
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
						Value: MulOperator{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,
			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var:   VarWithType{"x", TypeLiteral("int64"), false},
						Value: IntLiteral(-4),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
						Value: ArrayLiteral{
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
						Value: ArrayLiteral{
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
						Variable: VarWithType{"n[3]", TypeLiteral("int"), false},
						Value:    IntLiteral(2),
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
		ProcDecl{
			Name: "changer",
			Args: []VarWithType{
				VarWithType{"x", TypeLiteral("int"), true},
				VarWithType{"y", TypeLiteral("int"), false},
			},
			Return: []VarWithType{
				VarWithType{Name: "", Typ: TypeLiteral("int")},
			},
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

		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
						Value: FuncCall{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						Value: ArrayLiteral{
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{"n",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						Value: ArrayLiteral{
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
						},
					},
					LetStmt{
						Var: VarWithType{"n2",
							SliceType{
								Base: TypeLiteral("int"),
							},
							false,
						},
						Value: VarWithType{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
						Variable: VarWithType{"n[3]", TypeLiteral("int"), false},
						Value:    IntLiteral(2),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

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
						Value: ArrayLiteral{
							IntLiteral(44),
							IntLiteral(55),
							IntLiteral(88),
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
		ProcDecl{
			Name: "PrintASlice",
			Args: []VarWithType{
				{Name: "A", Typ: SliceType{Base: TypeLiteral("byte")}},
			},
			Return: nil,

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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"fd",
							TypeLiteral("uint64"),
							false,
						},
						Value: FuncCall{
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
							IntLiteral(0),
							IntLiteral(1),
							IntLiteral(2),
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
						},
					},
					LetStmt{
						Var: VarWithType{
							"n",
							TypeLiteral("uint64"),
							false,
						},
						Value: FuncCall{
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SliceType{TypeLiteral("string")},
							false,
						},
						Value: ArrayLiteral{
							StringLiteral("3"),
							StringLiteral("foo"),
							StringLiteral("hello"),
							StringLiteral("world"),
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SliceType{TypeLiteral("int")},
							false,
						},
						Value: ArrayLiteral{
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
						Value: ArrayValue{
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

/*
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
		ProcDecl{
			Name:   "main",
			Args:   nil,
			Return: nil,

			Body: BlockStmt{
				[]Node{
					LetStmt{
						Var: VarWithType{
							"x",
							SliceType{TypeLiteral("int")},
							false,
						},
						Value: ArrayLiteral{
							IntLiteral(3),
							IntLiteral(4),
							IntLiteral(5),
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
						Value: ArrayValue{
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
*/
