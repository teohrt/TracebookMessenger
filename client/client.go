package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var exit = make(chan bool)

func main() {
	var address string
	var name string
	//address:= "10.26.181.239:8080"

	fmt.Println("--------------------------------------------")
	fmt.Println("|      Welcome to TracebookMessenger!      |")
	fmt.Println("--------------------------------------------")
	fmt.Print("Enter server's address: ")
	fmt.Scanln(&address)
	fmt.Print("What is your name? : ")
	fmt.Scanln(&name)
	fmt.Println("--------------------------------------------")
	fmt.Println("|  Type a message and press ENTER to chat. |")
	fmt.Println("--------------------------------------------")

	// Connect to server socket and send name
	conn, _ := net.Dial("tcp", address)
	fmt.Fprintf(conn, name+"\n")

	go sendMessage(conn, name)
	go recieveMessage(conn)

	// Blocking operation
	// Allows the go routines to excecute indefinitely
	<-exit
}

func sendMessage(c net.Conn, n string) {
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
		fmt.Print(msg)
	}
}
