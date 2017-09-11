package invalidprograms

// WrongType is a program which tries to assign an invalid value to a
// variable.
const WrongType = `proc main() () {
	let x int = "string"
	PrintInt(x)
}
`

const InvalidType = `proc main() () {
	let x fint = 3
	PrintInt(x)
}
`

const WrongUserType = `type fint int
proc main() () {
	let x int = 3
	let y fint = x
	PrintInt(x)
}
`
