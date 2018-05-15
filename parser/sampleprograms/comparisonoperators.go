package sampleprograms

const EqualComparison = `func main() () : io {
	mutable a int = 3
	let b int = 3
	if a == b {
		PrintString("true\n")
	} else {
		PrintString("false\n")
	}

	while a == b {
		PrintInt(a)
		PrintString("\n")
		a = a + 1
	}
}`

const NotEqualComparison = `func main() () : io {
	mutable a int = 3
	let b int = 3
	if a != b {
		PrintString("true\n")
	} else {
		PrintString("false\n")
	}

	while a != b {
		PrintInt(a)
		PrintString("\n")
		a = a + 1
	}
}`

const GreaterComparison = `func main() () : io {
	mutable a int = 4
	let b int = 3
	if a > b {
		PrintString("true\n")
	} else {
		PrintString("false\n")

	}

	while a > b {
		PrintInt(a)
		PrintString("\n")
		a = a - 1
	}
}`

const GreaterOrEqualComparison = `func main() () : io {
	mutable a int = 4
	let b int = 3
	if a >= b {
		PrintString("true\n")

	} else {
		PrintString("false\n")
	}

	while a >= b {
		PrintInt(a)
		PrintString("\n")
		a = a - 1
	}
}`

const LessThanComparison = `func main() () : io {
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
}`

const LessThanOrEqualComparison = `func main() () : io {
	mutable a int = 1
	let b int = 3
	if a <= b {
		PrintString("true\n")

	} else {
		PrintString("false\n")
	}

	while a <= b {
		PrintInt(a)
		PrintString("\n")
		a = a + 1
	}
}`
