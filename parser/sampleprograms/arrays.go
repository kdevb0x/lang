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
