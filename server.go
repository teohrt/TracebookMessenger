package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn     net.Conn
	loggedIn bool
}

var clients = []Client{}

func main() {

	fmt.Println("Welcome to TracebookMessenger!")

	ln, _ := net.Listen("tcp", ":8080")

	for {
		// Accept connection
		conn, _ := ln.Accept()
		fmt.Println("New Client Connected.")

		// Add to list of clients
		c := Client{conn, true}
		clients = append(clients, c)

		go c.start()
	}
}

// Handles client traffic
func (c Client) start() {
	for {
		// Listens for message that ends in a newline character
		msg, _ := bufio.NewReader(c.conn).ReadString('\n')

		// Log incomming message
		fmt.Print(string(msg))

		// Echo back to all clients except the sender
		for _, client := range clients {
			if client.loggedIn && client != c {
				client.conn.Write([]byte(msg))
			}
		}
	}
}
