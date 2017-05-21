This is just me playing around with creating a toy language, which I
really have no business doing.

It's only online as a backup for myself.  Pay no attention to this
repo.

If you refuse to listen to the above and are curious about this
language, see the `parser/sampleprograms` directory for simple
programs that act as tests.  As currently written, they'll only
compile under Plan9/AMD64, but it shouldn't be too much work to get
the generated assembly files to compile with the Go tool chain instead
of the Plan 9 C toolchain (it'll just result in much bigger binaries,
since you can't tell Go to not link in the Go runtime.)

My main goal with this language is to have a simple, easy to use
language with sum types, a strict delineation of pure and impure
functions, and variables that are immutable by default, and then
experimenting with whatever seems interesting from there.  (I plan on
add compiler time evaluation for pure functions with known arguments
and tail call optimization, too, but right now this code base doesn't
even have the most basic of type systems or operator precedence.)

