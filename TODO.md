(Note: If you've stumbled across this file, it will likely only make sense to @driusan. Ask if you're curious
about something.)

# New features TODOs

- Implement "import" / package namespaces
	- refactor builtins into separate standard library package
- Casting (Syntax: `cast(val) as type`)
- Implement multiple dispatch
- Compile time evaluation of pure functions with constant arguments

# Other TODOs

- Add better documentation
	- Add proper documentation to the compiler code base
- Write some non-test sample programs and fix bugs or unergonomic language design (elf linker? autoformatter?)
- Add better test cases tail call optimization (esp. with different stack sizes)

# Syntactic sugar TODOs

- Method invocation syntax. x.foo().bar() should be equivalent to bar(foo(x))
- Add let x match y {} ... syntax for extracting value from sum type and assigning it to a variable?
- Add foreach loop or for maybe just normal for loops. (Syntax needs design.)

# Design TODOs

- Generic functions/macros?
- Investigate and decide on what other types should be implemented:
	- tuples
		- Implement multiple return values from a function.
	- (singly linked) lists?
	- interfaces?
	- structs? (are 2 product types necessary if there's already a tuple?)
	- pointers? (pointers with GC? Just references? Are sum types, generics and a maybe monad enough to not have pointers?)
	- float? dec64?
	- remove "int" and force an explicit size?
- Refactor/rewrite the IR into multiple IRs (at least a forward only/optimizing IR which gets transformed into a better 
  architecture specific IR.)
	- Add optimizations (ie obvious optimizations: dead code elimination, inlining...)
	- Add a new WASM architecture IR, in addition to amd64
