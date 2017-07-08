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
