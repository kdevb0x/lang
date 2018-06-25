// ArrayMutation tests mutating an array value.
func main () () -> affects(IO) {
	mutable n = { 1, 2, 3, 4, 5 }
	PrintInt(n[3])
	PrintString("\n")
	n[3] = 2
	PrintInt(n[3])
	PrintString("\n")
	PrintInt(n[2])
}
