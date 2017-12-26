package main

import (
	"fmt"
	"os"
)

var retstr string

func main() {
	str1 := os.Args[1]
	var array = []byte(str1)               // str1转换
	for i := len(array) - 1; i >= 0; i-- { // 翻转
		retstr = retstr + string(array[i])
	}
	fmt.Println(retstr)
}
