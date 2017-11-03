(Note: If you've stumbled across this file, it will likely only make sense to @driusan. Ask if you're curious
about something.)

# Bugfix TODOs
- functions should reserve the correct stack size
- len should work on all slice types, and arrays, and strings
- add test cases for unsigned comparisons in wasm (all operators)

# New features TODOs

- Implement multiple dispatch
- Implement "import" / package namespaces
	- refactor builtins into separate standard library package
- Casting (Syntax: `cast(val) as type`)
- Compile time evaluation of pure functions with constant arguments
- Allow let statements in conditions to bind to appropriate scope (ie while (let x = Read()) > 0 { ... }

# Other TODOs

- Add better documentation
	- Add proper documentation to the compiler code base
- Write some non-test sample programs and fix bugs or unergonomic language design (elf linker? autoformatter?)
- Add better test cases tail call optimization (esp. with different stack sizes)

# Syntactic sugar TODOs

- Method invocation syntax. x.foo().bar() should be equivalent to bar(foo(x))
- Add foreach loop or for maybe just normal for loops. (Syntax needs design.)

# Design TODOs

- Generic functions/macros?
- Investigate and decide on what other types should be implemented:
	- tuples
		- Implement multiple return values from a function.
	- (singly linked) lists?
	- interfaces?
	- structs? (are 2 product types necessary if there's already a tuple?)
	- float? dec64?
- Refactor/rewrite the IR into multiple IRs (at least a forward only/optimizing IR which gets transformed into a better 
  architecture specific IR.)
	- Add optimizations (ie obvious optimizations: dead code elimination, inlining...)
