package invalidprograms

// LetAssignment is a program that tries to mutate a let variable,
// rather than shadowing it.
const LetAssignment = `proc main() () {
	let x int = 3
	x = x + 5
	printf("%d\n", x)
}
`

// MutStatementShadow creates a mutable variable, and then tries to shadow
// it, which is illegal.
const MutStatementShadow = `proc main() () {
	mutable n int = 5
	print("%d\n", n)
	mutable n string = "hello"
	print("%s\n", n)
}`

// MutStatementShadow creates a mutable variable, and then tries to shadow
// it with a let statement, which is still illegal.
const MutStatementShadow2 = `proc main() () {
	mutable n int = 5
	print("%d\n", n)
	let n string = "hello"
	print("%s\n", n)
}`

// MutStatementScopeShadow creates a mutable variable, and tries to shadow
// it in a different scope, which is still illegal.
const MutStatementScopeShadow = `proc main() () {
	mutable n int = 5
	print("%d\n", n)
	if n == 5 {
		mutable n string = "hello"
		print("%s\n", n)
	}
}`

// MutStatementScopeShadow creates a mutable variable, and tries to shadow
// it with a let variable in a different scope, which is still illegal.
const MutStatementScopeShadow2 = `proc main() () {
	mutable n int = 5
	print("%d\n", n)
	if n == 5 {
		let n string = "hello"
		print("%s\n", n)
	}
}`
