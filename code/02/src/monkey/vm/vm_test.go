package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input) //词法+语法

		comp := compiler.New()
		err := comp.Compile(program) //编译
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode()) //虚拟机创建
		err = vm.Run()             //运行
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop() //弹栈

		testExpectedObject(t, tt.expected, stackElem)
	}
}

// 词法+语法分析，返回AST
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	return p.ParseProgram()
}

// 验证是期望的对象
func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual object.Object,
) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value.got=%d, want=%d", result.Value, expected)
	}
	return nil
}
