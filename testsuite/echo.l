func main(args []string) () -> affects(IO) {
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
