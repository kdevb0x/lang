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
		"if", "else", "else if", "return":
		return true
	}
	return false
}

type Type string

func (t Type) IsValid() bool {
	switch t {
	case "int", "bool", "string":
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
		"=": // assignment
		return true
	}
	return false
}
func (t Operator) String() string {
	return string(t)
}

// BuiltIn functions (only until the language is well-defined enough
// to have a standard library.)
type BuiltIn string

func (bi BuiltIn) IsValid() bool {
	switch bi {
	case "println":
		// prints a builtin type to the screen
		return true
	case "Read":
		// reads an int or string from stdin
		return true
	case "extract":
		// Converts an io value to the underlying type.
		return true
	}
	return false
}

func (t BuiltIn) String() string {
	return string(t)
}

type Char string

func (c Char) IsValid() bool {
	return len(c) == 1
}

func (t Char) String() string {
	return string(t)
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
