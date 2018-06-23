package sampleprograms

const SumTypeDefn = `
type Foo = int | string
`

const SumTypeFuncCall = `
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
`

const SumTypeFuncReturn = `
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
`

const UserSumTypeDefn = `
// Tests a user defined type which is a sum type containing a user defined enum
enum Keyword = While | Mutable
type Token = Keyword | string
`
