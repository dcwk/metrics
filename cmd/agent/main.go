package main

import (
	"fmt"
	"time"
)

func main() {
	for true {
		time.Sleep(2 * time.Second)
		fmt.Println("Test time")
	}
}
