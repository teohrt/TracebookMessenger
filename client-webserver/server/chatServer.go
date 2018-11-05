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

func runChatServer(chatPort string) {

	// Listen for incomming tcp requests
	ln, _ := net.Listen("tcp", chatPort)

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

// Handles outgoing server messages to clients
func handleOutgoing() {
	for {
		// Grab user input for message
		input := bufio.NewReader(os.Stdin)
		msg, _ := input.ReadString('\n')
		history = append(history, msg)
		for _, client := range clients {
			client.conn.Write([]byte("( SERVER ) : " + msg))
		}
	}
}

// Log and echo back to all clients except the sender specified
func logAndSend(msg string, c Client) {
	fmt.Print(string(msg))
	history = append(history, msg)
	for _, client := range clients {
		if client.loggedIn && client != c {
			client.conn.Write([]byte(msg))
		}
	}
}

func welcomeMessage(cp string, wp string) {
	// Display IP and port of the chat server
	self, _ := net.Dial("udp", "8.8.8.8:80")
	localAddr := self.LocalAddr().(*net.UDPAddr)

	logger("Web  server hosted at : " + localAddr.IP.String() + wp)
	logger("Chat server hosted at : " + localAddr.IP.String() + cp)
	logger("--------------------------------------------")
	logger("|      Welcome to TracebookMessenger!      |")
	logger("|        Listening for connections.        |")
	logger("--------------------------------------------")
}