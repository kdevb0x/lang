package sampleprograms

const UserDefinedType = `
type Foo int

func main() () : io {
	let x Foo = 4
	PrintInt(x)
}`

const EnumType = `
data Foo = A | B

func main() () : io {
	let a Foo = A
	match a {
	case A:
		PrintString("I am A!\n")
	case B:
		PrintString("I am B!\n")
	}
}`

const EnumTypeInferred = `
data Foo = A | B

func main() () : io {
	let a = B
	match a {
	case A:
		PrintString("I am A!\n")
	case B:
		PrintString("I am B!\n")
	}
}`

const GenericEnumType = `
data Maybe a = Nothing | Just a

func DoSomething(x int) (Maybe int) {
	if x > 3 {
		return Nothing
	}
	return Just 5
}

func main() () : io {
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
`
