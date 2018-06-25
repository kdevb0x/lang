// SliceMutation tests mutating a slice value.
func main() () -> affects(IO) {
	mutable n []int = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
	PrintString("\n")
	n[3] = 2
	PrintInt(n[3])
	PrintString("\n")
	PrintInt(n[2])
}
