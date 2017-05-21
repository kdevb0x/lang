package sampleprograms

// A simple hello world program using the built-in print.
// It has slightly more complex structure than the HelloWorld, but is
// still fairly straight-forward.
// It has multiple arguments, uses a printf style formatting string and
// duplicates a string
// literal (which should only get inserted into the compiled data section
// once.)
const HelloWorld2 = `proc main() () {
	print("%s %s\n %s", "Hello, world!\n", "World??", "Hello, world!\n")
}`
