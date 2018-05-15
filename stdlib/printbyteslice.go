package stdlib

const PrintByteSlice = `func PrintByteSlice(buf []byte) () : io {
	Write(1, buf)
}
`

const PrintString = `func PrintString(str string) () : io {
	Write(1, cast(str) as []byte)
}
`
