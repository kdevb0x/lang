func foo(x [5]byte) () {
	PrintInt(x[0])
	assert(x[0] % 5 == 1)
	assert(x[1] % 5 == 2)
	assert(x[2] % 5 == 3)
	assert(x[3] % 5 == 4)
	assert(x[4] % 5 == 0)
}

func main () () {
	let n [5]byte = { 1, 2, 3, 4, 5 }
	mutable n2 [5]byte = { 6, 7, 8, 9, 10 }
	foo(n)
	foo(n2)
}
