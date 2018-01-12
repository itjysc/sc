package main

import (
	"net/http"
	"log"
)

func main() {
	http.Handle("/",http.FileServer(http.Dir("."))) /*注意后面加了斜线表示其是一个目录。
	FileServer(http.Dir(".")）表示将以"/"为根目录，我们只需要把静态页面放这里就可以啦！*/
	log.Fatal(http.ListenAndServe(":8888",nil))

}
