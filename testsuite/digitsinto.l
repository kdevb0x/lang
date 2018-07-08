// This test is a prototype of the basis for PrintInt. It converts
// a integer to an array of byte characters.
func DigitsInto(x int, mutable digits [20]byte) (int) { // main.s:182
	if x == 0 {
		mutable i = 0
		while i < 20 {
			digits[i] = 0
			i = i + 1
		}
		digits[0] = 48
		return 1
	}

	mutable remainder = x // LV0+24
	if remainder < 0 {
		remainder = 0 - remainder
	}
	// Keep dividing by 10 and taking remainder mod 10 to get the digits (reversed)
	mutable i = 0 // LV1+32
	while remainder != 0 {
		assert(i < 20)
		let digit = cast(remainder % 10) as byte

		// 48 == '0'
		digits[i] = digit + 48

		remainder = remainder / 10
		i = i + 1
	}

	// Add the minus sign if applicable
	if x < 0 {
		digits[i] = 45 // 45 == '-'
		i = i + 1
	}
	let size = digits.Reverse(i)
	return size
}

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
	DigitsInto(678, digits)


	// Test for the value that we used.
	// 48 == '0'
	assert(digits[0] == 6+48)
	assert(digits[1] == 7+48)
	assert(digits[2] == 8+48)
	assert(digits[3] == 0)

	// Test for 0. Ensure that other digits were wiped out.
	DigitsInto(0, digits)
	assert(digits[0] == 0+48)
	assert(digits[1] == 0)
	assert(digits[2] == 0)
	assert(digits[3] == 0)

	// 678 was even length, so test an odd length
	// to ensure the swapping works for both cases.
	DigitsInto(1234, digits)
	assert(digits[0] == 1+48)
	assert(digits[1] == 2+48)
	assert(digits[2] == 3+48)
	assert(digits[3] == 4+48)
	assert(digits[4] == 0)

	// Try a negative number
	DigitsInto(-106, digits)
	assert(digits[0] == 45) // 45 = "-"
	assert(digits[1] == 1+48)
	assert(digits[2] == 0+48)
	assert(digits[3] == 6+48)
	assert(digits[4] == 0)

	// And an even length negative number
	let len = DigitsInto(-14, digits)
	assert(digits[0] == 45) // 45 = "-"
	assert(digits[1] == 1+48)
	assert(digits[2] == 4+48)
	assert(digits[3] == 0)

	DigitsInto(-1046, digits)
	assert(digits[0] == 45) // 45 = "-"
	assert(digits[1] == 1+48)
	assert(digits[2] == 0+48)
	assert(digits[3] == 4+48)
	assert(digits[4] == 6+48)
	assert(digits[5] == 0)
}
