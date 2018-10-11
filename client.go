package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var exit = make(chan bool)

func main() {
	fmt.Println("Type a message and press ENTER to chat.")

	// Connect to server socket
	conn, _ := net.Dial("tcp", ":8080")

	go sendMessage(conn)
	go recieveMessage(conn)

	// Blocking operation
	// Allows the go routines to excecute indefinitely
	<-exit
}

func sendMessage(c net.Conn) {
	for {
		// Grab user input for message
		input := bufio.NewReader(os.Stdin)
		msg, _ := input.ReadString('\n')
		// Send to server socket
		fmt.Fprintf(c, msg+"\n")
	}
}

func recieveMessage(c net.Conn) {
	for {
		// Listen for reply
		msg, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Println("Server: " + msg)
	}
}
