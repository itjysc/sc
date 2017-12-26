package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "hello world ，你好，世界")
	} else {
		fmt.Fprintf(w, "method not allowed")
	}
}

func main() {
	addr := ":8080"
	http.HandleFunc("/helloworld", handler)
	http.ListenAndServe(addr, nil)
}