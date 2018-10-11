package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to server socket
	conn, _ := net.Dial("tcp", ":8080")

	for {
		// Grab input for message
		input := bufio.NewReader(os.Stdin)
		fmt.Print("Text Message: ")
		msg, _ := input.ReadString('\n')
		// Send to socket
		fmt.Fprintf(conn, msg+"\n")
		// Listen for reply
		msg, _ = bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Server: " + msg)
	}
}
