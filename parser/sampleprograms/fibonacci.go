package sampleprograms

const Fibonacci = `func fib_rec(n uint64, n1 uint64) (uint64) : io {
	let n2 = n + n1
	if n2 >= 200 {
		return n1
	}
	PrintInt(n2)
	PrintString("\n")
	return fib_rec(n1, n2)
}

func main() () : io {
	let _ = fib_rec(1, 1)
}
`
