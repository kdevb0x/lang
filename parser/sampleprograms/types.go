package sampleprograms

const UserDefinedType = `
type Foo int

proc main() () {
	let x Foo = 4
	print("%d\n", x)
}`

const EnumType = `
data Foo = A | B

proc main() () {
	let a Foo = A
	match a {
	case A:
		print("I am A!\n")
	case B:
		print("I am B!\n")
	}
}`

const EnumTypeInferred = `
data Foo = A | B

proc main() () {
	let a = B
	match a {
	case A:
		print("I am A!\n")
	case B:
		print("I am B!\n")
	}
}`

const GenericEnumType = `
data Maybe a = Nothing | Just a

func DoSomething(x int) (Maybe int) {
	if x > 3 {
		return Nothing
	}
	return Just 3
}

proc main() () {
	let x = DoSomething(3)
	match x {
	case Nothing:
		print("I am nothing!")
	case Just n:
		print("%d", n)
	}
	let x = DoSomething(4)
	match x {
	case Nothing:
		print("I am nothing!")
	case Just n:
		print("%d", n)
	}
}
`
