package sampleprograms

// WriteStringByte ensures that Write works with both strings and
// bytes (implying that they have the same representation when passed
// as a parameter.
const WriteStringByte = `func main() () -> affects(IO) {
	let str string = "hello"
	let bty []byte= { 104, 101,  108, 108, 111 }
	Write(1, cast(str) as []byte)
	Write(1, bty)
}
`
