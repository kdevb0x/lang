// SumToTen creates a program which calls a pure function which sums the
// numbers from 1 to 10, written in a procedural fashion.
func sum(x int) (int) {
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
}
