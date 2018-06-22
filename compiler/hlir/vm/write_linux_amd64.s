TEXT ·Write(SB), 20, $0-24
	MOVQ fd+0(FP), DI

	// Strings have the format
	// struct{ len int64, buf *char} 
	MOVQ nbytes+8(FP), DX // nbytes
	MOVQ buf+16(FP), SI  // buf
	MOVQ $1, AX // write syscall
	SYSCALL
	RET
