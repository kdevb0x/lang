func PrintSecond(args []string) () -> affects(IO) {
	PrintString(args[1])
}

func main() () -> affects(IO) {
	let aslice []string = {"foo", "bar", "baz" }
	PrintSecond(aslice)
}
