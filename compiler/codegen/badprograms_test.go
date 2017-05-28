package codegen

import (
	"fmt"

	"github.com/driusan/lang/parser/sampleprograms/invalidprograms"
)

func ExampleTooManyArgs() {
	if err := RunProgram("toomanyargs", invalidprograms.TooManyArguments); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Unexpected number of parameters to aFunc: got 1 want 0.
}

func ExampleTooFewArgs() {
	if err := RunProgram("toofewargs", invalidprograms.TooFewArguments); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Unexpected number of parameters to aFunc: got 0 want 1.
}

func ExampleBadProcCall() {
	if err := RunProgram("badproccall", invalidprograms.BadProcCall); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not call procedure from pure function.
}


func ExampleBadLetAssignment() {
	if err := RunProgram("badletassignment", invalidprograms.LetAssignment); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Can not assign to immutable let variable "x".
}

func ExampleWrongType() {
	if err := RunProgram("wrongtype", invalidprograms.WrongType); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Incompatible type assignment: can not assign string to int for variable "x".
}

func ExampleUndefinedVariable() {
	if err := RunProgram("undefinedvar", invalidprograms.UndefinedVariable); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Use of undefined variable "x".
}

func ExampleVariableDefinedLater() {
	if err := RunProgram("varlater", invalidprograms.VariableDefinedLater); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Use of undefined variable "x".
}

func ExampleWrongScope() {
	if err := RunProgram("wrongscope", invalidprograms.WrongScope); err != nil {
		fmt.Println(err.Error())
	}
	// Output: Use of undefined variable "x".
}
