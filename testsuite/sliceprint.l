// Slightly more thorough testing of slices.
// Prints "Foo" to stdout and "Bar" to stderr.
func main () () {
	let x [6]byte = { 70, 111, 111, 66, 97, 114 }
	let y = x[3:6]

	// "Foo"
	Write(1, x[0:3])
	// "Bar"
	Write(2, y)
}
