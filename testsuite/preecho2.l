// PreEcho2 is like Echo, but ensures the argument passing of
// slices works correctly.
func PrintSlice(args []string) () -> affects(IO) {
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

func main() () -> affects(IO) {
	let args []string = { "foo", "bar", "baz" }
	PrintSlice(args)
}
