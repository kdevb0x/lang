# Language Spec

This specifies the unnamed language that has a compiler in this directory. The
goal of the language is to be a relatively simple compiled language with small
binaries, sum types, and an explicit delineation between pure and impure
functions.

## Variables

Variables are typed and come in 2 variations: mutable, and immutable.

Immutable variables are defined with the `let` keyword while mutable
ones are defined with the `mutable` keyword.

Both take a name, an optional type, an equal sign, and then the initial
value such as:

```
let x int = 0
```

or

```
mutable x int = 0
```

The type ("int" in the above) is optional, and if left out will be inferred
from the value being assigned to it. Both `let` and `mutable` statements
require a value to be initialized to.

let and mutable variables have the following differences:

1. mutable variables can be assigned a new value with the `=` operator, but
   can never be shadowed by another variable of the same name (even in a
   different scope).
2. let variables can be shadowed by a new variable with the same name, but
   the value of the variable can never change. (Shadowing a variable declares
   a new variable which just happens to have the same name.)

For instance, this is legal, and creates a variable x, then changes its value:

```
mutable x = 3
x = 5
```

while this is not, because x is immutable:

```
let x = 3
x = 5
```

Similarly, this is legal, and declares 2 variables with the same name (
although in this case, only the string variation can ever be referenced..):

```
let x = 3
let x = "string"
```

while this is not, because it attempts to redeclare a mutable variable x:

```
mutable x = 3
mutable x = "string"
```

mutable variables can *never* be shadowed, even by let variables or in a
different scope.

This is illegal:

```
mutable x = 3
let x = 4
```

as the mutable x can not be shadowed.

As is this, even though the if block is a different scope (variables
are block scoped):

```
mutable x = 3
if x == 3 {
	let x = 4
}
```

However, this is legal:

```
let x = 3
if x == 3 {
	mutable x = 4
	x = 5
}
```

Since the initial x can be shadowed by the mutable x, and then the mutable
x can be modified.

## Types

Variables are typed (even if the type is usually inferred while declaring
a variable.)

Built in types are int (system word size, currently always 8 bytes), uint,
uint8/byte, int8, int16, int32, int64, uint16, uint32, uint64, string and bool.

New user types can be defined in 2 ways. The simplest way is to define a new
type as a different kind of some existing type, using the "type" keyword  such
as in:

```
type NewType OldType
```

To define NewType as a kind of OldType.

A more concrete example might be:

```
type Password string
```

To define Password as a type of string. Types defined in this way can not be
used interchangably in variable assignments or function parameters, but creates
a new type which has the same memory layout and characteristics as old type.
It can be thought of as roughly analogous to the "type" keyword in Go. 

