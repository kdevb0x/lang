package sampleprograms

// SimpleAlgorithm is a port of the non-cheating C version of the algorithm used at
// https://www.fpcomplete.com/blog/2017/07/iterators-streams-rust-haskell
// The proc main ony goes up to 10 instead of 1000000 since the purpose in this context
// is using the algorithm as a test case for the parser/code generator, not as a benchmark.
const SimpleAlgorithm = `func loop(high int) (int) {
	mutable total = 0
	mutable i = 0
	let high = high * 2
	i = 1 
	while i < high {
		if i % 2 == 0 {
			total = total +  i*2
		}
		i = i + 1
	}
	return total
}

func main () () -> affects(IO) {
	PrintInt(loop(10))
}`
