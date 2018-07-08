// Tests the Swap portion of the Reverse portion of DigitsInto.
func Swap(mutable digits [20]byte, i int, size int) () {
	assert(i < size)
	assert(i >= 0)

	let tmp = digits[i]
	// FIXME: This shouldn't need the parenthesis, the order of operations
	// is wrong.
	let tmp2 = digits[(size-i)-1]

	digits[(size-i)-1] = tmp
	digits[i] = tmp2
}


func main () () {
	// 20 bytes = the largest uint64
	mutable digits [20]byte = {
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0
	}
	digits[0] = 1
	digits[1] = 2
	digits[2] = 3
	//PrintString("123\n")
	digits.Swap(0, 3)

	assert(digits[0] == 3)
	assert(digits[1] == 2)
	assert(digits[2] == 1)
/*

	// Test an even length
	digits[0] = 1
	digits[1] = 2
	digits[2] = 3
	digits[3] = 4
	PrintString("\n\n1234\n")
	digits.Reverse(4)

	assert(digits[0] == 4)
	assert(digits[1] == 3)
	assert(digits[2] == 2)
	assert(digits[3] == 1)
	*/
}
