package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const ( //数字作为常量池的索引
	//iota将定义常量（操作码）按递增赋予int初值
	OpConstant Opcode = iota //指令——生成常量池索引：OpConstant int数字
	OpAdd
)

type Instructions []byte //编译后存字节的集合
type Opcode byte         //操作码uint8

type Definition struct {
	Name          string //操作码的名字 mov(操作码) AX,BX（操作数）
	OperandWidths []int  //每位操作数占的字节。数组下标：第几个，值：占字节数
}

// 操作码byte->操作码类型信息 映射
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, //操作码的名字和每个操作数占用2个字节
	OpAdd:      {"OpAdd", []int{}},
}

// 通过字节码op查找对应操作码，返回该操作码struct类
// 例如op=byte(OpConstant) ，OpConstant=0，字节码为0
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// 返回操作码00000000 操作数... 编译后的字节流，第一个字节必对应操作码，后面为操作数，operands是操作数对象的在常量池的索引int
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op] //返回byte对应的操作码类型信息，字节->指令
	if !ok {
		return []byte{}
	}

	instructionLen := 1                   //统计总字节数，初始操作码占1B
	for _, w := range def.OperandWidths { //range返回数组下标和值
		instructionLen += w //累加每位操作数占的字节
	}

	instruction := make([]byte, instructionLen) //开辟对应长度空间
	instruction[0] = byte(op)                   //操作码转B 存入

	offset := 1
	for i, o := range operands { //每位操作数，转B存入。operands操作数
		width := def.OperandWidths[i] //每个操作数占几字节
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o)) //大端编码，存入[]byte
		}
		offset += width
	}
	return instruction
}

// 反汇编-操作数byte->int
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths)) //存放操作数转换为int后的切片
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2: //操作数占2个字节
			operands[i] = int(ReadUint16(ins[offset:])) //按2B转对应int
		}
		offset += width //读下一个操作数
	}
	return operands, offset
}

// 公共函数-ReadUint16
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins) //将切片的前 2 个字节转换为 uint16 值，字节是大端存储
}

// 打印字符串 返回易读的字节流
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i]) // 通过字节码op查找对应操作码，返回该操作码struct类
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		//operands转换int的操作数切片，read：读了几位
		operands, read := ReadOperands(def, ins[i+1:]) //反汇编-操作数byte->[]int
		//输出第i个指令，
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

// 格式化输出 对应的操作码名称和操作数，并验证操作数对应个数正确
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths) //操作数个数

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)

	}

	switch operandCount { //操作数个数
	case 0:
		return def.Name
	case 1: //格式化输出一个参数的指令名+操作数int
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
