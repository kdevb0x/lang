package sampleprograms

const AssertionFail = `func main () () {
	assert(false)
}`

const AssertionPass = `func main () () {
	assert(true)
}`

const AssertionFailWithMessage = `func main () () {
	assert(false, "This always fails")
}`

const AssertionPassWithMessage = `func main () () {
	assert(true, "You should never see this")
}`

const AssertionFailWithVariable = `func main () () {
	let x = 3
	assert(x > 3)
}`
