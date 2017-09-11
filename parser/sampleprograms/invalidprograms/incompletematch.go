package invalidprograms

const IncompleteMatch = `
data Foo = A | B | C

proc main() () {
	let x = A
	match x {
	case A:
		PrintString("I am A\n")
	case B:
		PrintString("I am B\n")
	}
}`
