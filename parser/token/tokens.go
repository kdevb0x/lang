package token

type Token interface {
	// Returns true iff the token is a valid token according
	// to its type
	IsValid() bool
	String() string
}

type Keyword string

func (t Keyword) String() string {
	return string(t)
}

func (k Keyword) IsValid() bool {
	switch k {
	case "proc", "while", "mut", "let", "func",
		"if", "else", "else if", "return",
		"type", "match", "data", "case":
		return true
	}
	return false
}

type Type string

func (t Type) IsValid() bool {
	switch t {
	case "int", "bool", "string",
		"uint8", "uint16", "uint32", "uint64",
		"int8", "int16", "int32", "int64":
		return true
	}
	return false

}

func (t Type) String() string {
	return string(t)
}

type Operator string

func (o Operator) IsValid() bool {
	switch o {
	case "+", "-", "*", "/", "%", // math
		"<=", "<", "==", ">", ">=", "!=", // comparison
		"=", // assignment
		"|": // other
		return true
	}
	return false
}
func (t Operator) String() string {
	return string(t)
	//return "Operator(" +string(t) + ")"
}

type Char string

func (c Char) IsValid() bool {
	return len(c) == 1
}

func (t Char) String() string {
	return string(t)
	// return "Char(" + string(t) + ")"
}

type String string

func (s String) IsValid() bool {
	return true
}

func (t String) String() string {
	return string(t)
}

type Whitespace string

func (w Whitespace) IsValid() bool {
	// FIXME: Should verify that it's whitespace.
	return true
}
func (t Whitespace) String() string {
	return string(t)
}

type Unknown string

func (i Unknown) IsValid() bool {
	return true
}
func (t Unknown) String() string {
	return string(t)
}
