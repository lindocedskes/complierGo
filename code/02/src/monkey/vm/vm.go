package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object   //常量池
	instructions code.Instructions //编译后的字节码

	stack []object.Object
	sp    int //始终指向栈中的空闲槽，栈顶为sp-1
}

// 创建虚拟机
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Insttuctions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// 运行虚拟机 ，通过ip（指令计数器）指针遍历，进行取值
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip]) //取指令，指向每条指令的操作码 1字节

		switch op {
		case code.OpConstant: //解码-指令解释执行
			constIndex := code.ReadUint16(vm.instructions[ip+1:]) //读取一个操作数占2字节-16bit的,常量池索引
			ip += 2                                               //下一条指令的操作码
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd: //解码-指令解释执行
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			//翻译成汇编命令
			//fmt.Printf("ADD %d ,%d \n", leftValue, rightValue)

			result := leftValue + rightValue
			vm.push(&object.Integer{Value: result})
		}

	}
	return nil
}

// 查看栈顶元素
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// 压栈
func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

// 出栈
func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}
