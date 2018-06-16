package ast

func nTabs(lvl int) string {
	var ret string
	for i := 0; i < lvl; i++ {
		ret += "\t"
	}
	return ret
}
