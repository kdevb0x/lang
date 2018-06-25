// Test the write syscall to a hardcoded file descriptor (stderr)
func main () () -> affects(IO, Filesystem) {
	Write(1, cast("Stdout!") as []byte)
	Write(2, cast("Stderr!") as []byte)
}
