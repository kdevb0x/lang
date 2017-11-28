package sampleprograms

// Cat implements the unix "cat" command.
// This implementation always uses an 1 byte buffer,
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
		Close(file)

		i = i + 1
	}
}
`

// UnbufferedCat2 is the same as UnbufferedCat, but uses
// let statement bindings in the while condition.
const UnbufferedCat2 = `proc main (args []string) () {
	mutable buf []byte = {0}

	let i = 0
	while (let i = i + 1) < len(args) {
		let file = Open(args[i])
		while (let n = Read(file, buf)) > 0 {
			PrintByteSlice(buf)
		}
		Close(file)
	}
}
`
