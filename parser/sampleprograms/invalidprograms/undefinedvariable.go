package invalidprograms

// UndefinedVariable is a program which tries to use a variable
// that has not been defined.
const UndefinedVariable = `func main() () : io {
	PrintInt(x)
}
`

// VariableDefinedLater is a program which tries to use a variable
// that's declared later than its usage.
const VariableDefinedLater = `func main() () : io {
	PrintInt(x)
	let x int = 3
}
`

// WrongScope is a program which tries to use a variable
// that's declared in a different scope.
const WrongScope = `func main() () : io {
	if true {
		let x int = 3
	}
	PrintInt(x)
}
`
