(Note: If you've stumbled across this file, it will likely only make sense to @driusan. Ask if you're curious
about something.)

# Bugfix TODOs
- len should work on all slice types, and arrays, and strings
- add test cases for unsigned comparisons in wasm (all operators)

# New features TODOs

- Implement multiple dispatch
- Implement "import" / package namespaces
	- refactor builtins into separate standard library package
- Casting (Syntax: `cast(val) as type`)
- Compile time evaluation of pure functions with constant arguments
- Allow let statements in conditions to bind to appropriate scope (ie while (let x = Read()) > 0 { ... }

# HLIR Optimization TODOs
- Add optimization pass
- Arithmetic simplification for constants
- Comparison simplification for constants
- Constant propagation
- Eliminate unused instructions / blocks
	- Dst unused:
		- MOV
		- ADD, SUB, DIV, MUL, MOD
		- EQ, NEQ, GEQ, GT, LT, LTE
	- if(false) { .. } else { ..2 } => ..2
	- if(true) { .. } else { ..2} => ..
	- while(false) { .. } => eliminate
- Pure function evaluation for constant arguments
- Inline functions
	- Small functions
	- Functions with at least 1 constant argument

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
