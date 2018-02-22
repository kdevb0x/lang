package stdlib

const PrintByteSlice = `func PrintByteSlice(buf []byte) () -> affects(IO) {
	Write(1, buf)
}
`

const PrintString = `func PrintString(str string) () -> affects(IO) {
	Write(1, cast(str) as []byte)
}
`
