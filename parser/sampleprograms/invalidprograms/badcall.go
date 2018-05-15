package invalidprograms

// TooManyArguments is a program which tries to call a function
// with too many arguments.
const TooManyArguments = `func main() () : io {
	let x int = aFunc(3)
	printf("%d\n", x)
}

func aFunc() (int) {
	return aProc()
}
`

// TooFewArguments is a program which tries to call a function
// without enough arguments.
const TooFewArguments = `func main() () : io {
	let x int = aFunc()
	printf("%d\n", x)
}

func aFunc(x int) (int) {
	return x
}
`
