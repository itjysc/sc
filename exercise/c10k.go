package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func handler(conn net.Conn) {
	defer conn.Close()
	fmt.Fprintf(conn, "%s", time.Now().String())
}

func main() {
	addr := ":8080"
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handler(conn)

	}

}