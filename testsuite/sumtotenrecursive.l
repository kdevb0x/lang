// SumToTenRecursive does the same thing as SumToTen, but is written using
// tail call recursion instead of loops.
func sum(x int) (int) {
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

