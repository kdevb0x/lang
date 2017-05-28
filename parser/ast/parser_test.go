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
		LetStmt, MutStmt,
		AdditionOperator, SubtractionOperator, AssignmentOperator,
		MulOperator, DivOperator,
		EqualityComparison, NotEqualsComparison, GreaterComparison,
		GreaterOrEqualComparison, LessThanComparison, LessThanOrEqualComparison:
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
	panic(fmt.Sprintf("Unimplemented type for compare %v vs %v", reflect.TypeOf(v1), reflect.TypeOf(v2)))
}

func TestParseFizzBuzz(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.Fizzbuzz))
	if err != nil {
		t.Fatal(err)
	}
	ast, err := Construct(tokens)
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
						Var:          VarWithType{Name: "terminate", Typ: "bool"},
						InitialValue: BoolLiteral(false),
					},
					MutStmt{
						Var:          VarWithType{Name: "i", Typ: "int"},
						InitialValue: IntLiteral(1),
					},
					WhileLoop{
						Condition: NotEqualsComparison{
							Left:  VarWithType{"terminate", "bool"},
							Right: BoolLiteral(true),
						},
						Body: BlockStmt{
							[]Node{
								IfStmt{
									Condition: EqualityComparison{
										Left: ModOperator{
											Left:  VarWithType{Variable("i"), "int"},
											Right: IntLiteral(15),
										},
										Right: IntLiteral(0),
									},
									Body: BlockStmt{
										[]Node{
											FuncCall{
												Name: "print",
												UserArgs: []Value{
													StringLiteral(`fizzbuzz\n`),
												},
											},
										},
									},
									Else: BlockStmt{
										[]Node{
											IfStmt{
												Condition: EqualityComparison{
													Left: ModOperator{
														Left:  VarWithType{Variable("i"), "int"},
														Right: IntLiteral(5),
													},
													Right: IntLiteral(0),
												},
												Body: BlockStmt{
													[]Node{
														FuncCall{
															Name: "print",
															UserArgs: []Value{
																StringLiteral(`buzz\n`),
															},
														},
													},
												},
												Else: BlockStmt{
													[]Node{
														IfStmt{
															Condition: EqualityComparison{
																Left: ModOperator{
																	Left:  VarWithType{"i", "int"},
																	Right: IntLiteral(3),
																},
																Right: IntLiteral(0),
															},
															Body: BlockStmt{
																[]Node{
																	FuncCall{
																		Name: "print",
																		UserArgs: []Value{
																			StringLiteral(`fizz\n`),
																		},
																	},
																},
															},
															Else: BlockStmt{
																[]Node{
																	FuncCall{
																		Name: "print",
																		UserArgs: []Value{
																			StringLiteral(`%d\n`),
																			VarWithType{"i", "int"},
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
								AssignmentOperator{
									Variable: VarWithType{"i", "int"},
									Value: AdditionOperator{
										Left:  VarWithType{"i", "int"},
										Right: IntLiteral(1),
									},
								},
								IfStmt{
									Condition: GreaterOrEqualComparison{
										Left:  VarWithType{"i", "int"},
										Right: IntLiteral(100),
									},
									Body: BlockStmt{
										[]Node{
											AssignmentOperator{
												Variable: VarWithType{"terminate", "bool"},
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
	ast, err := Construct(tokens)
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
						Name: "print",
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
			t.Errorf("empty function: got %v want %v", ast[i], v)
		}
	}
}

func TestEmptyMain(t *testing.T) {
	tokens, err := token.Tokenize(strings.NewReader(sampleprograms.EmptyMain))
	if err != nil {
		t.Fatal(err)
	}
	ast, err := Construct(tokens)
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
	ast, err := Construct(tokens)
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
						Var:   VarWithType{"n", "int"},
						Value: IntLiteral(5),
					},
					FuncCall{
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d\n`),
							VarWithType{"n", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"x", "int"},
						InitialValue: IntLiteral(3),
					},
					MutStmt{
						Var: VarWithType{"y", "int"},
						InitialValue: AdditionOperator{
							Left:  VarWithType{"x", "int"},
							Right: IntLiteral(1),
						},
					},
					AssignmentOperator{
						Variable: VarWithType{"x", "int"},
						Value: AdditionOperator{
							Left: VarWithType{"x", "int"},
							Right: AdditionOperator{
								Left:  VarWithType{"y", "int"},
								Right: IntLiteral(1),
							},
						},
					},
					FuncCall{
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d\n`),
							VarWithType{"x", "int"},
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
	ast, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		ProcDecl{
			Name: "foo",
			Args: nil,
			Return: []VarWithType{
				VarWithType{Name: "", Typ: "int"},
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
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d`),
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
	ast, err := Construct(tokens)
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
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d`),
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
				VarWithType{Name: "", Typ: "int"},
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
	ast, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		ProcDecl{
			Name: "sum",
			Args: []VarWithType{
				{Name: "x", Typ: "int"},
			},
			Return: []VarWithType{
				{Typ: "int"},
			},
			Body: BlockStmt{
				[]Node{
					MutStmt{
						Var:          VarWithType{Name: "val", Typ: "int"},
						InitialValue: VarWithType{"x", "int"},
					},
					MutStmt{
						Var:          VarWithType{Name: "sum", Typ: "int"},
						InitialValue: IntLiteral(0),
					},
					WhileLoop{
						Condition: GreaterComparison{
							Left:  VarWithType{"val", "int"},
							Right: IntLiteral(0),
						},
						Body: BlockStmt{
							[]Node{
								AssignmentOperator{
									Variable: VarWithType{"sum", "int"},
									Value: AdditionOperator{
										Left:  VarWithType{"sum", "int"},
										Right: VarWithType{"val", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"val", "int"},
									Value: SubtractionOperator{
										Left:  VarWithType{"val", "int"},
										Right: IntLiteral(1),
									},
								},
							},
						},
					},
					ReturnStmt{VarWithType{"sum", "int"}},
				},
			},
		},
		ProcDecl{
			Name: "main",
			Body: BlockStmt{
				[]Node{
					FuncCall{
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d\n`),
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
	ast, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}
	expected := []Node{
		FuncDecl{
			Name: "foo",
			Args: nil,
			Return: []VarWithType{
				VarWithType{Name: "", Typ: "int"},
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
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d`),
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
	ast, err := Construct(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []Node{
		FuncDecl{
			Name: "sum",
			Args: []VarWithType{
				{Name: "x", Typ: "int"},
			},
			Return: []VarWithType{
				{Typ: "int"},
			},
			Body: BlockStmt{
				[]Node{
					ReturnStmt{
						Val: FuncCall{
							Name: "partial_sum",
							UserArgs: []Value{
								IntLiteral(0),
								VarWithType{"x", "int"},
							},
						},
					},
				},
			},
		},
		FuncDecl{
			Name: "partial_sum",
			Args: []VarWithType{
				{Name: "partial", Typ: "int"},
				{Name: "x", Typ: "int"},
			},
			Return: []VarWithType{
				{Typ: "int"},
			},
			Body: BlockStmt{
				[]Node{
					IfStmt{
						Condition: EqualityComparison{
							Left:  VarWithType{"x", "int"},
							Right: IntLiteral(0),
						},
						Body: BlockStmt{
							[]Node{
								ReturnStmt{
									Val: VarWithType{"partial", "int"},
								},
							},
						},
					},
					ReturnStmt{
						Val: FuncCall{
							Name: "partial_sum",
							UserArgs: []Value{
								AdditionOperator{
									Left:  VarWithType{"partial", "int"},
									Right: VarWithType{"x", "int"},
								},
								SubtractionOperator{
									Left:  VarWithType{"x", "int"},
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
						Name: "print",
						UserArgs: []Value{
							StringLiteral(`%d\n`),
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
	ast, err := Construct(tokens)
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
						Var: VarWithType{Name: "add", Typ: "int"},
						Value: AdditionOperator{
							Left:  IntLiteral(1),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "sub", Typ: "int"},
						Value: SubtractionOperator{
							Left:  IntLiteral(1),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "mul", Typ: "int"},
						Value: MulOperator{
							Left:  IntLiteral(2),
							Right: IntLiteral(3),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "div", Typ: "int"},
						Value: DivOperator{
							Left:  IntLiteral(6),
							Right: IntLiteral(2),
						},
					},
					LetStmt{
						Var: VarWithType{Name: "x", Typ: "int"},
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
						Name: "print",

						UserArgs: []Value{
							StringLiteral(`Add: %d\n`),
							VarWithType{"add", "int"},
						},
					},
					FuncCall{
						Name: "print",

						UserArgs: []Value{
							StringLiteral(`Sub: %d\n`),
							VarWithType{"sub", "int"},
						},
					},
					FuncCall{
						Name: "print",

						UserArgs: []Value{
							StringLiteral(`Mul: %d\n`),
							VarWithType{"mul", "int"},
						},
					},
					FuncCall{
						Name: "print",

						UserArgs: []Value{
							StringLiteral(`Div: %d\n`),
							VarWithType{"div", "int"},
						},
					},
					FuncCall{
						Name: "print",

						UserArgs: []Value{
							StringLiteral(`Complex: %d\n`),
							VarWithType{"x", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"a", "int"},
						InitialValue: IntLiteral(3),
					},
					LetStmt{
						Var:   VarWithType{"b", "int"},
						Value: IntLiteral(3),
					},
					IfStmt{
						Condition: EqualityComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: EqualityComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "print",
									UserArgs: []Value{
										StringLiteral(`%d\n`),
										VarWithType{"a", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", "int"},
									Value: AdditionOperator{
										Left:  VarWithType{"a", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"a", "int"},
						InitialValue: IntLiteral(3),
					},
					LetStmt{
						Var:   VarWithType{"b", "int"},
						Value: IntLiteral(3),
					},
					IfStmt{
						Condition: NotEqualsComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: NotEqualsComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "print",
									UserArgs: []Value{
										StringLiteral(`%d\n`),
										VarWithType{"a", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", "int"},
									Value: AdditionOperator{
										Left:  VarWithType{"a", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"a", "int"},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var:   VarWithType{"b", "int"},
						Value: IntLiteral(3),
					},
					IfStmt{
						Condition: GreaterComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: GreaterComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "print",
									UserArgs: []Value{
										StringLiteral(`%d\n`),
										VarWithType{"a", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", "int"},
									Value: SubtractionOperator{
										Left:  VarWithType{"a", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"a", "int"},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var:   VarWithType{"b", "int"},
						Value: IntLiteral(3),
					},
					IfStmt{
						Condition: GreaterOrEqualComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: GreaterOrEqualComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "print",
									UserArgs: []Value{
										StringLiteral(`%d\n`),
										VarWithType{"a", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", "int"},
									Value: SubtractionOperator{
										Left:  VarWithType{"a", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"a", "int"},
						InitialValue: IntLiteral(4),
					},
					LetStmt{
						Var:   VarWithType{"b", "int"},
						Value: IntLiteral(3),
					},
					IfStmt{
						Condition: LessThanComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: LessThanComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "print",
									UserArgs: []Value{
										StringLiteral(`%d\n`),
										VarWithType{"a", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", "int"},
									Value: AdditionOperator{
										Left:  VarWithType{"a", "int"},
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
	ast, err := Construct(tokens)
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
						Var:          VarWithType{"a", "int"},
						InitialValue: IntLiteral(1),
					},
					LetStmt{
						Var:   VarWithType{"b", "int"},
						Value: IntLiteral(3),
					},
					IfStmt{
						Condition: LessThanOrEqualComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`true\n`)},
								},
							},
						},
						Else: BlockStmt{
							[]Node{
								FuncCall{
									Name:     "print",
									UserArgs: []Value{StringLiteral(`false\n`)},
								},
							},
						},
					},
					WhileLoop{
						Condition: LessThanOrEqualComparison{
							Left:  VarWithType{"a", "int"},
							Right: VarWithType{"b", "int"},
						},
						Body: BlockStmt{
							[]Node{
								FuncCall{
									Name: "print",
									UserArgs: []Value{
										StringLiteral(`%d\n`),
										VarWithType{"a", "int"},
									},
								},
								AssignmentOperator{
									Variable: VarWithType{"a", "int"},
									Value: AdditionOperator{
										Left:  VarWithType{"a", "int"},
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
