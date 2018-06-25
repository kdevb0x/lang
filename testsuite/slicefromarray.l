// Tests that a slice can be taken from an array
func main () () {
	let x [5]byte = { 1, 2, 3, 4, 5 }
	let y = x[2:4]

	assert(len(y) == 2)
	assert(y[0] == 3)
	assert(y[1] == 4)
}
