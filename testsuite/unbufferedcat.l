// Cat implements the unix "cat" command.
// This implementation always uses an 1 byte buffer.
func main (args []string) () -> affects(IO) {
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
