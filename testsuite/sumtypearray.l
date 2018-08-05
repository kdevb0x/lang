type TestType = string | int

func main () () {
	let foo [2]TestType = {
		"stringy",
		33
	}
	assert(len(foo) == 2)
}
