func foo(x bool) (int | string) {
	if x {
		return 3
	}
	return "not3"
}

func main () () {
	let x = foo(false)
	match x {
	case int:
		PrintInt(x)
	case string:
		PrintString(x)
	}

	let x = foo(true)
	match x {
	case int:
		PrintInt(x)
	case string:
		PrintString(x)
	}
}
