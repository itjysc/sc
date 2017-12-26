package main

import "fmt"

func main() {
	var x int
	x = 1
	var y int
	y = 2

	fmt.Println(x, y)
	//x和y的初始化指针
	fmt.Println(&x, &y)
	swap(&x, &y)
	fmt.Println(x, y)
	//x和y的值对调后
	fmt.Println(&x, &y)

	swap_a(&x, &y)
	//x和y的指针对调后
	fmt.Println(&x, &y)
	fmt.Println(x, y)
}

//将x和y的值对调
func swap(p *int, q *int) {
	var t = *p
	*p = *q
	*q = t
}

//x和y的指针对调，值不变
func swap_a(p *int, q *int) {
	var t = p
	p = q
	q = t
}