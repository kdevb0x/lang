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
	PrintInt(foo(1))
	PrintString(", ")
	PrintInt(foo(3))
	PrintString("\n")
}`
