package invalidprograms

// BrokenMain is a program consisting of valid tokens, but missing
// a closing bracket. It should result in an error in the AST parser,
// but not in the tokenizer.
const BrokenMain = `proc main() () {
`
