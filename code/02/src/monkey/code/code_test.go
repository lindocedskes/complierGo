package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode //操作码 push
		operands []int  //操作数 AX
		expected []byte //期望make后的值
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}}, //没有操作数
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...) //返回3个字节，第一个操作码，2和3为数字的大端存储

		if len(instruction) != len(tt.expected) { //几个字节-byte
			t.Errorf("instruction has wrong length. want=%d,got=%d",
				len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

// 打印信息字节码转易读的内容（反汇编）
func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2), //[0x00, 0x02]
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`
	concatted := Instructions{}
	for _, ins := range instructions { //字节流切片的切片->字节流切片
		//...运算符可用于将切片或数组“展开”到可变参数列表中
		//一个byte切片展开，就是很多个byte类型的对象
		concatted = append(concatted, ins...) //通过使用...展开ins切片 例如ins：展开后：0x00, 0x01
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, concatted.String())
	}
}

// 测试操作数的反汇编，将Make的byte结果反向解回原来值，易读
func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int //期望的操作数 切片
		bytesRead int   //期望读到的字节长度
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...) //Make

		def, err := Lookup(byte(tt.op)) //通过字节码查找对应操作码，返回该操作码struct类
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}
		//operands读取值int，n读取了几个字节
		operandsRead, n := ReadOperands(def, instruction[1:]) //读操作数，-操作数的反汇编
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
