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

const WrongArgType = `
func foo(s int) (int) {
	return s+5
}

proc main() () {
	foo("hello")
}
`

const WrongArgUserType = `
type fint int
func foo(s fint) (int) {
	return s+5
}

proc main() () {
	let x int = 5
	foo(x)
}
`
