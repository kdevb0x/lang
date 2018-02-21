package stdlib

const PrintByteSlice = `proc PrintByteSlice(buf []byte) () {
	Write(1, buf)
}
`

const PrintString = `proc PrintString(str string) () {
	Write(1, str)
}
`
