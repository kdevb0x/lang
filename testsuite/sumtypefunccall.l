func foo (x int | string) () -> affects (IO) {
	match x {
	case int:
		PrintInt(x)
	case string:
		PrintString(x)
	}
}

func main () () {
	foo("bar")
	foo(3)
}
