package sampleprograms

const EqualComparison = `proc main() () {
	mut a int = 3
	let b int = 3
	if a == b {
		print("true")
	} else {
		print("false")
	}

	while a == b {
		print("%d\n", a)
		a = a + 1
	}
}`

const NotEqualComparison = `proc main() () {
	mut a int = 3
	let b int = 3
	if a != b {
		print("true")
	} else {
		print("false")
	}

	while a != b {
		print("%d\n", a)
		a = a + 1
	}
}`

const GreaterComparison = `proc main() () {
	mut a int = 4
	let b int = 3
	if a > b {
		print("true")
	} else {
		print("false")
	}

	while a > b {
		print("%d\n", a)
		a = a - 1
	}
}`

const GreaterOrEqualComparison = `proc main() () {
	mut a int = 4
	let b int = 3
	if a >= b {
		print("true")
	} else {
		print("false")
	}

	while a >= b {
		print("%d\n", a)
		a = a - 1
	}
}`

const LessThanComparison = `proc main() () {
	mut a int = 4
	let b int = 3
	if a < b {
		print("true")
	} else {
		print("false")
	}

	while a < b {
		print("%d\n", a)
		a = a + 1
	}
}`

const LessThanOrEqualComparison = `proc main() () {
	mut a int = 1
	let b int = 3
	if a <= b {
		print("true")
	} else {
		print("false")
	}

	while a <= b {
		print("%d\n", a)
		a = a + 1
	}
}`
