// MutAddition tests that basic mutable variables
// and addition work. It should print "8", but does
// it in a convuluted way.
func main() () -> affects(IO) {
	mutable x int = 3
	mutable y int = x + 1
	x = x + y + 1
	PrintInt(x)
}
