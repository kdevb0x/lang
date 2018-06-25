enum Foo = A | B

func main() () -> affects(IO) {
	let a Foo = A
	match a {
	case A:
		PrintString("I am A!\n")
	case B:
		PrintString("I am B!\n")
	}
}
