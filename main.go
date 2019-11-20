package main

import (
	"fmt"
	"time"
)

var week time.Duration

func main() {
	aki := "where"
	fmt.Println(aki)
	fmt.Println(&aki)
	var akiPtr *string
	akiPtr = &aki
	fmt.Println(akiPtr)
}
