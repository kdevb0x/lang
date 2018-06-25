// Test that equality comparisons work
func main() () -> affects(IO) {
	mutable a int = 3
	let b int = 3
	if a == b {
		assert(true)
		PrintString("true\n")
	} else {
		assert(false)
		PrintString("false\n")
	}

	while a == b {
		PrintInt(a)
		PrintString("\n")
		a = a + 1
		assert(a != b)
	}
}
