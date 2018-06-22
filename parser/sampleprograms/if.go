package sampleprograms

// IfBool tests that "if x" works where if is a boolean variable.
const IfBool = `
func foo(x bool) (int) {
	if x {
		return 3
	}
	return 7
}

func main () () {
	PrintInt(foo(false))
	PrintInt(foo(true))
}
`
