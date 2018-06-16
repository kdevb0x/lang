package invalidprograms

const IncompleteMatch = `
enum Foo = A | B | C

func main() () -> affects(IO) {
	let x = A
	match x {
	case A:
		PrintString("I am A\n")
	case B:
		PrintString("I am B\n")
	}
}`
