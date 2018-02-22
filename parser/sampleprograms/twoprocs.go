package sampleprograms

// TwoProcs defines two trivial procedures in the same file.
const TwoProcs = `func foo () (int) {
	return 3
}

func main () () -> affects(IO) {
	PrintInt(foo())
}`
