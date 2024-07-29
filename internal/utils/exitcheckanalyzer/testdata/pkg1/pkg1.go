package main

import (
	"fmt"
	"os"
)

func calc(a int) int {
	b := (a + 8) * 2
	return b
}

func main() {
	fmt.Println(calc(3), calc(5))
	fmt.Println("Hello, world!")
	os.Exit(0) // want "can't use osExit with exit code"
}
