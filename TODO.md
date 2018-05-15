(Note: If you've stumbled across this file, it will likely only make sense to @driusan. Ask if you're curious
about something.)

# Bugfix TODOs
- can not return []byte -- claims "invalid argument"
- Need to enforce that effects are either handled or propagated from function

# New features TODOs

- Implement multiple dispatch
- Implement "import" / package namespaces
	- refactor builtins into separate standard library package
- Compile time evaluation of pure functions with constant arguments
- Add comments and convert samples from Go strings to files

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

# Syntactic sugar TODOs

- Add foreach loop or for maybe just normal for loops? (Syntax needs design.)

# Design TODOs

- Generic functions/macros?
- Investigate and decide on what other types should be implemented:
	- tuples
		- Implement multiple return values from a function.
	- (singly linked) lists?
	- interfaces?
	- structs? (are 2 product types necessary if there's already a tuple?)
	- float? dec64?

# Tests TODO
- Add better tests for invalid shadowing or assignments inside conditionals
	- do not allow mutable statement or assignment in conditional
	- ensure variables declared in condition don't outlive scope
- Add test cases for unsigned comparisons in wasm (all operators)
- Add better test cases tail call optimization (esp. with different stack sizes)
- Need better tests for invalid types.. ie typos like "let digits []buf = .." fail for the wrong reasons..
