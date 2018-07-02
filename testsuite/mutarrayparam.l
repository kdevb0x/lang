func foo(mutable x [5]byte) () {
	x[3] = 99
	return
}

func main () () {
	mutable n [5]byte = { 1, 2, 3, 4, 5 }
	assert(n[3] == 4)
	foo(n)
	assert(n[3] == 99)
	return
}
