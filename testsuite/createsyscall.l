func main () () -> affects(IO, Filesystem) {
	let fd = Create("foo.txt")
	Write(fd, cast("Hello\n") as []byte)
	Close(fd)
}
