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

	// Display IP and port of the server
	self, _ := net.Dial("udp", "8.8.8.8:80")
	localAddr := self.LocalAddr().(*net.UDPAddr)
	fmt.Println("Server hosted at: " + localAddr.IP.String() + port)

	fmt.Println("--------------------------------------------")
	fmt.Println("|      Welcome to TracebookMessenger!      |")
	fmt.Println("|        Listening for connections.        |")
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
		msg = "( " + c.name + " ) : " + msg

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
