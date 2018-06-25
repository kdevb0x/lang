// OutOfOrder defines func which is called before it's
// defined in the source.
func main() () -> affects(IO) {
	PrintInt(foo())
}

func foo() (int) {
	return 3
}
