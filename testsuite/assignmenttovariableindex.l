func main () () -> affects(IO) {
	mutable x = { 1, 3, 4, 5 }
	let y = x[0]
	x[y] = 6
	PrintInt(x[y])
	PrintInt(x[y+1])
}
