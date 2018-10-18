package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn     net.Conn
	loggedIn bool
	name     string
}

var clients = []Client{}

func main() {

	port := ":8080"

	// Display IP and port of the server
	self, _ := net.Dial("udp", "8.8.8.8:80")
	localAddr := self.LocalAddr().(*net.UDPAddr)
	fmt.Println("Server hosted at: " + localAddr.IP.String() + port)

	fmt.Println("--------------------------------------------")
	fmt.Println("|      Welcome to TracebookMessenger!      |")
	fmt.Println("|        Listening for connections.        |")
	fmt.Println("--------------------------------------------")

	ln, _ := net.Listen("tcp", port)

	go sendMessage()

	for {
		// Accept connection and grab client name
		conn, _ := ln.Accept()
		name, _ := bufio.NewReader(conn).ReadString('\n')
		name = name[:len(name)-1] // strips the newline character from input

		// Add to list of clients
		c := Client{conn, true, name}
		clients = append(clients, c)

		logAndSend("***New Client Connected : "+c.name+"\n", c)

		go c.start()
	}
}

// Handles client traffic
func (c Client) start() {
	for {
		// Listens for message that ends in a newline character
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			msg = "***" + c.name + " left the chat.\n"
		} else {
			msg = ("( " + c.name + " ) : " + msg)
		}

		logAndSend(msg, c)

		if err != nil {
			break
		}
	}
}

// Log and echo back to all clients except the sender specified
func logAndSend(msg string, c Client) {
	fmt.Print(string(msg))
	for _, client := range clients {
		if client.loggedIn && client != c {
			client.conn.Write([]byte(msg))
		}
	}
}

func sendMessage() {
	for {
		// Grab user input for message
		input := bufio.NewReader(os.Stdin)
		msg, _ := input.ReadString('\n')

		fmt.Print(string(msg))
		for _, client := range clients {
			client.conn.Write([]byte("( SERVER ) : " + msg))
		}
	}
}
