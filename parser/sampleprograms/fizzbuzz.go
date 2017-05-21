package sampleprograms

// Fizzbuzz is a simple, well formatted fizzbuzz program
// to use for testing.
const Fizzbuzz = `proc main() () {
	mut terminate bool = false
	mut i int = 1
	while terminate != true {
		if i % 15 == 0 {
			print("fizzbuzz\n")
		} else if i % 5 == 0 {
			print("buzz\n")
		} else if i % 3 == 0 {
			print("fizz\n")
		} else {
			print("%d\n", i)
		}

		i = i + 1
		if i >= 100 {
			terminate = true
		}
	}
}`
