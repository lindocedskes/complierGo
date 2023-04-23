package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}       //期望常量池
	expectedInsttuctions []code.Instructions //期望字节切片集合
}

// 测试整数计算
func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{ //公共测试集合
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInsttuctions: []code.Instructions{ //字节流切片的切片
				code.Make(code.OpConstant, 0), //OpConstant int整数 ->
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd), //将栈中2个add
			},
		},
	}

	runCompilerTest(t, tests) //测试辅助函数-编译测试
}

// 测试辅助函数-编译测试
func runCompilerTest(t *testing.T, tests []compilerTestCase) {
	//Helper将调用函数标记为测试辅助函数。在打印文件和行信息时，该函数将被跳过。Helper可以同时从多个goroutine中调用
	t.Helper() //?

	for _, tt := range tests {
		program := parse(tt.input) //词法分析、语法分析返回AST

		compiler := New()
		err := compiler.Compile(program) //编译器编译
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		//Bytecode封装编译后的字节流切片+求值后的常量池，传给虚拟机
		bytecode := compiler.Bytecode()
		//判断字节流切片是否正确
		err = testInstructions(tt.expectedInsttuctions, bytecode.Insttuctions)
		if err != nil {
			t.Fatalf("testInstructions error: %s", err)
		}
		//判断常量池是否正确
		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants error: %s", err)
		}

	}
}

// 词法分析、语法分析返回AST
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// 判断字节流切片是否正确
func testInstructions(
	expected []code.Instructions, //字节流切片的切片
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected) //[]code.Instructions->code.Instructions

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instruction length.\nwant=%q\ngot =%q", concatted, actual)
	}

	for i, ins := range concatted { //按字节比较
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction length.\nwant=%q\ngot =%q", concatted, actual)
		}
	}
	return nil
}

// testInstructions的辅助函数——输入expected是字节切片的切片,转为一个字节切片
func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{} //空的字节切片

	for _, ins := range s { //顺序存
		out = append(out, ins...)
	}
	return out
}

// 判断常量池 expected：[]interface{}{1, 2}
func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constant.got=%d,want =%d", len(actual), len(expected))
	}

	for i, constant := range expected { //遍历期望的常量
		switch constant := constant.(type) { //常量池中查找
		case int: //该常量是int
			err := testIntegerObject(int64(constant), actual[i]) //判断是否是int对象
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}
	return nil
}

// 判断是否是int对象
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}
