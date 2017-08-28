package sampleprograms

// SimpleArray creates and initializes an array
const ReferenceVariable = `proc changer(mutable x int, y int) (int) {
	x = 4
	return x + y
}

proc main() () {
	mutable var = 3
	PrintInt(var)
	PrintString("\n")
	let sum = changer(var, 3)

	PrintInt(var)
	PrintString("\n")

	PrintInt(sum)
}`