The second way to define a type, is to use the "data" keyword to create an
enumerated (sum) type or a type of generic container. The "data" keyword is
more closely inspired by Haskell syntax (but isn't as robust or powerful.)

You can create a simple enumeration by separating the options with a `|`
character. For instance:

```
data StopLight = Red | Yellow | Green
```

This introduces StopLight as a new type, which can have the values of either
`Red`, `Yellow`, or `Green`. A variable of type StopLight can be assigned the
value by using the barewords introduced above. For instance:

```
mutable NearMyHouse StopLight = Red
```

If there are any identifiers before the equals sign in the data declaration,
then it creates a generic family of types where the identifiers must be
concretized for instances of the type. For instance, you can create a family
of option types with:

```
data Maybe x = Nothing | Just x
```

When a variable is defined, the "x" must be specified. For instance, based
on the above you could create an optional int type as:

```
let x Maybe int = Just 3
```

When declaring the variable, the type was concretized from the generic "Maybe x"
to a concrete type of "Maybe int" for the variable "x"

### Arrays and Slices

Variables can also be arrays, or slices of arrays. Arrays are declared with
a Go-like syntax of prefixing the type with "[size]", while array literals are
created by surrounding the values with `{` and `}`. For instance to declare
an array of three bytes, containing the values of 0, 1, and 2:

```
let firstthree [3]byte = { 0, 1, 2 }
```

Slices are like arrays, but also store the size of the backing array in memory.
They can be thought of as variable size arrays in this language. They're created
by omitting the size in the type:

```
let somevalues []byte = { 0, 1, 2, 4 }
```

N.B.: Slices are implemented differently than in most languages. This will
likely be eventually fixed, but for now they mostly exist to solve the problem
of allowing variable-sized arrays to be passed to functions such as `Read()`
They also lack the ability to do some obvious things like get the length.

## Functions

Functions, like most things, come in two varieties: a pure, and an impure form.

The pure form is declared with the `func` keyword, while the impure form is
declared with the `proc` keyword. If a proc with the name `main` is defined,
it's used as the entry point to the program. procs can call funcs, but funcs
can never call a proc. The tradeoff is that funcs can be theoretically evaluted
at compile time if the arguments are known at compile time. (N.B. Compile time
evaluation is planned, but not implemented in the reference compiler.)

Functions of both types take a tuple of arguments, and a tuple of return types.
(N.B. The syntax supports multiple return values, but support is currently only
implemented for a single return value in the reference compiler.)
Both tuples are required, even if empty. The function signature is followed by a
code block enclosed in curly brackets. The return value is specified with the
return keyword and is required (unless the function returns the empty tuple.)

For instance, to create a function named "three" that returns 3:

```
func three () (int) {
	return 3
}
```

To create a main procedure which calls the (builtin, see below) PrintString
procedure.

```
proc main () () {
	PrintString("Hello, world!")
}
```

Or putting it together, we can call a function from a procedure (and throw
an argument into the mix for fun).

```
proc main () () {
	PrintInt(threemore(5))
}

func threemore (x int) (int) {
	return x + 3
}
```

(Note that in the above there was need to no forward declare `threemore` before
calling it from main.)

### Reference parameters (mutable keyword)

While arguments are generally passed by value, procs (but not funcs) can
declare a parameter passed to them as a reference by prefixing the variable with
the "mutable" keyword. (References should generally only be used for things that
are, in fact, mutated, so there's no way to declare something as a non-mutated
reference.) If a parameter is a reference parameter, only mutable variables
can be passed to that parameter.

## Control Flow

### Valid comparison operators

The following comparison operators are defined for all types: `==`, `!=`.

Something is equal if it has the same value and is of the same type. Something
is not equal if equals is false.

The following comparison operators are defined for variables of the same integer
base type, and have the obvious meanings: `<=`, `<`, `==`, `!=`, `>`, `>=`. They
are used infix and return a bool. Their behaviour is undefined for non-integer
or incompatible types.

### if 

If statements are defined with the `if` keyword and take something that
evaluates to a `bool` type as a parameter. Behaviour is (currently) undefined
for other types being passed (but will likely turn into an error in the future.)

The condition for an `if` statement is followed by a block in curly brackets
(and must be a block, not a statement.) The code block can optionally be
followed by an `else` keyword followed by either a code block or another `if`
statement to create an `else if` chain.

Some examples:

```
if true {
	PrintString("true\n")
}

let x = 3
if x >= 2 {
	PrintString(">= 2\n")
} else {
	PrintString("< 2\n")
}

if x == 3 {
	PrintString("x == 3\n")
} else if x == 2 {
	PrintString("x == 2\n")
} else {
	PrintString("x is neither 3, nor 2!\n")
}
```

### while

`while` statements take a `bool` argument and a code block and repeat that block
until the argument is false. 

For instance, the following will print "543210"

```
mutable x = 5
while x >= 0 {
	PrintInt(x)
	x = x - 1
}
```

### match

The match keyword is somewhat akin to the `switch` keyword in Go (without the
ability to combine cases with a comma or take an initializer), but has the added
constraint that if the argument passed to it is an enumerated type defined with
the `data` keyword, the cases must be exhaustive of all options, otherwise it's
a compile time error.

In the simplest case, you give it a variable as an argument and check the cases
against the possible results. (There is no fallthrough)

The following will print "Three!":

```
let x = 3
match x {
case 2: PrintString("Two!")
case 3: PrintString("Three!")
default: PrintString("Not 2 or 3!")
}
```

The value before the case list is optional, and can be omitted to use `match`
as an else if block chain (with prettier indentation.)

For instance, to normalize a variable x to -1, 0, or 1 (perhaps for some sort
of sorting algorithm):

```
match {
case x < 0:
	return -1
case x > 1:
	return 1
default:
	return 0
}
```

More idiomatically, however, would be to return an enumerated like
`data Ord = EQ | LT | GT`, since matching against `data` types require all
possible options to be covered.

The final use of the match statement is to extract the value from a generic
`data` container type. The easiest way to describe this is probably with an
example.

If we were to use the `Maybe int` type defined above, and pass it to a function,
we can use match to get the "int" portion out.

```
data Maybe x = Nothing | Just x
func Example(x Maybe int) (string) {
	match x {
	case Nothing:
		return "Oh no!"
	case Just y:
		if y > 0 {
			return "Greater than zero"
		} else if y < 0 {
			return "Less than zero"
		else {
			return "Zero"
		}
	}
}
```

In the above example, "y" is a valid variable inside of the match case covering
the "Just" branch, and inside of the case it refers to the value inside of x,
while "x" itself has the value of "Just x". If we were to call the function as
`Example(Just 3)` `y` would be `3` and have type int, while x would be `Just 3`
and have type `Maybe int`.

## Builtins

There are a number of builtins which are primarily intended to support the test
suite by making it possible to do some kind of I/O that we can test against.
There are builtins for printing, and file I/O. They will likely be moved outside
of the language and into a separate standard library once some kind of
namespacing/packaging system is implemented.

All of the builtins are procs, since they all have side-effects.

### PrintInt, PrintString, PrintByteSlice

`PrintInt` prints int types to stdout, `PrintString` prints string types, and
`PrintByteSlice` prints byte slices (interpreted as a string.)

### Open, Create, Read, Write, Close

`Open` opens a file for reading and returns a file descriptor. It has the
signature `proc Open (file string) (int)`

`Create` opens a file for writing, creating it with mode 0666 (before umask)
and truncating it if it exists. It has the signature
`proc Create(file string) (int)`

`Read` takes a file descriptor and a mutable byte slice, and reads the length
of the slice into it, overwriting its previous contents. It returns the number
of bytes read from the file descriptor, and 0 should be interpreted as EOF.
It has the signature `proc Read(fd int, mutable data []byte) (int)`

`Write` takes a file descriptor and a byte slice, and writes the contents of
the the byte slice into the file. It returns the number of bytes written. It
has the signature `proc Write(fd int, data []byte) (int)`

`Close` closes an open file descriptor. It has the signature:
`proc Close (fd int) ()`

N.B. You can currently only Write to files opened with Create, and only Read
files opened with Open. There is no more generic way to open a file read/write,
or to open a file for writing without truncating it (yet.)
