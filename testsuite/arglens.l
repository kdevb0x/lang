func main(args []string) () -> affects(IO) {
	mutable i = 1
	while i < args.len() {
		PrintInt(args[i].len())
		if i != len(args) - 1 {
			PrintString(" ")
		}

		i = i + 1
	}
	PrintString("\n")
}
