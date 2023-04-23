## 简易解释器 Go -参考书：Writing A Compiler In Go (Thorsten Ball)

learned：
+ 2.1第一条指令
  + 2.1.4 ：第一个编译器测试通过：定义了一个操作码：OpConstant 整数； code.go + compiler.go
  一个通过遍历AST发出OpConstant的编译器，并计算常量整数字面量已经添加到常量池。编译器接口能传递编译的结果给虚拟机。
  + 2.1.5 :简易虚拟机，解释执行字节码，OpConstant指令可被解释执行
  + 2.2 栈上加法：OpAdd指令实现，2数相加，栈实现。
  + 2.3 repl实现：删除之前的求值和环境，增加调用编译器+VM
  

