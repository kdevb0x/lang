package codegen

// Non-OS specific builtin functions
const (
	// FIXME: This could probably be better written.
	printint = `TEXT PrintInt(SB), 20, $32-16
	// CX = remaining digits (after div)
	// DX = last digit (after div)
	// DI = pointer to string 
	// R8 = string length
	// R9 = bool true if negative
	// R13 = unmodified string length
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
	INCQ R13
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
	MOVQ R8, R13
	CMPQ R9, $0
	// It was positive, there's no sign to add
	JE reverse
	// Add a "-" sign to the end before reversing the string
	MOVB $45, (DI)(R8*1)
	INCQ R8
	INCQ R13
reverse:
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
	MOVQ R13, 0(SP)
	MOVQ R10, 8(SP) // Move from R10 to arg0
	CALL PrintString(SB)
	ADDQ $32, SP
	RET
print0:
	MOVQ SP, DI
	MOVB $48, (DI)
	SUBQ $24, SP
	MOVQ $1, 0(SP)
	MOVQ DI, 8(SP)
	CALL PrintString(SB)
	ADDQ $24, SP
	RET
`

	slicelen = `
TEXT len(SB), 20, $0-16
	MOVQ len+0(FP), AX
	RET
`
)
