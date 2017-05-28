package invalidprograms

// LetAssignment is a program that tries to mutate a let variable,
// rather than shadowing it.
const LetAssignment = `proc main() () {
	let x int = 3
	x = x + 5
	printf("%d\n", x)
}
`
