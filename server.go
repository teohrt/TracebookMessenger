package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn     net.Conn
	loggedIn bool
	name     string
}

var clients = []Client{}

func main() {

	port := ":8080"

	fmt.Println("--------------------------------------------")
	fmt.Println("|      Welcome to TracebookMessenger!      |")
	fmt.Println("|  Listening for connections on port " + port + " |")
	fmt.Println("--------------------------------------------")

	ln, _ := net.Listen("tcp", port)

	for {
		// Accept connection and grab client name
		conn, _ := ln.Accept()
		name, _ := bufio.NewReader(conn).ReadString('\n')
		name = name[:len(name)-1] // strips the newline character from input

		// Add to list of clients
		c := Client{conn, true, name}
		clients = append(clients, c)

		fmt.Println("***New Client Connected : " + c.name)
		go c.start()
	}
}

// Handles client traffic
func (c Client) start() {
	for {
		// Listens for message that ends in a newline character
		msg, _ := bufio.NewReader(c.conn).ReadString('\n')

		// Log incomming message
		fmt.Print(string("( " + c.name + " ) : " + msg))

		// Echo back to all clients except the sender
		for _, client := range clients {
			if client.loggedIn && client != c {
				client.conn.Write([]byte(msg))
			}
		}
	}
}
