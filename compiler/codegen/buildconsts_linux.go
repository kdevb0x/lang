package codegen

const (
	entrypoint = `TEXT	_main(SB), 21, $144
	CALL	main(SB)

loop:
	MOVQ	$0, 0(SP)
	CALL	exits(SB)
	JMP	loop
`

	exits = `TEXT exits(SB), 20, $0
	MOVQ retcode+0(FP), DI
	MOVQ $60, AX
	SYSCALL
	RET // Unreached
`

	printstring = `// Strings have the format
// struct{ len int64, buf *char} 
TEXT PrintString(SB), 20, $0
	MOVQ str+0(FP), R8
	MOVQ $1, DI // fd
	LEAQ 8(R8), SI // buf
	MOVQ 0(R8), DX // nbytes
	MOVQ $1, AX // write syscall
	SYSCALL
	RET
`
)
