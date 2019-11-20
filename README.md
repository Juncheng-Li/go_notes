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
* %s ???
* %t boolean
* String
```go
s := "hel" + "lo,"
s += "world!"
fmt.Println(s) //输出 “hello, world!”

//查看s前缀是否是string
strings.HasPrefix(s, string) //bool
//查看s后缀是否是string
strings.HasSuffix(s, string) //bool
//包含
strings.Contains(s, string) //bool
//index
strings.Index(s, str string) //int
strings.LastIndex(s, str string) //int
strings.IndexRune(s string, r rune) //查询非ASCII编码的字符在字符串中的位置//int
//replace
strings.Replace(str, old, new, n) //n为替换前n个，n=-1就替换all
//字符串出现次数, str在s中出现的次数
strings.Count(s, str) //int
//重复字符串, 就是重复打印很多遍
strings.Repeat(s, count)
//大写小写转换
strings.ToLower(s)
strings.ToUpper(s)
//去除空格(首尾的空格)
strings.TrimSpace(s)
//去除首尾的"cut"
strings.Trim(s, "cut")
//分割字符串 (根据空格分，output是一个array)
strings.Fields(s)
//分割字符串（根据给的符号分，output是一个array）
strings.Split(s, ",")
//拼接字符串
strings.Join(sl, sep) //sl是一个array，sep是符号,sl是直接使用其他function的输出的注意，有可能function的输出并不支持所以造成的结果并没有加入sep合并

//字符转换
strconv.IntSize //告诉你系统int的size
strconv.Itoa(i) //i是int
val, err = strconv.Atoi(s) //s是string
strconv.ParseFloat(s string, bitSize int) (f float64, err error)
```

**时间**
```go
t := time.Now()
fmt.Println(t) // e.g. Wed Dec 21 09:52:14 +0100 RST 2011
fmt.Printf("%02d.%02d.%4d\n", t.Day(), t.Month(), t.Year())
// 21.12.2011
t = time.Now().UTC()
fmt.Println(t) // Wed Dec 21 08:52:14 +0000 UTC 2011
fmt.Println(time.Now()) // Wed Dec 21 09:52:14 +0100 RST 2011
// calculating times:
week = 60 * 60 * 24 * 7 * 1e9 // must be in nanosec
week_from_now := t.Add(time.Duration(week))
fmt.Println(week_from_now) // Wed Dec 28 08:52:14 +0000 UTC 2011
// formatting times:
fmt.Println(t.Format(time.RFC822)) // 21 Dec 11 0852 UTC
fmt.Println(t.Format(time.ANSIC)) // Wed Dec 21 08:56:34 2011
fmt.Println(t.Format("02 Jan 2006 15:04")) // 21 Dec 2011 08:52
s := t.Format("20060102")
fmt.Println(t, "=>", s) // Wed Dec 21 08:52:14 +0000 UTC 2011 => 20111221
```

**指针**
* 目的是为了廉价传递一个变量的引用
* 不能得到文字或者常量的地址（const）
* 指针算法是不合法的，因此 c = *p++ 在 Go 语言的代码中是不合法的，又比如pointer + 2在GO中也是不合法的
* 空指针的反向引用会报错
* 指针可以迭代（嵌套）
```go
//通过加 & 获得地址
fmt.Println(aki)
fmt.Println(&aki)
//通过指针前面加 * 获得指针指向地址的内容（叫做 反引用/内容引用/间接引用）
var akiPtr *string = &aki
fmt.Println(*akiPtr) //返回aki的内容
//定义指针
var intP *int
```



给type重新命名
```go
package main
import "fmt"

//把int重新命名为TZ
type TZ int

func main() {
	var a, b TZ = 3, 4
	c := a + b
	fmt.Printf("c has the value: %d", c) // 输出：c has the value: 7
}
```

注意点:
* 不同的类型不能混合使用，即便是int16和int32
* 你可以使用 a := uint64(0) 来同时完成类型转换和赋值操作
* 在格式化字符串里，%d 用于格式化整数（%x 和 %X 用于格式化 16 进制表示的数字），%g 用于格式化浮点型（%f 输出浮点数，%e 输出科学计数表示法），%0nd 用于规定输出长度为n的整数，其中开头的数字 0 是必须的。%n.mg 用于表示数字 n 并精确到小数点后 m 位，除了使用 g 之外，还可以使用 e 或者 f，例如：使用格式化字符串 %5.2e 来输出 3.4 的结果为 3.40e+00

**运算**
```go
i++
i += 1
```

运算细节：
* / 对于整数运算而言，结果依旧为整数，例如：9 / 4 -> 2。
* 取余运算符只能作用于整数：9 % 4 -> 1。
* 整数除以 0 可能导致程序崩溃，将会导致运行时的恐慌状态（如果除以 0 的行为在编译时就能被捕捉到，则会引发编译错误）；第 13 章将会详细讲解如何正确地处理此类情况。
* 浮点数除以 0.0 会返回一个无穷尽的结果，使用 +Inf 表示。

