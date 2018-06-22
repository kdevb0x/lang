package sampleprograms

const LineComment = `
// I am the documentation for main.
func main () () { // Hello I am a comment
	let x = 3
	PrintInt(x) // This is a comment about PrintInt
}
//} // (This would be invalid if it wasn't commented)
`

const BlockComment = `
/* I am the documentation for main. */
func main () () {
	let x = /* I am inline 4 + */ 3
	PrintInt(x)
	/* I
	span
	multiple
	lines
	*/
}
`
