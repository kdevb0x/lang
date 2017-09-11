package invalidprograms

// UndefinedVariable is a program which tries to use a variable
// that has not been defined.
const UndefinedVariable = `proc main() () {
	PrintInt(x)
}
`

// VariableDefinedLater is a program which tries to use a variable
// that's declared later than its usage.
const VariableDefinedLater = `proc main() () {
	PrintInt(x)
	let x int = 3
}
`

// WrongScope is a program which tries to use a variable
// that's declared in a different scope.
const WrongScope = `proc main() () {
	if true {
		let x int = 3
	}
	PrintInt(x)
}
`
