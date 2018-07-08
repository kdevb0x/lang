// UnbufferedCat2 is the same as UnbufferedCat, but uses
// let statement bindings in the while condition.
//
// (The syntax for this was not implemented when the UnbufferedCat
// test was first written.)
func main (args []string) () -> affects(IO) {
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
