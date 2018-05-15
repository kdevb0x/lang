package sampleprograms

// SomeMath does arbitrary math operations to ensure that they work.
const SomeMath = `func main() () : io {
	let add int = 1 + 2
	let sub int = 1 - 2
	let mul int = 2 * 3
	let div int = 6 / 2
	let x int = 1 + 2 * 3 - 4 / 2

	PrintString("Add: ")
	PrintInt(add)
	PrintString("\n")
	PrintString("Sub: ")
	PrintInt(sub)
	PrintString("\n")
	PrintString("Mul: ")
	PrintInt(mul)
	PrintString("\n")
	PrintString("Div: ")
	PrintInt(div)
	PrintString("\n")
	PrintString("Complex: ")
	PrintInt(x)
	PrintString("\n")
}
`

// Precedence tests that brackets properly adjust the precedence inside of arithmetic
// values.
//
// It should Print "-3"
const Precedence = `func main() () : io {
	let x = (1 + 2) * (3 - 4)
	PrintInt(x)
}
`
