package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type Compiler struct { //编译器
	instructions code.Instructions //切片类型
	constants    []object.Object   //常量池
}
type Bytecode struct { //?传给虚拟机，编译器内做断言
	Insttuctions code.Instructions //切片类型
	Constants    []object.Object   //常量池
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// 需要传给虚拟机的内容：编译后的字节流切片+求值后的常量池
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Insttuctions: c.instructions, //存放编译后的字节流切片
		Constants:    c.constants,    //存放求值后的常量池
	}
}

// 编译器遍历AST，递归子程序法
func (c *Compiler) Compiler(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program: //程序开始节点,向下递归
		for _, s := range node.Statements {
			err := c.Compiler(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement: //表达式语句
		err := c.Compiler(node.Expression) //传入表达式
		if err != nil {
			return err
		}
	case *ast.InfixExpression: //中缀
		err := c.Compiler(node.Left) //解析左节点
		if err != nil {
			return err
		}
		err = c.Compiler(node.Right) //解析右节点
		if err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		//对Inter字面量求值=>返回Integer对象
		integer := &object.Integer{Value: node.Value}
		// '发出'-生成和输出OpConstant指令，加入常量池，Integer对象返回索引
		c.emit(code.OpConstant, c.addConstant(integer))
	}

	return nil
}

// 将求解常量(一个实例对象)加入常量池，返回索引
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1 //返回在constants的索引
}

// '发出'-生成和输出指令，编译指令
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...) //编译指令
	pos := c.addInstruction(ins)      //已经有的字节流切片添加编译的切片，返回最大位置
	return pos
}

// 已经有的字节流切片添加编译的切片，返回最大位置
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
