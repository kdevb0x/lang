package sampleprograms

// SimpleSlice creates and initializes a slice
const SimpleSlice = `func main() () : io {
	let n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
}`

// SimpleSliceInference tests that slice types can be inferred.
// It needs one level of indirection to ensure it's not inferred
// as an array..
const SimpleSliceInference = `func main() () : io {
	let n []int = { 1, 2, 3, 4, 5 }
	let n2 = n
	PrintInt(n2[3])
}`

// ArrayMutation tests mutating an array value.
const SliceMutation = `func main() () : io {
	mutable n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
	PrintString("\n")
	n[3] = 2
	PrintInt(n[3])
	PrintString("\n")
	PrintInt(n[2])
}`

// SliceParam tests passing a slice as a parameter.
const SliceParam = `func main() () : io {
	let b []byte = { 44, 55, 88 }
	PrintASlice(b)
}

func PrintASlice(A []byte) () : io {
	PrintByteSlice(A)
}
`

const SliceStringParam = `func PrintSecond(args []string) () : io {
	PrintString(args[1])
}

func main() () : io {
	let aslice []string = {"foo", "bar", "baz" }
	PrintSecond(aslice)
}`

const SliceStringVariableParam = `func PrintSecond(args []string) () : io {
	let i = 1
	PrintString(args[i])
}

func main() () : io {
	let aslice []string = {"foo", "bar", "baz" }
	PrintSecond(aslice)
}`
