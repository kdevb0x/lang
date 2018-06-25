func main() () -> affects(IO) {
	mutable a int = 4
	let b int = 3
	if a < b {
		PrintString("true\n")
	} else {
		PrintString("false\n")
	}

	while a < b {
		PrintInt(a)
		PrintString("\n")
		a = a + 1
	}
}
