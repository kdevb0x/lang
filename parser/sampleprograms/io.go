package sampleprograms

const PrintString = `
proc main () () {
	PrintString("Success!")
}`

// Test the write syscall to a hardcoded file descriptor (stderr)
const WriteSyscall = `
proc main () () {
	Write(1, "Stdout!")
	Write(2, "Stderr!")
}`

// Test that the Open and Read syscalls work correctly. (Note: to use
// this test you need to know what's in the foo.txt file first.)
const ReadSyscall = `
proc main () () {
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
proc main () () {
	let fd = Create("foo.txt")
	Write(fd, "Hello\n")
	Close(fd)
}
`
