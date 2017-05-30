package sampleprograms

// LetStatement creates a let variable and prints it, to ensure the
// compiler is working.
const LetStatement = `proc main() () {
	let n int = 5
	print("%d\n", n)
}`

// LetStatementShadow creates a let statement, and shadows it with
// another let statement.
const LetStatementShadow = `proc main() () {
	let n int = 5
	print("%d\n", n)
	let n string = "hello"
	print("%s\n", n)
}`
