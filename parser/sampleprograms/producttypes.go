package sampleprograms

const ProductTypeDefn = `
type Foo = (x int, y string)
`

const ProductTypeValue = `
func main () () {
	let x (x int, y bool) = (3, false)
	PrintInt(x.x)
	PrintString("\n")
	PrintInt(x.y)
}
`

const UserProductTypeValue = `
type Foo = (x int, y string)
func main () () {
	let x Foo = (3, "hello\n")
	PrintString(x.y)
	PrintInt(x.x)
}
`
