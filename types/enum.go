package types

import (
	"github.com/driusan/lang/parser/ast"
)

func GetEnumMap(nodes []ast.Node) (EnumMap) {
	rv := make(EnumMap)
	for _, n := range nodes {
		switch f := n.(type) {
		case ast.EnumTypeDefn:
			for i, v := range f.Options {
				rv[v.Constructor] = EnumIndexInfo{i, f}
			}
		}
	}
	return rv
}
type EnumIndexInfo struct {
	Index int
	Defn  ast.EnumTypeDefn
}

type EnumMap map[string]EnumIndexInfo

func (em EnumMap) GetIndex(val ast.EnumOption) int {
	for key, info := range em {
		if val.Constructor == key {
			return info.Index
		}
	}
	// There should be no way for an EnumValue to be in this situation
	panic("Could not get index for enum value")
}
