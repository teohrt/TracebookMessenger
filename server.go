package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	fmt.Println("Welcome to TracebookMessenger!")

	ln, _ := net.Listen("tcp", ":8080")
	// Accept connection
	conn, _ := ln.Accept()

	for {

		// Listens for message that ends in a newline character
		msg, _ := bufio.NewReader(conn).ReadString('\n')

		// Log incomming message
		fmt.Print(string(msg))

		// Echo back to the client
		conn.Write([]byte("Echo: " + msg))
	}
}
