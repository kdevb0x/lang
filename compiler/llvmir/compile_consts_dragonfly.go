package llvmir

import "fmt"

const (
	O_CREAT  = 0x200
	O_WRONLY = 1
)

var CREATE_CONST = fmt.Sprintf("%d", O_CREAT|O_WRONLY)

const (
	SYS_EXIT  = "1"
	SYS_READ  = "3"
	SYS_WRITE = "4"
	SYS_OPEN  = "5"
	SYS_CLOSE = "6"
)
