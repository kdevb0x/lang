(Note: If you've stumbled across this file, it will likely only make sense to @driusan. Ask if you're curious
about something.)

Priorities (in order):
0. Bugs (any that are found preventing the below)
2. Heap allocation (required for both PrintInt and token package)
3. Enforce effects
4. Import (required for ast package)

# Bugfix TODOs
- can not return []byte -- claims "invalid argument"
- Need to enforce that effects are either handled or propagated from function
- need to determine what effect parameters that are mutated should declare
	- ReferenceParameter test
- VM should use PrintByteSlice and PrintString from stdlib, not hack.
- Types need refactoring into another package and enum types need a real type class.

# New features TODOs

- Implement type based function overloading
- Implement "import" / package namespaces
	- refactor builtins into separate standard library package
- Compile time evaluation of pure functions with constant arguments
- Add product type / tuple support
- Add interfaces/polymorphism

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
- Inline function calls
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
- Need way to allocate/return variables on the heap
	- Malloc/Free? GC? Ownership rules?

# Tests TODO
- Add better tests for invalid shadowing or assignments inside conditionals
	- do not allow mutable statement or assignment in conditional
	- ensure variables declared in condition don't outlive scope
- Add test cases for unsigned comparisons in wasm (all operators)
- Add better test cases tail call optimization (esp. with different stack sizes)
- Need better tests for invalid types.. ie typos like "let digits []buf = .." fail for the wrong reasons..
- Need better tests for sum types that aren't passed as functions (ie mutable x string | int, then assign to both string and int )
- Better tests for incompatible tuple values (wrong types, wrong size, access member that doesn't exist, assign to element in mutable tuple, etc.)
- Need tests for tuples in mutable variables and not just let, also as function call parameters

