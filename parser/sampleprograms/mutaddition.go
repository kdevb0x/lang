package sampleprograms

// MutAddition tests that basic mutable variables
// and addition work. It should print "8", but does
// it in a convuluted way.
const MutAddition = `proc main() () {
	mut x int = 3
	mut y int = x + 1
	x = x + y + 1
	print("%d\n", x)
}`
