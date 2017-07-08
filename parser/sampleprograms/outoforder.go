package sampleprograms

// OutOfOrder defines proc which is called before it's
// defined in the source. It should print "3".
const OutOfOrder = `proc main() () {
	PrintInt(foo())
}

proc foo() (int) {
	return 3
}
`
