package main

import (
	"fmt"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", ":8080")
	for {
		// Grab connection
		conn, _ := ln.Accept()
		var msg []byte

		// Reads from connection and stores in msg
		fmt.Fscan(conn, &msg)
		fmt.Println("Message:", string(msg))
	}
}