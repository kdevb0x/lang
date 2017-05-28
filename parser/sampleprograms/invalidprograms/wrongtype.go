package invalidprograms

// WrongType is a program which tries to assign an invalid value to a
// variable.
const WrongType = `proc main() () {
	let x int = "string"
	print("%d\n", x)
}
`

const InvalidType = `proc main() () {
	let x fint = 3
	print("%d\n", x)
}
`
