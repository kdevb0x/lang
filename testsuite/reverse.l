// Tests the reverse portion of DigitsInto.
func Reverse(mutable digits [20]byte, size int) (int) {
	let swapto = (size / 2)
	mutable i = 0
	while i < swapto {
		assert(i >= 0)
		assert(i < 20)
		let tmp = digits[i]
		// FIXME: This shouldn't need the parenthesis, the order of operations
		// is wrong.
		let tmp2 = digits[(size-i)-1]

		digits[(size-i)-1] = tmp
		digits[i] = tmp2

		i = i + 1
	}

	i = size
	while i < 20 {
		assert(i >= 0)
		assert(i < 20)
		digits[i] = 0
		i = i + 1
	}
	return i
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
	digits.Reverse(3)

	assert(digits[0] == 3)
	assert(digits[1] == 2)
	assert(digits[2] == 1)

	// Test an even length
	digits[0] = 1
	digits[1] = 2
	digits[2] = 3
	digits[3] = 4
	digits.Reverse(4)

	assert(digits[0] == 4)
	assert(digits[1] == 3)
	assert(digits[2] == 2)
	assert(digits[3] == 1)
}
