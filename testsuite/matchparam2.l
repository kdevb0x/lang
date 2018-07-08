enum Maybe x = Nothing | Just x

func foo (x Maybe int64) (int) -> affects(IO) {
	PrintString("x")
	match x {
	case Just n:
		return n
	case Nothing:
		return 0
	}
}

func main () () -> affects(IO) {
	PrintInt(foo(Just 5))
}
