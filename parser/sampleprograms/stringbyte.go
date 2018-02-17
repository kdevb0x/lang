package sampleprograms

// WriteStringByte ensures that Write works with both strings and
// bytes (implying that they have the same representation when passed
// as a parameter.
const WriteStringByte = `proc main() () {
	let str string = "hello"
	let bty []byte= { 104, 101,  108, 108, 111 }
	Write(1, str)
	Write(1, bty)
}
`
