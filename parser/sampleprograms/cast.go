package sampleprograms

const CastBuiltin = `func main() () : io {
	let foo []byte = { 70, 111, 111 }
	PrintString(cast(foo) as string)
}`

const CastBuiltin2 = `func main() () : io {
	let foo = "bar"
	PrintByteSlice(cast(foo) as []byte)
}`

const CastIntVariable = `proc main () () {
	let foo = 65
	let baz = cast(foo) as byte

	PrintInt(baz)
}
`
