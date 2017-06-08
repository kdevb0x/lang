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
