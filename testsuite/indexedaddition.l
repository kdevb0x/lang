func main() () -> affects(IO) {
	let x []int = { 3, 4, 5 }
	mutable n = x[1]
	n = n + x[2]
	let n2 = x[2] + x[0]
	PrintInt(n)
	PrintString("\n")
	PrintInt(n2)
}
