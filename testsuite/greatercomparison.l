func main() () -> affects(IO) {
	mutable a int = 4
	let b int = 3
	if a > b {
		assert(true)
	} else {
		assert(false, "Took the wrong branch of if statement")
	}

	while a > b {
		PrintInt(a)
		PrintString("\n")
		a = a - 1
	}
}
