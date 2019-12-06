package main

import (
	"fmt"
	"os"
)

var (
	h, H bool
	v    bool
	q    *bool
	D    string
	Conf string
)

func init() {

}
func main() {
	fmt.Println(os.Getgid())
	fmt.Println(os.Getuid())
	fmt.Println(os.Getgroups())
}
