package codegen

const (
	entrypoint = `#define NPRIVATES	16


GLOBL	argv0(SB), 8+16, $8
GLOBL	_tos(SB), 8+16, $8
GLOBL	_privates(SB), 8+16, $8
GLOBL	_nprivates(SB), 8+16, $4

TEXT	_main(SB), 21, $144
/*
	MOVQ	AX, _tos(SB)
	LEAQ	16(SP), AX
	MOVQ	AX, _privates(SB)
	MOVL	$NPRIVATES, _nprivates(SB)
	MOVL	inargc-8(FP), BP
	LEAQ	inargv+0(FP), AX
	MOVQ	AX, 8(SP)
*/
	CALL	main(SB)

loop:
	MOVQ	$0, 0(SP)
	CALL	exits(SB)
	JMP	loop
`

	exits = `TEXT exits(SB), 20, $0
	// MOVQ BP, a+0(FP)
	MOVQ $8, BP
	SYSCALL
	RET // Unreached
`

	printstring = `// Strings have the format
// struct{ len int64, buf *char} 
TEXT PrintString(SB), 20, $0
	MOVQ str+0(FP), BP
	MOVQ $0, offset+24(FP)
	MOVQ 0(BP), AX
	MOVQ AX, nbytes+16(FP)
	LEAQ 8(BP), AX
	MOVQ AX, buf+8(FP)
	MOVQ $1, fd+0(FP)
	MOVQ $51, BP // pwrite
	SYSCALL
	RET
`
)
