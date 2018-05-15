package sampleprograms

const SimpleMatch = `func main () () : io {
	let x = 3
	match x {
	case 1:
		PrintString("I am 1\n")
	case 2:
		PrintString("I am 2\n")
	case 4:
		PrintString("I am 4\n")
	case 3:
		PrintString("I am 3\n")
	}
}`

const IfElseMatch = `func main () () : io {
	let x = 3
	match {
	case x < 3:
		PrintString("x is less than 3\n")
	case x > 3:
		PrintString("x is greater than 3\n")
	case x < 4:
		PrintString("x is less than 4\n")
	}

}
`

const MatchParam = `data Maybe x = Nothing | Just x

func foo (x Maybe int) (int) {
	match x {
	case Just n:
		return n
	case Nothing:
		return 0
	}
}

func main () () : io {
	PrintInt(foo(Just 5))
}`

// Same as above, but print "x".
//
// (There was a bug where func calls didn't work if the string param was a single character long.)
const MatchParam2 = `data Maybe x = Nothing | Just x

func foo (x Maybe int) (int) : io {
	PrintString("x")
	match x {
	case Just n:
		return n
	case Nothing:
		return 0
	}
}

func main () () : io {
	PrintInt(foo(Just 5))
}`
