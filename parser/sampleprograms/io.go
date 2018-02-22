package sampleprograms

const PrintString = `
func main () () -> affects(IO) {
	PrintString("Success!")
}`

// Test the write syscall to a hardcoded file descriptor (stderr)
const WriteSyscall = `
func main () () -> affects(IO, Filesystem) {
	Write(1, cast("Stdout!") as []byte)
	Write(2, cast("Stderr!") as []byte)
}`

// Test that the Open and Read syscalls work correctly. (Note: to use
// this test you need to know what's in the foo.txt file first.)
const ReadSyscall = `
func main () () -> affects(IO, Filesystem) {
		let fd = Open("foo.txt")
		mutable dta []byte = {0, 1, 2, 3, 4, 5}
		let n = Read(fd, dta)
		PrintByteSlice(dta)
		Close(fd)
}
`

// Tests that the Create and Write syscalls work (Note: to use this
// as a test you need to be able to read foo.txt after.)
const CreateSyscall = `
func main () () -> affects(IO, Filesystem) {
	let fd = Create("foo.txt")
	Write(fd, cast("Hello\n") as []byte)
	Close(fd)
}
`
