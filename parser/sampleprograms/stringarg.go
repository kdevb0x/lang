package sampleprograms

// Stringarg tests passing a string as a parameter.
const StringArg = `proc main() () {
	let b string = "foobar"
	PrintAString(b)
}

proc PrintAString(str string) () {
	PrintString(str)

}
`
