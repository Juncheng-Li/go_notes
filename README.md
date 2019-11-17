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


Print
```go
//函数 fmt.Sprintf 与 Printf 的作用是完全相同的，不过前者将格式化后的字符串以返回值的形式返回给调用者
fmt.Print("Hello:", 23)
```