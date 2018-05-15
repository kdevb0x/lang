package sampleprograms

// MethodSyntax tests foo.x() method invocation
// syntax. It should print "10"
const MethodSyntax = `func main() () : io {
	let foo = 3
	let y = foo.add3().add(4)
	PrintInt(y)
}

func add3(val int) (int) {
	return val + 3
}

func add(x int, y int) (int) {
	return x + y
}
`
