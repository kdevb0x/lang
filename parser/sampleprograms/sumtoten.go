package sampleprograms

// SumToTen creates a program which calls a pure function which sums the
// numbers from 1 to 10, written in a procedural fashion.
const SumToTen = `func sum(x int) (int) {
	mutable val int = x
	mutable sum int = 0
	while val > 0 {
		sum = sum + val
		val = val - 1
	}
	return sum
}

func main() () -> affects(IO) {
	PrintInt(sum(10))
}`

// SumToTenRecursive does the same thing as SumToTen, but is written using
// tail call recursion instead of loops.
const SumToTenRecursive = `func sum(x int) (int) {
	return partial_sum(0, x)
}

func partial_sum(partial int, x int) (int) {
	if x == 0 {
		return partial
	}

	return partial_sum(partial + x, x - 1)
}

func main() () -> affects(IO) {
	PrintInt(sum(10))
}
`
