// SliceParam tests passing a slice as a parameter.
func main() () -> affects(IO) {
	let b []byte = { 44, 55, 88 }
	PrintASlice(b)
}

func PrintASlice(A []byte) () -> affects(IO) {
	PrintByteSlice(A)
}
