package ast

import (
	"fmt"

	"github.com/driusan/lang/parser/sampleprograms/invalidprograms"
)

func buildAST(src string) error {
	_, _, _, err := Parse(src)
	if err != nil {
		return err
	}
	return nil
}

func ExampleTooManyArgs() {
	if err := buildAST(invalidprograms.TooManyArguments); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Unexpected number of parameters to aFunc: got 1 want 0.
}

func ExampleTooFewArgs() {
	if err := buildAST(invalidprograms.TooFewArguments); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Unexpected number of parameters to aFunc: got 0 want 1.
}

func ExampleBadLetAssignment() {
	if err := buildAST(invalidprograms.LetAssignment); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not assign to immutable let variable "x".
}

func ExampleWrongType() {
	if err := buildAST(invalidprograms.WrongType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Incompatible assignment for variable "x": Can not assign string to int.
}

func ExampleUndefinedVariable() {
	if err := buildAST(invalidprograms.UndefinedVariable); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Use of undefined variable "x".
}

func ExampleVariableDefinedLater() {
	if err := buildAST(invalidprograms.VariableDefinedLater); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Use of undefined variable "x".
}

func ExampleWrongScope() {
	if err := buildAST(invalidprograms.WrongScope); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Use of undefined variable "x".
}

func ExampleInvalidType() {
	if err := buildAST(invalidprograms.InvalidType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Invalid type: fint
}

func ExampleWrongUsertype() {
	if err := buildAST(invalidprograms.WrongUserType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Incompatible assignment for variable "y": can not assign int to fint.
}

func ExampleMutStatementShadow() {
	if err := buildAST(invalidprograms.MutStatementShadow); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not shadow mutable variable "n".
}

func ExampleMutStatementShadow2() {
	if err := buildAST(invalidprograms.MutStatementShadow2); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not shadow mutable variable "n".
}

func ExampleMutStatementScopeShadow() {
	if err := buildAST(invalidprograms.MutStatementScopeShadow); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not shadow mutable variable "n".
}

func ExampleMutStatementScopeShadow2() {
	if err := buildAST(invalidprograms.MutStatementScopeShadow2); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not shadow mutable variable "n".
}

func ExampleTooBigUInt8() {
	if err := buildAST(invalidprograms.TooBigUint8); err != nil {
		fmt.Println(err.Error())
	}

	// Output: Incompatible assignment for variable "y": value (256) must be between 0 and 255.
}

func ExampleIncompleteMatch() {
	if err := buildAST(invalidprograms.IncompleteMatch); err != nil {
		fmt.Println(err.Error())
	}

	// Output: Inexhaustive match for enum type "Foo": Missing case "C".
}

func ExampleWrongArgType() {
	if err := buildAST(invalidprograms.WrongArgType); err != nil {
		fmt.Println(err.Error())
	}

	// Output: Incompatible call to foo: argument s must be of type int (got string)
}

func ExampleWrongArgUserType() {
	if err := buildAST(invalidprograms.WrongArgUserType); err != nil {
		fmt.Println(err.Error())
	}

	// Output: Incompatible call to foo: argument s must be of type fint (got int)
}
