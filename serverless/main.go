package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	var exit = make(chan bool)

	var broadcastPort string
	var nodeAddress string

	self, _ := net.Dial("udp", "8.8.8.8:80")
	localAddr := self.LocalAddr().(*net.UDPAddr)

	fmt.Print("Enter port for this node: ")
	fmt.Scanln(&broadcastPort)
	broadcastPort = ":" + broadcastPort
	//broadcastPort := ":8080"

	fmt.Print("Enter address of node to connect to: ")
	fmt.Scanln(&nodeAddress)
	fmt.Println("Node hosted at : " + localAddr.IP.String() + broadcastPort)

	initialConnection(nodeAddress)

	go listen(broadcastPort)

	// Blocking operation
	// Allows the go routines to excecute indefinitely
	<-exit
}

// Connects to a node to get updated
// chat history and node network
func initialConnection(address string) {
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println("no connection")
	} else {
		fmt.Fprintf(conn, "TEST!\n")
	}
}

func listen(thisPort string) {
	// Port format: ":8080"
	// Listen for incomming tcp requests
	ln, _ := net.Listen("tcp", thisPort)

	for {
		// Accept connection and grab msg
		conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		msg = msg[:len(msg)-1] // strips the newline character from input

		fmt.Println(msg)
		// Add to list of clients
		// c := Client{conn, true, name}
		// clients = append(clients, c)

		//logAndSend("***New Client Connected : "+c.name+"\n", c)

		//go c.start()
	}
}
