package sampleprograms

// TwoProcs defines two trivial procedures in the same file.
const TwoProcs = `proc foo() (int) {
	return 3
}

proc main() () {
	print("%d", foo())
}`
