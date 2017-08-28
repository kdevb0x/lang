package sampleprograms

const EqualComparison = `proc main() () {
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

const NotEqualComparison = `proc main() () {
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

const GreaterComparison = `proc main() () {
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

const GreaterOrEqualComparison = `proc main() () {
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

const LessThanComparison = `proc main() () {
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

const LessThanOrEqualComparison = `proc main() () {
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
