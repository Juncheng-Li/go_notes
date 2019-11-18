# go_notes

learning notes for programming language go

* Lower-case file name
* Space replaced by _
* package name should be lowercase
* 大写字母开头 public（包外部可见）， 小写字母开头 private （包内部可见）
* 导入包没有使用会报错：imported and not used: os
* 注释使用 // 或者 /* */

**Alias import**

```go
package main

import fm "fmt" // alias3

func main() {
   fm.Println("hello, world")
}
```

**Define variable**
* 声明的变量必须要被使用
```go
var a, b, c int
```
```go
var a int
var b bool
var str string
```
```go
var (
	a int
	b bool
	str string
)
```
```go
//根据变量值自动推断类型
var a = 15
var b = false
var str = "Go says hello to the world!"

//声明function中的局部变量
//首选的快捷方式
a := 1
```
```go
//骚操作
a, b, c := 5, 7, "abc"
a, b, c, = 6, 8, "cba"
```


**Print**
```go
//函数 fmt.Sprintf 与 Printf 的作用是完全相同的，不过前者将格式化后的字符串以返回值的形式返回给调用者
fmt.Print("Hello:", 23)
```

**init function**
```go
//会在main之前运行
package trans

import "math"

var Pi float64

func init() {
   Pi = 4 * math.Atan(1) // init() function computes Pi
}
```

**常用类型**
* %d 整数(计算最快) int8, int16, int32, int64，uint8, uint16, uint32, uint64, 默认的int是int64
* %g , %f, %e 浮点型 float32（小数点后7位）, float64（小数点后面15位， 尽量用这个因为math中的运算都会要求接收这个类型）
* 你可以使用 a := uint64(0) 来同时完成类型转换和赋值操作


注意点:
* 不同的类型不能混合使用，即便是int16和int32
* 在格式化字符串里，%d 用于格式化整数（%x 和 %X 用于格式化 16 进制表示的数字），%g 用于格式化浮点型（%f 输出浮点数，%e 输出科学计数表示法），%0nd 用于规定输出长度为n的整数，其中开头的数字 0 是必须的。%n.mg 用于表示数字 n 并精确到小数点后 m 位，除了使用 g 之外，还可以使用 e 或者 f，例如：使用格式化字符串 %5.2e 来输出 3.4 的结果为 3.40e+00
