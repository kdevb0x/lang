package sampleprograms

const IndexAssignment = `proc main() () {
	let x []int = { 3, 4, 5 }
	mutable n = x[1]
	let n2 = x[2]
	PrintInt(n)
	PrintString("\n")
	PrintInt(n2)
}`

const IndexedAddition = `proc main() () {
	let x []int = { 3, 4, 5 }
	mutable n = x[1]
	n = n + x[2]
	let n2 = x[2] + x[0]
	PrintInt(n)
	PrintString("\n")
	PrintInt(n2)
}
`

const AssignmentToConstantIndex = `proc main () () {
	mutable x = { 3, 4, 5 }
	x[1] = 6
	PrintInt(x[0])
	PrintInt(x[1])
	PrintInt(x[2])
}`

const AssignmentToVariableIndex = `proc main () () {
	mutable x = { 1, 3, 4, 5 }
	let y = x[0]
	x[y] = 6
	PrintInt(x[y])
	PrintInt(x[y+1])
}`

const AssignmentToSliceConstantIndex = `proc main () () {
	mutable x []byte = { 3, 4, 5 }
	x[1] = 6
	PrintInt(x[0])
	PrintInt(x[1])
	PrintInt(x[2])
}`

const AssignmentToSliceVariableIndex = `proc main () () {
	mutable x []byte = { 1, 3, 4, 5 }
	let y = x[0]
	x[y] = 6
	PrintInt(x[y])
	PrintInt(x[y+1])
}`
