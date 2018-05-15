package sampleprograms

// EmptyMain is the simplest valid program. It contains
// an empty proc main.
const EmptyMain = `func main () () {
}`

// EmptyReturn is the same as EmptyMain, but contains a naked
// return statement.
const EmptyReturn = `proc main() () {
	return
}`
