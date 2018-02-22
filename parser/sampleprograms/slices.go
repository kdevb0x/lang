package sampleprograms

// SimpleSlice creates and initializes a slice
const SimpleSlice = `func main() () -> affects(IO) {
	let n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
}`

// SimpleSliceInference tests that slice types can be inferred.
// It needs one level of indirection to ensure it's not inferred
// as an array..
const SimpleSliceInference = `func main() () -> affects(IO) {
	let n []int = { 1, 2, 3, 4, 5 }
	let n2 = n
	PrintInt(n2[3])
}`

// ArrayMutation tests mutating an array value.
const SliceMutation = `func main() () -> affects(IO) {
	mutable n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
	PrintString("\n")
	n[3] = 2
	PrintInt(n[3])
	PrintString("\n")
	PrintInt(n[2])
}`

// SliceParam tests passing a slice as a parameter.
const SliceParam = `func main() () -> affects(IO) {
	let b []byte = { 44, 55, 88 }
	PrintASlice(b)
}

func PrintASlice(A []byte) () -> affects(IO) {
	PrintByteSlice(A)
}
`

const SliceStringParam = `func PrintSecond(args []string) () -> affects(IO) {
	PrintString(args[1])
}

func main() () -> affects(IO) {
	let aslice []string = {"foo", "bar", "baz" }
	PrintSecond(aslice)
}`

const SliceStringVariableParam = `func PrintSecond(args []string) () -> affects(IO) {
	let i = 1
	PrintString(args[i])
}

func main() () -> affects(IO) {
	let aslice []string = {"foo", "bar", "baz" }
	PrintSecond(aslice)
}`
