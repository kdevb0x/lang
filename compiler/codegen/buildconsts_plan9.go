package codegen

const (
	entrypoint = `#define NPRIVATES 16
GLOBL argv(SB), 8+16, $8
GLOBL _tos(SB), 8+16, $8
GLOBL _privates(SB), 8+16, $8
GLOBL _nprivates(SB), 8+16, $4

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
TEXT exits(SB), 20, $-16
	MOVQ $8, BP
	SYSCALL
	RET // Unreached
`

	// wrapper around syscall ssize_t write(int fd, const void *buf, size_t nbytes);
	// Strings are of the format struct{size, [size]char}, so we need to swap
	// the order of the params in the syscall
	write = `
TEXT Write(SB), 20, $0-24
	MOVQ nbytes+8(FP), DX
	MOVQ buf+16(FP), SI

	MOVQ DX, buf+16(FP)
	MOVQ SI, nbytes+8(FP)

	MOVQ $-1, offset+24(FP) // Offset

	MOVQ $51, BP // pwrite syscall
	SYSCALL
	RET
`

	// wrapper around syscall ssize_t write(int fd, const void *buf, size_t nbytes);
	// Strings are of the format struct{size, [size]char}, so we need to swap
	// the order of the params in the syscall
	read = `
TEXT Read(SB), 20, $0-24
	MOVQ $-1, offset+24(FP)
	// MOVQ fd+0(FP), DI

	MOVQ buf+16(FP), SI // buf
	MOVQ len+8(FP), DX // nbytes
	MOVQ DX, nbytes+16(FP)
	MOVQ SI, buf+8(FP)
	MOVQ $50, BP // pread syscall
	SYSCALL
	RET
`

	// wrapper around int open(char *file, int omode)
	// the syscall expects a C string, so we need to make
	// sure the string parameter is null terminated
	open = `
TEXT Open(SB), 20, $0-24
	// Move (the C string portion) into the first arg to the syscall
	MOVQ file+8(FP), DI
	// Move the length into a register, so that we can index by it
	MOVQ len+0(FP), CX
	// Ensure it's nil terminated
	MOVQ $0, file+8(FP)(CX*1)

	MOVQ DI, file+0(FP)
	MOVQ $0, omode+8(FP) // omode = 0 = OREAD
	MOVQ $14, BP // open syscall
	SYSCALL
	RET
`
	// wrapper around int open(char *file, int omode)
	// the syscall expects a C string, so we need to make
	// sure the string parameter is null terminated
	createf = `
TEXT Create(SB), 20, $0-24
	// Move (the C string portion) into the first arg to the syscall
	MOVQ file+8(FP), DI
	// Move the length into a register, so that we can index by it
	MOVQ len+0(FP), CX
	// Ensure it's nil terminated
	MOVQ $0, file+8(FP)(CX*1)

	MOVQ DI, file+0(FP)
	MOVQ $%d, omode+8(FP) // open mode = O_WRONLY|O_CREAT
	MOVQ $438, perms+16(FP) // fileperms. 438 decimal = 0666 octal.
	MOVQ $22, BP // create syscall
	SYSCALL
	RET
`

	closestr = `
TEXT Close(SB), 20, $0-8
//	MOVQ fd+0(FP), DI
	MOVQ $4, BP // close syscall
	SYSCALL
	RET
`
)
