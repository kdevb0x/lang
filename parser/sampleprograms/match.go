package sampleprograms

const SimpleMatch = `proc main() () {
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

const IfElseMatch = `proc main() () {
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
