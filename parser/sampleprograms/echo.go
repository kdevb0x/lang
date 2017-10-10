package sampleprograms

// echo echos its arguments to stdout. It's the simplest
// program to test command line parameters..
const Echo = `proc main(args []string) () {
	mutable i = 1
	let length = len(args)
	while i < length {
		PrintString(args[i])

		i = i + 1

		if i != length {
			PrintString(" ")
		}
	}
	PrintString("\n")
}`

// PreEcho is like Echo, but it doesn't take command line
// arguments and has the parameters hardcoded, to make sure
// any bugs in echo are from the argument passing, not the
// program logic.
const PreEcho = `proc main() () {
	let args []string = { "foo", "bar", "baz" }
	mutable i = 1
	let length = len(args)
	while i < length {
		PrintString(args[i])

		i = i + 1

		if i != length {
			PrintString(" ")
		}
	}
	PrintString("\n")
}`

// PreEcho2 is like Echo, but ensures the argument passing of
// slices works correctly.
const PreEcho2 = `proc PrintSlice(args []string) () {
	mutable i = 1
	let length = len(args)
	while i < length {
		PrintString(args[i])

		i = i + 1

		if i != length {
			PrintString(" ")
		}
	}
	PrintString("\n")
}

proc main() () {
	let args []string = { "foo", "bar", "baz" }
	PrintSlice(args)
}
`
