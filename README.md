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
//实际上，fmt 包（第 4.4.3 节）最简单的打印函数也有 2 个返回值：
count, err := fmt.Println(x) // number of bytes printed, nil or 0, error
```

**init function**
```go
//会在main之前运行
//不能有任何的参数和返回值
//强烈建议一个package里面只有一个init
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



**给type重新命名**
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
//绝对值
Abs(x int)


```
运算细节：
* / 对于整数运算而言，结果依旧为整数，例如：9 / 4 -> 2。
* 取余运算符只能作用于整数：9 % 4 -> 1。
* 整数除以 0 可能导致程序崩溃，将会导致运行时的恐慌状态（如果除以 0 的行为在编译时就能被捕捉到，则会引发编译错误）；第 13 章将会详细讲解如何正确地处理此类情况。
* 浮点数除以 0.0 会返回一个无穷尽的结果，使用 +Inf 表示。


**if-else**
* 注意事项 不要同时在 if-else 结构的两个分支里都使用 return 语句，这将导致编译报错 function ends without a return statement（你可以认为这是一个编译器的 Bug 或者特性）。（ 译者注：该问题已经在 Go 1.1 中被修复或者说改进 ）
* 判断string是否为空的2种方法：
    if str == "" { ... }
    if len(str) == 0 {...}
```go
if condition {
	// do something	
} else {
	// do something	
}
```
```go
if condition1 {
	// do something	
} else if condition2 {
	// do something else	
} else {
	// catch-all or default
}
```
```go
if initialization; condition { //感觉可以用作try catch //也可以用作等待读取完成
	// do something
}

//e.g. 类似try catch
if err := file.Chmod(0664); err != nil {
	fmt.Println(err)
	return err
}

//e.g. 等待读取完成
if value, ok := readData(); ok {
…
}

```

**switch**
* 不需要break语句，自带隐藏break
```go
switch num1 {
	case 98, 99:
		fmt.Println("It's equal to 98")
	case 100: 
		fmt.Println("It's equal to 100")
	default:
		fmt.Println("It's not equal to 98 or 100")
	}
```
```go
//fallthrough可以去掉自隐藏的break
switch i {
	case 0: fallthrough
	case 1:
		f() // 当 i == 0 时函数也会被调用
}
```
```go
//带初始化
switch result := calculate() {
	case result < 0:
		...
	case result > 0:
		...
	default:
		// 0
}
```

**for**
```go
for i := 0; i < 5; i++ {
		fmt.Printf("This is the %d iteration\n", i)
}
```
```go
//在循环中同时使用多个计数器
for i, j := 0, N; i < j; i, j = i+1, j-1 {}

for i, j, s := 0, 5, "a"; i < 3 && j < 100 && s != "aaaaa"; i, j,
		s = i+1, j+1, s + "a" {
		fmt.Println("Value of i, j, s:", i, j, s)
	}
```
```go
//无限循环
for {}
for ;; { }
for i:=0;;i++ {}
```
```go
//亚洲字体在print中的使用
str := "Go is a beautiful language!"
fmt.Printf("The length of str is: %d\n", len(str))
for pos, char := range str {
	fmt.Printf("Character on position %d is: %c \n", pos, char)
}
fmt.Println()
str2 := "Chinese: 日本語"
fmt.Printf("The length of str2 is: %d\n", len(str2))
for pos, char := range str2 {
	fmt.Printf("character %c starts at byte position %d\n", char, pos)
}
fmt.Println()
fmt.Println("index int(rune) rune    char bytes")
for index, rune1 := range str2 {
	fmt.Printf("%-2d      %d      %U '%c' % X\n", index, rune1, rune1, rune1, []byte(string(rune1)))
}
```
for-range
```go

```

**错误**
判断有没有错误
```go
if err != nil {}
```

**array**
* slice 的cap是slice的开始到母arr的结尾的长度
* slice 的len就是slice本身的长度

**defer**
* defer的东西会在函数推出前执行，相当于finally
```go
func ReadWrite() bool {
	file.Open("file")
	defer file.Close()
	if failureX {
		return false
	}
	if failureY {
		return false
	}
	return true
}
```

**struct**
```go
type person struct {
	name string
	age int
}

//声明
var P person

//使用
P.name = "Astaxie"  // 赋值"Astaxie"给P的name属性.
P.age = 25  // 赋值"25"给变量P的age属性
fmt.Printf("The person's name is %s", P.name)  // 访问P的name属性.

// 赋值初始化
tom.name, tom.age = "Tom", 18

// 两个字段都写清楚的初始化
bob := person{age:25, name:"Bob"}

// 按照struct定义顺序初始化值
paul := person{"Paul", 43}
```
* 匿名字段 字段继承
* 字段继承如果有相同名字的，优先访问外层的
```go
package main

import "fmt"

type Human struct {
	name string
	age int
	weight int
}

type Student struct {
	Human  // 匿名字段，那么默认Student就包含了Human的所有字段
	speciality string
}

func main() {
	// 我们初始化一个学生
	mark := Student{Human{"Mark", 25, 120}, "Computer Science"}

	// 我们访问相应的字段
	fmt.Println("His name is ", mark.name)
	fmt.Println("His age is ", mark.age)
	fmt.Println("His weight is ", mark.weight)
	fmt.Println("His speciality is ", mark.speciality)
	// 修改对应的备注信息
	mark.speciality = "AI"
	fmt.Println("Mark changed his speciality")
	fmt.Println("His speciality is ", mark.speciality)
	// 修改他的年龄信息
	fmt.Println("Mark become old")
	mark.age = 46
	fmt.Println("His age is", mark.age)
	// 修改他的体重信息
	fmt.Println("Mark is not an athlet anymore")
	mark.weight += 60
	fmt.Println("His weight is", mark.weight)
}
```

**method**
* 有点像java的overloading (同一个函数名，不同的函数签名)
```go
package main

import (
	"fmt"
	"math"
)

type Rectangle struct {
	width, height float64
}

type Circle struct {
	radius float64
}

func (r Rectangle) area() float64 {
	return r.width*r.height
}

func (c Circle) area() float64 {
	return c.radius * c.radius * math.Pi
}


func main() {
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	c1 := Circle{10}
	c2 := Circle{25}

	fmt.Println("Area of r1 is: ", r1.area())
	fmt.Println("Area of r2 is: ", r2.area())
	fmt.Println("Area of c1 is: ", c1.area())
	fmt.Println("Area of c2 is: ", c2.area())
}
```

**自定义类型**
* 也可以理解为对类型的重命名
```go
//这里int就变成了age
type ages int

type money float32

type months map[string]int

m := months {
	"January":31,
	"February":28,
	...
	"December":31,
}
```

**channels**
```go
ci := make(chan int)
cs := make(chan string)
cf := make(chan interface{})

ch <- v    // 发送v到channel ch.
v := <-ch  // 从ch中接收数据，并赋值给v
```

**多协程**
* 必须设置GOMAXPROCS大于1来激活多个线程
* 经验指示协程的数量 n - 1 能获得最佳性能（n为cpu核心的数量）
