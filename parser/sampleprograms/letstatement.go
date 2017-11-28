package sampleprograms

// LetStatement creates a let variable and prints it, to ensure the
// compiler is working.
const LetStatement = `proc main() () {
	let n int = 5
	PrintInt(n)
}`

// LetStatementShadow creates a let statement, and shadows it with
// another let statement.
const LetStatementShadow = `proc main() () {
	let n int = 5
	PrintInt(n)
	PrintString("\n")
	let n string = "hello"
	PrintString(n)
}`

// LetCondition tests creating let statements in a conditional.
const LetCondition = `proc main() () {
	let i = 0
	if (let i = i + 1) == 1 {
		PrintInt(i)
	} else {
		PrintInt(-1)
	}

	if (let i = i + 1) != 1 {
		PrintInt(i)
	} else {
		PrintInt(-1)
	}

	while (let i = i + 1) < 3 {
		PrintInt(i)
	}
}
`
