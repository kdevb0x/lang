func PrintSecond(args []string) () -> affects(IO) {
	let i = 1
	PrintString(args[i])
}

func main() () -> affects(IO) {
	let aslice []string = {"foo", "bar", "baz" }
	PrintSecond(aslice)
}
