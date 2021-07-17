package main

import (
	"errors"
	"fmt"
)

func main() {
	txt := "hi"
	msg(txt)
}

func msg(s string) {
	defer fmt.Println("before panic")
	if s == "hi"{
		panic(errors.New("panic here"))
	}
	defer fmt.Println("after panic")
}