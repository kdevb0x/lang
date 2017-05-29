package sampleprograms

const TypeInference = `
func foo(x int) (int) {
	mut a = x
	a = a + 1

	let x = a + 1
	if x > 3 {
		return a
	}
	return 0
}

proc main() () {
	print("%d, %d\n", foo(1), foo(3))
}`
