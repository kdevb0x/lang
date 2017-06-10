package invalidprograms

const IncompleteMatch = `
data Foo = A | B | C

proc main() () {
	let x = A
	match x {
	case A:
		print("I am A\n")
	case B:
		print("I am B\n")
	}
}`
