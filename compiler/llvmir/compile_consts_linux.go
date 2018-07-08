package codegen

const (
	O_CREAT  = 0100
	O_WRONLY = 1
)

var CREATE_CONST = fmt.Sprintf("%d", O_CREAT|O_WRONLY)

const (
	SYS_EXIT  = "60"
	SYS_WRITE = "1"
	SYS_READ  = "0"
	SYS_OPEN  = "2"
	SYS_CLOSE = "3"
)
