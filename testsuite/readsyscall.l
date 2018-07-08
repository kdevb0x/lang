// Test that the Open and Read syscalls work correctly. (Note: to use
// this test you need to know what's in the foo.txt file first.)
func main () () -> affects(IO, Filesystem) {
		let fd = Open("foo.txt")
		mutable dta []byte = {0, 1, 2, 3, 4, 5}
		let n = Read(fd, dta)
		PrintByteSlice(dta)
		Close(fd)
}
