// Test that not equals comparisons work as expected
func main() () -> affects(IO) {
	mutable a int = 3
	let b int = 3
	if a != b {
		assert(false, "Took the wrong branch of if statement")
	} else {
		assert(true)
	}

	while a != b {
		assert(false, "Unexpectedly entered while loop")
		PrintInt(a)
		PrintString("\n")
		a = a + 1
	}
}
