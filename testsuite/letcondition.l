// LetCondition tests creating let statements in a conditional.
func main() () -> affects(IO) {
	let i = 0
	if (let i = i + 1) == 1 {
		PrintInt(i)
	} else {
		PrintInt(-1)
	}

	if (let i = i + 1) != 1 {
		PrintInt(i)
	} else {
		PrintInt(-1)
	}

	while (let i = i + 1) < 3 {
		PrintInt(i)
	}
}
