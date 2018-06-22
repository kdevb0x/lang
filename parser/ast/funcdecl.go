package ast

import (
	"fmt"
)

type Callable interface {
	Type
	GetArgs() TupleType
	ReturnTuple() TupleType
}

type TupleType []VarWithType

func (t TupleType) Type() Type {
	if len(t) == 0 {
		return nil
	}
	return t //t[0].Type()
}

func (t TupleType) Info() TypeInfo {
	ti := TypeInfo{}
	for _, piece := range t {
		ti.Size += piece.Type().Info().Size
	}
	return ti
}

func (t TupleType) TypeName() string {
	ret := "("
	for i, piece := range t {
		if i == 0 {
			ret += piece.Type().TypeName()
		} else {
			ret += ", " + piece.Type().TypeName()
		}
	}
	ret += ")"
	return ret
}
func (t TupleType) Node() Node {
	return t
}
func (t TupleType) PrettyPrint(lvl int) string {
	panic("Not implemented")
}
func (t TupleType) Components() []Type {
	var v []Type
	for _, sub := range t {
		v = append(v, sub.Type())
	}
	return v
}

// type Function should be the same as procedure, but
// until the statements are settled we're just have Funcedure
type FuncDecl struct {
	Name    string
	Args    TupleType
	Return  TupleType
	Effects []Effect

	Body BlockStmt
}

func (pd FuncDecl) Node() Node {
	return pd
}

func (pd FuncDecl) GetArgs() TupleType {
	return pd.Args
}

func (fd FuncDecl) String() string {
	return fmt.Sprintf("FuncDecl{\n\tName: %v,\n\tArgs: %v,\n\tReturn: %v,\n\tEffects:: %v\n\tBody: %v}", fd.Name, fd.Args, fd.Return, fd.Effects, fd.Body)
}

func (fd FuncDecl) TypeName() string {
	return fmt.Sprintf("func %v %v", fd.Args.TypeName(), fd.Return.TypeName())
}
func (fd FuncDecl) Type() Type {
	return fd.Return.Type()
}

func (fd FuncDecl) Info() TypeInfo {
	return fd.Return.Info()
}

func (fd FuncDecl) ReturnTuple() TupleType {
	return fd.Return
}

func (fd FuncDecl) PrettyPrint(lvl int) string {
	panic("Unimplemented")
}

func (fd FuncDecl) Components() []Type {
	var v []Type
	for _, t := range fd.Return {
		v = append(v, t.Type())
	}
	return v
}
