package sampleprograms

const SimpleMatch = `proc main() () {
	let x = 3
	match x {
	case 1:
		print("I am 1\n")
	case 2:
		print("I am 2\n")
	case 4:
		print("I am 4\n")
	case 3:
		print("I am 3\n")
	}
}`

const IfElseMatch = `proc main() () {
	let x = 3
	match {
	case x < 3:
		print("x is less than 3\n")
	case x > 3:
		print("x is greater than 3\n")
	case x < 4:
		print("x is less than 4\n")
	}
}
`
