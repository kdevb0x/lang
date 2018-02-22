package sampleprograms

// Fizzbuzz is a simple, well formatted fizzbuzz program
// to use for testing.
const Fizzbuzz = `func main () () -> affects(IO) {
	mutable terminate bool = false
	mutable i int = 1
	while terminate != true {
		if i % 15 == 0 {
			PrintString("fizzbuzz")
		} else if i % 5 == 0 {
			PrintString("buzz")
		} else if i % 3 == 0 {
			PrintString("fizz")
		} else {
			PrintInt(i)
		}
		PrintString("\n")

		i = i + 1
		if i >= 100 {
			terminate = true
		}
	}
}`
