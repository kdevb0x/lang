package codegen

// Non-OS specific builtin functions
const (
	// FIXME: This could probably be better written.
	printint = `TEXT PrintInt(SB), 20, $16
	// CX = remaining digits (after div)
	// DX = last digit (after div)
	// DI = pointer to string 
	// R8 = string length
	// R9 = bool true if negative
	MOVQ arg0+0(FP), AX

	// If negative, set R9 and set DI number to the absolute value
	CMPQ AX, $0
	JL neg
	// Positive, set R9 appropriately and begin
	MOVQ $0, R9
	JMP pos

neg:
	NEGQ AX
	MOVQ $1, R9
	// AX is now positive and R9 is true
pos:
	MOVQ $0, R8
	CMPQ AX, $0
	JE print0
	MOVQ $10, CX
	MOVQ SP, DI

div10:
	MOVQ $0, DX
	IDIVQ CX
	ADDQ $48, DX // $48 = '0' char in ASCII
	MOVB DX, (DI)(R8*1)
	INCQ R8
	CMPQ AX, $0
	JE  addsign
	JMP div10

addsign:
	CMPQ R9, $0
	// It was positive, there's no sign to add
	JE reverse
	// Add a "-" sign to the end before reversing the string
	MOVB $45, (DI)(R8*1)
	INCQ R8

reverse:
	PUSHQ R8
	MOVQ SP, R10 // store the location of the string we just built on the stack
	SUBQ $32, SP // make sure we don't overwrite it

	MOVQ R8, R11 // R11 = end char to swap idx
	DECQ R11
	MOVQ $0, R12 // R12 = start char to swap idx
	SARQ $1, R8 // R8 /= 2.. otherwise we'd reverse twice
rloop:
	CMPQ R8, $0
	JE print 

	MOVB (DI)(R11*1), AX
	MOVB (DI)(R12*1), CX
	MOVB AX, (DI)(R12*1)
	MOVB CX, (DI)(R11*1) 
	INCQ R12
	DECQ R11
	DECQ R8
	JMP rloop


print:
	MOVQ R10, 0(SP) // Move from R10 to arg0
	CALL PrintString(SB)
	POPQ AX
	ADDQ $32, SP
	RET
print0:
	MOVQ SP, DI
	MOVB $48, (DI)
	PUSHQ $1
	MOVQ SP, BP
	SUBQ $32, SP
	MOVQ BP, 0(SP)
	CALL PrintString(SB)
	POPQ AX
	ADDQ $32, SP
	RET
`

	printstring = `
TEXT PrintString(SB), 20, $24
	MOVQ $1, 0(SP) // fd
	MOVQ str+0(FP), AX
	MOVQ AX, 8(SP)
	CALL Write(SB)
	RET
`
)
