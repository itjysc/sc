package main


/*
可参考：https://godoc.org/net/http

*/

import (
	"log"
	"net/http"
	"io"
)

func Hello(w http.ResponseWriter, r *http.Request) {     //定义一个名称为Hello的处理函数。
	io.WriteString(w,"尹正杰来也！")
}






func main() {
	http.HandleFunc("/hello",Hello) //给这个http服务绑定一个处理函数，用关键字"hello"来绑定Hello方法。
	log.Fatal(http.ListenAndServe(":8080",nil)) //启动一个http服务使其监听本地所有的8080端口。
}
