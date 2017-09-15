package sampleprograms

// SimpleSlice creates and initializes a slice
const SimpleSlice = `proc main() () {
	let n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
}`

// SimpleSliceInference tests that slice types can be inferred.
// It needs one level of indirection to ensure it's not inferred
// as an array..
const SimpleSliceInference = `proc main() () {
	let n []int = { 1, 2, 3, 4, 5 }
	let n2 = n
	PrintInt(n2[3])
}`

// ArrayMutation tests mutating an array value.
const SliceMutation = `proc main() () {
	mutable n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
	PrintString("\n")
	n[3] = 2
	PrintInt(n[3])
	PrintString("\n")
	PrintInt(n[2])
}`
