func foo(x int) (int) {
	mutable a = x
	a = a + 1

	let x = a + 1
	if x > 3 {
		return a
	}
	return 0
}

func main() () -> affects(IO) {
	PrintInt(foo(1))
	PrintString(", ")
	PrintInt(foo(3))
	PrintString("\n")
}
