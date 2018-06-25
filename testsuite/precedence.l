// Precedence tests that brackets properly adjust the precedence inside of
// arithmetic values.
//
// It should Print "-3"
func main() () -> affects(IO) {
	let x = (1 + 2) * (3 - 4)
	PrintInt(x)
}
