package codegen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/driusan/lang/compiler/irgen"

	"github.com/driusan/lang/parser/ast"
	"github.com/driusan/lang/parser/sampleprograms"
)

func RunProgram(name, p string) error {
	d, err := ioutil.TempDir("", "langtest"+name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(d)

	f, err := os.Create(d + "/main.s")
	if err != nil {
		return err
	}
	defer f.Close()

	prog, err := ast.Parse(p)
	if err != nil {
		return err
	}

	for _, v := range prog {
		switch v.(type) {
		case ast.FuncDecl, ast.ProcDecl:
			fnc, err := irgen.GenerateIR(v)
			if err != nil {
				return err
			}
			if err := Compile(f, fnc); err != nil {
				return err
			}
		case ast.TypeDefn:
			// No IR for types, we've already verified them.
		default:
			panic("Unhandled AST node type for code generation")
		}
	}

	// FIXME: Make this more robust, or at least move it to a helper. It
	// will only work on Plan 9 right now.
	cmd := exec.Command("6a", "-o", d+"/main.6", d+"/main.s")
	val, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	cmd = exec.Command("6l", "-o", d+"/main", d+"/main.6")
	val, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Finally, run the compiled binary!
	cmd = exec.Command(d + "/main")
	val, err = cmd.Output()
	fmt.Println(string(val))
	return err
}

func TestCompileHelloWorld(t *testing.T) {
	prgAst, err := ast.Parse(sampleprograms.HelloWorld)
	if err != nil {
		t.Fatal(err)
	}

	var w bytes.Buffer
	fnc, err := irgen.GenerateIR(prgAst[0])
	if err != nil {
		t.Fatal(err)
	}

	if err := Compile(&w, fnc); err != nil {
		t.Fatal(err)
	}

	// HelloWorld is simple enough that we can hand compile it. Most other
	// programs we just compile it, run it and check the output to ensure
	// that they work.
	expected := `#pragma lib "libstdio.a"
#pragma lib "libc.a"

TEXT main(SB), $32
	DATA .string0<>+0(SB)/8, $"Hello, w"
	DATA .string0<>+8(SB)/8, $"orld!\n\z\z"
	GLOBL .string0<>+0(SB), $16
	MOVQ $.string0<>+0(SB), BP
	CALL printf+0(SB)
	RET
`

	if val := w.Bytes(); string(val) != expected {
		t.Errorf("Unexpected ASM output compiling HelloWorld. got %s; want %s", val, expected)
	}
}

func ExampleCompileEmptyMain() {
	if err := RunProgram("emptymain", sampleprograms.EmptyMain); err != nil {
		// fmt.Println(err.Error())
	}
	// Output:
}

func ExampleHelloWorld() {
	if err := RunProgram("helloworld", sampleprograms.HelloWorld); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: Hello, world!
}

func ExampleLetStatement() {
	if err := RunProgram("letstatement", sampleprograms.LetStatement); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 5
}
func ExampleLetStatementShadow() {
	if err := RunProgram("letstatementshadow", sampleprograms.LetStatementShadow); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 5
	// hello
}

func ExampleCompileHelloWorld2() {
	if err := RunProgram("helloworld2", sampleprograms.HelloWorld2); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: Hello, world!
	//  World??
	//  Hello, world!
}

func ExampleCompileTwoProcs() {
	if err := RunProgram("twoprocs", sampleprograms.TwoProcs); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 3
}

func ExampleOutOfOrder() {
	if err := RunProgram("outoforder", sampleprograms.OutOfOrder); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 3
}

func ExampleMutAddition() {
	if err := RunProgram("mutaddition", sampleprograms.MutAddition); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 8
}

func ExampleSimpleFunc() {
	if err := RunProgram("simplefunc", sampleprograms.SimpleFunc); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 3
}

func ExampleSumToTen() {
	if err := RunProgram("sumtoten", sampleprograms.SumToTen); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 55
}

func ExampleSumToTenRecursive() {
	if err := RunProgram("sumtotenrecursive", sampleprograms.SumToTenRecursive); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 55
}

func ExampleFizzBuzz() {
	if err := RunProgram("fizzbuzz", sampleprograms.Fizzbuzz); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: 1
	// 2
	// fizz
	// 4
	// buzz
	// fizz
	// 7
	// 8
	// fizz
	// buzz
	// 11
	// fizz
	// 13
	// 14
	// fizzbuzz
	// 16
	// 17
	// fizz
	// 19
	// buzz
	// fizz
	// 22
	// 23
	// fizz
	// buzz
	// 26
	// fizz
	// 28
	// 29
	// fizzbuzz
	// 31
	// 32
	// fizz
	// 34
	// buzz
	// fizz
	// 37
	// 38
	// fizz
	// buzz
	// 41
	// fizz
	// 43
	// 44
	// fizzbuzz
	// 46
	// 47
	// fizz
	// 49
	// buzz
	// fizz
	// 52
	// 53
	// fizz
	// buzz
	// 56
	// fizz
	// 58
	// 59
	// fizzbuzz
	// 61
	// 62
	// fizz
	// 64
	// buzz
	// fizz
	// 67
	// 68
	// fizz
	// buzz
	// 71
	// fizz
	// 73
	// 74
	// fizzbuzz
	// 76
	// 77
	// fizz
	// 79
	// buzz
	// fizz
	// 82
	// 83
	// fizz
	// buzz
	// 86
	// fizz
	// 88
	// 89
	// fizzbuzz
	// 91
	// 92
	// fizz
	// 94
	// buzz
	// fizz
	// 97
	// 98
	// fizz
}

func ExampleSomeMath() {
	if err := RunProgram("somemath", sampleprograms.SomeMath); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: Add: 3
	// Sub: -1
	// Mul: 6
	// Div: 3
	// Complex: 5
}

func ExampleEqualComparison() {
	if err := RunProgram("equalcompare", sampleprograms.EqualComparison); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: true
	// 3
}

func ExampleNotEqualComparison() {
	if err := RunProgram("notequalcompare", sampleprograms.NotEqualComparison); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: false
}

func ExampleGreaterComparison() {
	if err := RunProgram("greatercompare", sampleprograms.GreaterComparison); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: true
	// 4
}

func ExampleGreaterOrEqualComparison() {
	if err := RunProgram("greaterorequalcompare", sampleprograms.GreaterOrEqualComparison); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: true
	// 4
	// 3
}

func ExampleLessThanEqualComparison() {
	if err := RunProgram("lessthanequalcompare", sampleprograms.LessThanComparison); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: false
}

func ExampleLessThanOrEqualComparison() {
	if err := RunProgram("lessthanorequalcompare", sampleprograms.LessThanOrEqualComparison); err != nil {
		// fmt.Println(err.Error())
	}
	// Output: true
	// 1
	// 2
	// 3
}

func ExampleUserDefinedType() {
	if err := RunProgram("usertype", sampleprograms.UserDefinedType); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 4
}
func ExampleTypeInference() {
	if err := RunProgram("typeinference", sampleprograms.TypeInference); err != nil {
		//fmt.Println(err.Error())
	}
	// Output: 0, 4
}
