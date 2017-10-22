package sampleprograms

// Cat implements the unix "cat" command.
// This implementation always uses an 8 byte buffer,
// primarily because the only way to declare a slice is
// currently with an array literal and there's no equivalent
// of malloc() or make(). This should be updated once it's
// implemented..
const UnbufferedCat = `proc main (args []string) () {
	mutable buf []byte = {0}

	mutable i = 1
	let length = len(args)
	while i < length {
		let file = Open(args[i])
		mutable n = Read(file, buf)
		PrintByteSlice(buf)
		while n > 0 {
			n = Read(file, buf)
			if n > 0 {
				PrintByteSlice(buf)
			}
		}

		i = i + 1
	}
}
`
