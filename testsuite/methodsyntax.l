// MethodSyntax tests foo.x() method invocation
// syntax. It should print "10"
func main() () -> affects(IO) {
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
