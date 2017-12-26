package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func handle(conn net.Conn) {
	fmt.Fprint(conn, "Welcome to Yin Zhengjie's home page", time.Now().String())
	conn.Close()
}

func main() {
	l, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}