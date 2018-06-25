// PreEcho is like Echo, but it doesn't take command line
// arguments and has the parameters hardcoded, to make sure
// any bugs in echo are from the argument passing, not the
// program logic.
func main() () -> affects(IO) {
	let args []string = { "foo", "bar", "baz" }
	mutable i = 1
	let length = len(args)
	while i < length {
		PrintString(args[i])

		i = i + 1

		if i != length {
			PrintString(" ")
		}
	}
	PrintString("\n")
}
