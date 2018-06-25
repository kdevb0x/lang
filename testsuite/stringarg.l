// Stringarg tests passing a string as a parameter.
func main() () ->affects(IO) {
	let b string = "foobar"
	PrintAString(b)
}

func PrintAString(str string) () ->affects(IO) {
	PrintString(str)

}
