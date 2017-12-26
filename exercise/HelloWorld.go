package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("Hello World!")
	io.WriteString(os.Stdout, "Hello World!!")
}