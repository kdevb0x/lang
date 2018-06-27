// Tests that a slice can be taken from an array
func main () () {
	let x []byte = { 1, 2, 3, 4, 5 }
	let y = x[3:5]

	PrintInt(y[0])
	PrintInt(y[1])
	assert(len(y) == 2)
	assert(y[0] == 4)
	assert(y[1] == 5)
}
