// Tests that a slice can be taken from an array
func main () () {
	mutable x [5]byte = { 1, 2, 3, 4, 5 }
	mutable y = x[2:4]

	PrintInt(y[0])
	PrintInt(y[1])
	assert(len(y) == 2)
	assert(y[0] == 3)
	assert(y[1] == 4)
}
