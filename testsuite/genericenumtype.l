enum Maybe a = Nothing | Just a

func DoSomething(x int) (Maybe int) {
	if x > 3 {
		return Nothing
	}
	return Just 5
}

func main() () -> affects(IO) {
	let x = DoSomething(3)
	match x {
	case Nothing:
		PrintString("I am nothing!\n")
	case Just n:
		PrintInt(n)
		PrintString("\n")
	}
	let x = DoSomething(4)
	match x {
	case Nothing:
		PrintString("I am nothing!\n")
	case Just n:
		PrintInt(n)
		PrintString("\n")
	}
}

