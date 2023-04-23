package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	//env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		//删除之前的求值和环境
		//evaluated := evaluator.Eval(program, env)
		//if evaluated != nil {
		//	io.WriteString(out, evaluated.Inspect())
		//	io.WriteString(out, "\n")
		//}

		//调用编译器
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! 编译错误:\n %s\n", err)
			continue
		}

		//调用虚拟机,输出栈顶元素
		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! 虚拟机执行字节码错误:\n %s\n", err)
			continue
		}

		//输出栈顶元素，输出
		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")

	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-""lyt""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
