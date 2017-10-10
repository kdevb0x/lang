package sampleprograms

// SimpleArray creates and initializes an array
const SimpleArray = `proc main() () {
	let n [5]int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
}`

// SimpleArrayInference tests type inference on an array literal
const SimpleArrayInference = `proc main() () {
	let n = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
}`

// ArrayMutation tests mutating an array value.
const ArrayMutation = `proc main() () {
	mutable n = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
	PrintString("\n")
	n[3] = 2
	PrintInt(n[3])
	PrintString("\n")
	PrintInt(n[2])
}`

// ArrayIndex tests indexing into an array by a variable
const ArrayIndex = `
proc main () () {
	let x = 3
	let n = { 1, 2, 3, 4, 5 }
	mutable n2 = { 1, 2, 3, 4, 5 }
	PrintInt(n[x])
	PrintString("\n")
	PrintInt(n2[x+1])
}
`

const StringArray = `
proc main () () {
	let args = { "foo", "bar" }
	PrintString(args[1])
	PrintString("\n")
	PrintString(args[0])
}
`
