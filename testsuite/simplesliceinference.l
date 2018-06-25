// SimpleSliceInference tests that slice types can be inferred.
// It needs one level of indirection to ensure it's not inferred
// as an array..
func main() () -> affects(IO) {
	let n []int = { 1, 2, 3, 4, 5 }
	let n2 = n
	PrintInt(n2[3])
}
