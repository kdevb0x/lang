package sampleprograms

const IndexAssignment = `func main() () -> affects(IO){
	let x []int = { 3, 4, 5 }
	mutable n = x[1]
	let n2 = x[2]
	PrintInt(n)
	PrintString("\n")
	PrintInt(n2)
}`

const IndexedAddition = `func main() () -> affects(IO) {
	let x []int = { 3, 4, 5 }
	mutable n = x[1]
	n = n + x[2]
	let n2 = x[2] + x[0]
	PrintInt(n)
	PrintString("\n")
	PrintInt(n2)
}
`

const AssignmentToConstantIndex = `func main () () -> affects(IO) {
	mutable x = { 3, 4, 5 }
	x[1] = 6
	PrintInt(x[0])
	PrintInt(x[1])
	PrintInt(x[2])
}`

const AssignmentToVariableIndex = `func main () () -> affects(IO) {
	mutable x = { 1, 3, 4, 5 }
	let y = x[0]
	x[y] = 6
	PrintInt(x[y])
	PrintInt(x[y+1])
}`

const AssignmentToSliceConstantIndex = `func main () () -> affects(IO) {
	mutable x []byte = { 3, 4, 5 }
	x[1] = 6
	PrintInt(x[0])
	PrintInt(x[1])
	PrintInt(x[2])
}`

const AssignmentToSliceVariableIndex = `func main () () -> affects(IO) {
	mutable x []byte = { 1, 3, 4, 5 }
	let y = x[0]
	x[y] = 6
	PrintInt(x[y])
	PrintInt(x[y+1])
}`
