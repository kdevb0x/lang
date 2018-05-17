package ast

type Context struct {
	Variables   map[string]VarWithType
	Mutables    map[string]VarWithType
	Functions   map[string]Callable
	Types       map[string]TypeDefn
	PureContext bool // true if inside a pure function.
	EnumOptions map[string]EnumOption
	CurFunc     Callable
}

func NewContext() Context {
	return Context{
		Variables: make(map[string]VarWithType),
		Functions: map[string]Callable{
			// FIXME: These should be replaced by a
			// multiple dispatch Print function and/or
			// moved to a standard library, not build in.
			"PrintString": FuncDecl{
				Name: "PrintString",
				Args: []VarWithType{
					{"str", TypeLiteral("string"), false},
				},
				Effects: []Effect{"IO"},
			},
			"PrintInt": FuncDecl{
				Name: "PrintInt",
				Args: []VarWithType{
					{"x", TypeLiteral("int"), false},
				},
				Effects: []Effect{"IO"},
			},
			"PrintByteSlice": FuncDecl{
				Name: "PrintByteSlice",
				Args: []VarWithType{
					{"slice", SliceType{TypeLiteral("byte")}, false},
				},
				Effects: []Effect{"IO"},
			},
			"len": FuncDecl{
				Name: "len",
				Args: []VarWithType{
					// FIXME: This should be any slice, not just string slices, but
					// string slices take priority until there's some sort of generic func
					// decl, because they're passed to main..
					{"slice", SliceType{TypeLiteral("string")}, false},
				},
				Return: []VarWithType{
					{"", TypeLiteral("uint64"), false},
				},
			},
			// FIXME: These should be moved out of the compiler
			// and into a standard library, once enough of the
			// compiler is implemented to have a standard
			// library.
			"Write": FuncDecl{
				Name: "Write",
				Args: []VarWithType{
					{"fd", TypeLiteral("uint64"), false},
					{"val", SliceType{TypeLiteral("byte")}, false},
				},
				Effects: []Effect{"IO", "Filesystem"},
			},
			"Read": FuncDecl{
				Args: []VarWithType{
					{"fd", TypeLiteral("uint64"), false},
					// NB. this should be []byte, once arrays are implemented.
					// NB2. This will read exactly the length of string from fd int
					//     dst and overwrite what's there.
					//     It needs a way to mark parameters mutable.
					{"dst", SliceType{TypeLiteral("byte")}, true},
				},
				Return: []VarWithType{{"", TypeLiteral("uint64"), false}},
				Effects: []Effect{"Filesystem"},
			},
			"Open": FuncDecl{
				Name: "Open",
				Effects: []Effect{"FD"},
				Args: []VarWithType{
					{"val", TypeLiteral("string"), false},
					/*
						Follow Go/Plan9 conventions. Open just opens, Create
						just creates with a default umask, rather than having
						options on a generic Open like POSIX does.
						This makes it easier to port to Plan 9.
						{"mode", TypeLiteral("int"), false},
						{"createperms", TypeLiteral("int"), false},
					*/
				},
				Return: []VarWithType{{"", TypeLiteral("uint64"), false}},
			},
			"Create": FuncDecl{
				Name: "Create",
				Effects: []Effect{"FD", "Filesystem"},
				Args: []VarWithType{
					{"val", TypeLiteral("string"), false},
					/*
						Follow Go/Plan9 conventions. Open just opens, Create
						just creates with a default umask, rather than having
						options on a generic Open like POSIX does.
						This makes it easier to port to Plan 9.
						{"mode", TypeLiteral("int"), false},
						{"createperms", TypeLiteral("int"), false},
					*/
				},
				Return: []VarWithType{{"", TypeLiteral("uint64"), false}},
			},
			"Close": FuncDecl{
				Name: "Close",
				Effects: []Effect{"FD"},
				Args: []VarWithType{
					{"val", TypeLiteral("uint64"), false},
				},
			},
		},
		Mutables: make(map[string]VarWithType),
		Types: map[string]TypeDefn{
			"int":    TypeDefn{TypeLiteral("int"), TypeLiteral("int"), nil},
			"uint":   TypeDefn{TypeLiteral("uint"), TypeLiteral("uint"), nil},
			"uint8":  TypeDefn{TypeLiteral("uint8"), TypeLiteral("uint8"), nil},
			"byte":   TypeDefn{TypeLiteral("byte"), TypeLiteral("byte"), nil},
			"uint16": TypeDefn{TypeLiteral("uint16"), TypeLiteral("uint16"), nil},
			"uint32": TypeDefn{TypeLiteral("uint32"), TypeLiteral("uint32"), nil},
			"uint64": TypeDefn{TypeLiteral("uint64"), TypeLiteral("uint64"), nil},
			"int8":   TypeDefn{TypeLiteral("int8"), TypeLiteral("int8"), nil},
			"int16":  TypeDefn{TypeLiteral("int16"), TypeLiteral("int16"), nil},
			"int32":  TypeDefn{TypeLiteral("int32"), TypeLiteral("int32"), nil},
			"int64":  TypeDefn{TypeLiteral("int64"), TypeLiteral("int64"), nil},
			"string": TypeDefn{TypeLiteral("string"), TypeLiteral("string"), nil},
			"bool":   TypeDefn{TypeLiteral("bool"), TypeLiteral("bool"), nil},
		},
		PureContext: false,
		EnumOptions: make(map[string]EnumOption),
	}
}

func (c Context) Clone() Context {
	var c2 Context
	c2.Variables = make(map[string]VarWithType)
	c2.Functions = make(map[string]Callable)
	c2.Mutables = make(map[string]VarWithType)
	c2.Types = make(map[string]TypeDefn)
	c2.EnumOptions = make(map[string]EnumOption)
	for k, v := range c.Variables {
		c2.Variables[k] = v
	}
	for k, v := range c.Mutables {
		c2.Mutables[k] = v
	}

	for k, v := range c.Functions {
		c2.Functions[k] = v
	}
	for k, v := range c.Types {
		c2.Types[k] = v
	}
	for k, v := range c.EnumOptions {
		c2.EnumOptions[k] = v
	}
	c2.PureContext = c.PureContext
	c2.CurFunc = c.CurFunc
	return c2
}

func (c Context) IsVariable(s string) bool {
	for _, v := range c.Variables {
		if string(v.Name) == s {
			return true
		}
	}
	return false
}

func (c Context) IsMutable(s string) bool {
	for _, v := range c.Mutables {
		if string(v.Name) == s {
			return true
		}
	}
	return false
}

func (c Context) IsFunction(s string) bool {
	for k := range c.Functions {
		if k == s {
			return true
		}
	}
	return false
}

func (c Context) ValidType(t Type) bool {
	if _, ok := c.Types[t.Type()]; ok {
		return true
	}
	return false
}

func (c Context) EnumeratedOption(t string) *EnumOption {
	if eo, ok := c.EnumOptions[t]; ok {
		return &eo
	}
	return nil
}
