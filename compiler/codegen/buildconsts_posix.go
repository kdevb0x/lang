// +build dragonfly linux darwin

package codegen

const (
	entrypoint = `
TEXT	_main(SB), 21, $-8
	// For some reason that I haven't been able to figure out, 0(SP) is sometimes argc,
	// and 8(SP) is sometimes argc when invoked. It probably has something to do with
	// memory alignment, but for now this hack just detects if the SP is off by 8 and
	// adjusts it when it is.
	CMPQ 0(SP), $0
	JNE hackdone
	ADDQ $8, SP
hackdone:
	// 0(SP) is the number of arguments, including cmd
	// Followed by n pointers to C style strings, followed
	// by a 0.
	// We need to convert these to a []string structure.
	// Slices are structs of the form struct{len, *void}
	// and strings are structs of the form struct{len, [n]char}
	MOVQ SP, R8   // R8 = the original stack pointer that we're converting. It grows up.
	MOVQ SP, BP   // BP = the new stack pointer after moving args. It grows down.
	MOVQ (R8), BX // BX = argc, not mutated. It's the length of the slice at the end.
moveargs:
	ADDQ $8, R8
	CMPQ (R8), $0 // If R8 is 0, we've reached the end of argv
	JE mainstart
	MOVQ (R8), DX // DX = argv[i]. We need to copy it below BP.

//strlenstart:
	MOVQ $0, CX
strlen:
	CMPB (DX)(CX*1), $0
	JE donestrlen
	INCQ CX
	JMP strlen
donestrlen:
	// Copy the string
	// Make room on the (new) stack
	SUBQ CX, BP
	// Align the start of the string.
	ANDQ $~7, BP
	MOVQ CX, R9 // MOVSB is going to destroy CX, so back it up in R9
	MOVQ DX, SI
	MOVQ BP, DI
	CLD
	REP; MOVSB // Copy the string
	
	// Copy the string length.
	SUBQ $8, BP   
	MOVQ R9, 0(BP) // Strlen
	// BP now has the string. Replace the *char at argv[DX] with the string at BP.
	MOVQ BP, (R8)
	SUBQ $24, BP
	JMP moveargs
	
mainstart:
	// Finally, convert SP to a slice, after making room on the stack for it.
	SUBQ $32, BP
	LEAQ 8(SP), R8 // The R8 = the original argv, with the pointers converted from *char to string.
	MOVQ R8, 8(BP)
	MOVQ BX, 0(BP) // BX = argc, still.
	MOVQ BP, SP
	CALL	main(SB)
loop:
	MOVQ	$0, 0(SP)
	CALL	exits(SB)
	JMP	loop

`

	exits = `
TEXT exits(SB), 20, $0
	MOVQ retcode+0(FP), DI
	MOVQ $` + SYS_EXIT + `, AX
	SYSCALL
	RET // Unreached
`

	// wrapper around syscall ssize_t write(int fd, const void *buf, size_t nbytes);
	// Strings are of the format struct{size, [size]char}, so we need to swap
	// the order of the params in the syscall
	write = `
TEXT Write(SB), 20, $0-24
	MOVQ fd+0(FP), DI

	// Strings have the format
	// struct{ len int64, buf *char} 
	MOVQ str+8(FP), R8
	LEAQ 8(R8), SI // buf
	MOVQ 0(R8), DX // nbytes
	MOVQ $` + SYS_WRITE + `, AX // write syscall
	SYSCALL
	RET
`

	// wrapper around syscall ssize_t write(int fd, const void *buf, size_t nbytes);
	// Strings are of the format struct{size, [size]char}, so we need to swap
	// the order of the params in the syscall
	read = `
TEXT Read(SB), 20, $0-24
	MOVQ fd+0(FP), DI

	// Strings have the format
	// struct{ len int64, buf *char} 
	MOVQ buf+16(FP), SI // buf
	MOVQ len+8(FP), DX // nbytes
	MOVQ $0, R10
	MOVQ $0, R8
	MOVQ $0, R9
	MOVQ $` + SYS_READ + `, AX // read syscall
	SYSCALL
	RET
`

	// wrapper around int open(char *file, int omode)
	// the syscall expects a C string, so we need to make
	// sure the string parameter is null terminated
	open = `
TEXT Open(SB), 20, $0-24
	// 0(FP) is the string, which has the format struct{n int, buf [n]byte}
	MOVQ file+0(FP), BX
	// Move (the C string portion) into the first arg to the syscall
	LEAQ 8(BX), DI
	// Move the length into a register, so that we can index by it
	MOVQ 0(BX), CX
	// Ensure the string is null terminated.
	// FIXME: This is segfaulting.
	// MOVB $0, (DI)(CX*1)
	MOVQ $0, SI // open mode = 0 = O_RDONLY
	MOVQ $0, DX // fileperms = irrelevant, since it's read only.. 
	MOVQ $` + SYS_OPEN + `, AX // open syscall
	SYSCALL
	RET
`
	// wrapper around int open(char *file, int omode)
	// the syscall expects a C string, so we need to make
	// sure the string parameter is null terminated
	createf = `
TEXT Create(SB), 20, $0-24
	// 0(FP) is the string, which has the format struct{n int, buf [n]byte}
	MOVQ file+0(FP), BX
	// Move (the C string portion) into the first arg to the syscall
	LEAQ 8(BX), DI
	// Move the length into a register, so that we can index by it
	MOVQ 0(BX), CX
	// Ensure the string is null terminated.
	// FIXME: This is segfaulting.
	// MOVB $0, (DI)(CX*1)
	MOVQ $%d, SI // open mode = O_WRONLY|O_CREAT
	MOVQ $438, DX // fileperms. 438 decimal = 0666 octal.
	MOVQ $` + SYS_OPEN + `, AX // open syscall
	SYSCALL
	RET
`

	closestr = `
TEXT Close(SB), 20, $0-8
	MOVQ fd+0(FP), DI
	MOVQ $` + SYS_CLOSE + `, AX // close syscall
	SYSCALL
	RET
`

	// FIXME: This should just be a wrapper to PrintString(), but
	// for some reason it's not working unless it's inlined..
	printbyteslice = `
TEXT PrintByteSlice(SB), 20, $16
	// wrapper around
	// write(1, *buf, nbytes) syscall.
	// Byte slices have the format
	// struct{ len int64, buf *char}, the inverse
	// of what we want
	MOVQ $1, DI // fd
	MOVQ buf+8(FP), SI // buf
	MOVQ len+0(FP), DX // nbytes
	MOVQ $` + SYS_WRITE + `, AX // write syscall
	SYSCALL
	RET
`
)
