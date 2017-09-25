package codegen

const (
	entrypoint = `
TEXT	_main(SB), 21, $144
	CALL	main(SB)
loop:
	MOVQ	$0, 0(SP)
	CALL	exits(SB)
	JMP	loop
`

	exits = `
TEXT exits(SB), 20, $0
	MOVQ retcode+0(FP), DI
	MOVQ $0x2000001, AX
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
	MOVQ $0x2000004, AX // write syscall
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
	MOVQ $0x2000003, AX // read syscall
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
	MOVQ $0x2000005, AX // open syscall
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
	MOVQ $0x2000005, AX // open syscall
	SYSCALL
	RET
`

	closestr = `
TEXT Close(SB), 20, $0-8
	MOVQ fd+0(FP), DI
	MOVQ $0x2000006, AX // close syscall
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
	MOVQ $0x2000004, AX // write syscall
	SYSCALL
	RET
`
)
