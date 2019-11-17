# go_notes

learning notes for programming language go

* Lower-case file name
* Space replaced by _
* package name should be lowercase
* 大写字母开头 public（包外部可见）， 小写字母开头 private （包内部可见）
* 导入包没有使用会报错：imported and not used: os
* 注释使用 // 或者 /* */

Alias import

```go
package main

import fm "fmt" // alias3

func main() {
   fm.Println("hello, world")
}
```

Define variable

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