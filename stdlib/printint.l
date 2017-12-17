func Reverse(buf []byte) ([]byte) {
	mutable fub = make([]byte, len(buf))
	
	mutable i = 0
	while i < len(buf) {
		fub[i] = buf[buf.len() - i]
		i = i + 1
	}
	return fub 
}

func Digits(x int) ([]byte) {
	mutable digits [20]byte = default

	mutable remainder = x
	mutable i = 0
	while remainder != 0 {
		let digit = cast(remainder % 10) as byte

		digits[i] = remainder + '0'

		remainder = remainder / 10
		i = i + 1
	}

	if x < 0 {
		digits[i] = '-'
	}
	return digits[0:i].Reverse()
}

// Before multiple dispatch
proc PrintInt(x int) () {
	PrintByteSlice(x.Digits())
}

// With multiple dispatch, assuming Print([]byte) is implemented
proc Print(x int) () {
	x.Digits().Print()
}
