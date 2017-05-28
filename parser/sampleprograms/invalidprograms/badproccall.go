package invalidprograms

// BadProcCall is a program which calls a procedure with side-effects
// from a pure function, but is otherwise valid.
const BadProcCall = `proc main() () {
	let x int = aFunc()
	print("%d\n", x)
}

func aFunc() (int) {
	return aProc()
}

proc aProc() (int) {
	return 3
}
`
