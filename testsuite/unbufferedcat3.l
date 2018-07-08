// UnbufferedCat3 is the same as UnbufferedCat2, but uses
// method invocation syntax too.
//
// (The method invocation syntax was implemented after let bindings.)
func main (args []string) () -> affects(IO) {
	mutable buf []byte = {0}

	let i = 0
	while (let i = i + 1) < args.len() {
		let file = Open(args[i])
		while (let n = file.Read(buf)) > 0 {
			PrintByteSlice(buf)
		}
		file.Close()
	}
}
