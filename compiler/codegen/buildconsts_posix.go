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
	// and strings are structs of the form struct{len, *void}
	//
	// SP, R8, and BP all start out as the current stack pointer.
	// SP does not get modified until the end, BP is a new stack pointer
	// that grows downwards as we calculate the strings, and R8 is the
	// index into the C style argv[][] that we're currently working with.

	MOVQ SP, BP   // BP = the new stack pointer after moving args. It grows down.
	MOVQ (SP), BX // BX = argc, not mutated. It's the length of the slice at the end.
	LEAQ (SP)(BX*8), R8 // R8 = the top of the C style *char[], since we're growing BP
			// downwards, we need to ensure that we convert the arguments in reverse,
			// and start at the end of argv.

	// Loop over the arguments
moveargs:
	CMPQ R8, SP // If R8 is the original SP, we've reached the end of the loop.
	JE mainstart
	MOVQ (R8), DX // DX = argv[i]. We need to store the pointer for after we've calculated
		      // the length.

//strlenstart:
	MOVQ $0, CX // CX = the string length. Starts at 0.
strlen:
	CMPB (DX)(CX*1), $0 // If DX[i] (where i is the current length) is 0, we've reached the end
			// of the string.
	JE donestrlen
	INCQ CX	// Increment length and check again.
	JMP strlen
donestrlen:
	// The string length is now in CX, and the pointer in DX.

	// Copy the string as struct{n, *char} to below BP
	// String pointer
	SUBQ $8, BP
	MOVQ DX, 0(BP)
	// String length.
	SUBQ $8, BP
	MOVQ CX, 0(BP)

	SUBQ $8, R8 // Move to the next argument
	JMP moveargs

mainstart:
	// BP is not the start of a slice of strings, but the slice header is missing.
	// BX is the slice length
	// 
	// Add the slice header for the []string. The pointer to the first string is currently at
	// BP, so temporarily store it in AX because we're about to modify BP.
	MOVQ BP, AX
	SUBQ $24, BP
	MOVQ AX, 8(BP) // Store the base pointer
	MOVQ BX, 0(BP) // Store the length

	// Make BP the new stack pointer for main now that we're done playing around with the
	// C style (int argc, char **argv) header that was in memory and converted it to a
	// []string
	MOVQ BP, SP
	CALL	main(SB)
loop:
	// If we fall through main, exit with a "success" exit status.
	MOVQ	$0, 0(SP)
	CALL	exits(SB)
	// Unreachable, but if it happens, just loop back keep calling exit.
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
	MOVQ nbytes+8(FP), DX // nbytes
	MOVQ buf+16(FP), SI  // buf
	MOVQ $` + SYS_WRITE + `, AX // write syscall
	SYSCALL
	RET
`

	// wrapper around syscall ssize_t read(int fd, const void *buf, size_t nbytes);
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
	// Move the C string pointer to the first register for the syscall
	MOVQ file+8(FP), DI
	// Move the length into any register, so that we can index by it
	MOVQ file+0(FP), CX
	// Ensure the string is null terminated.
	MOVB $0, file+0(FP)(CX*1)
	// Hardcoded parameters for the open syscall.
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
	// Move the C string pointer to the first register for the syscall
	MOVQ file+8(FP), DI
	// Move the length into any register, so that we can index by it
	MOVQ file+0(FP), CX
	// Ensure the string is null terminated.
	MOVB $0, file+0(FP)(CX*1)
	// Hardcoded parameters for the open syscall.
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
)
