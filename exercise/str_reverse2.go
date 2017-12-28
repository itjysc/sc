package main

import (
	"fmt"
)

func main() {
	strs := "北京欢迎你"
	num := []rune(strs)
	lang := len(num)
	for i, j := 0, lang-1; i < j; i, j = i+1, j-1 { //这种思路就是把最后一个字符和第一个字符互换，循环到最中间的那个就不做任何操作
		num[i], num[j] = num[j], num[i]
	}
	fmt.Printf("原始的字符串是：%s\n", strs)
	fmt.Printf("反转后的字符串是：%s\n", string(num))
}
