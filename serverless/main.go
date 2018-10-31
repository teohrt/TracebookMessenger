package main

import (
	"bufio"
	"fmt"
	"net"
)

type NodeInfo struct {
	chatHistory []string
	knownNodes  []string
}

var nodeInfo = NodeInfo{}

func main() {
	var exit = make(chan bool)

	var broadcastPort string
	var thisAddress string
	var initConnAddress string

	// Store this node's IP
	self, _ := net.Dial("udp", "8.8.8.8:80")
	localAddr := self.LocalAddr().(*net.UDPAddr)

	fmt.Print("Assign this node a port: ")
	fmt.Scanln(&broadcastPort)
	broadcastPort = ":" + broadcastPort
	//broadcastPort := ":8080"

	thisAddress = localAddr.IP.String() + broadcastPort

	fmt.Print("Enter address of node to connect to: ")
	fmt.Scanln(&initConnAddress)

	fmt.Println("Node hosted at : " + thisAddress)

	initialConnection(initConnAddress, thisAddress)

	go listen(broadcastPort)

	// Blocking operation
	// Allows the go routines to excecute indefinitely
	<-exit
}

// Connects to a node to get updated
// chat history and node network
func initialConnection(otherNodeAddress string, thisAddress string) {
	conn, err := net.Dial("tcp", otherNodeAddress)

	if err != nil {
		fmt.Println("Listening for connections...")

		// Set this node's data for testing purposes
		ch := []string{"test", "lol"}
		nodeInfo = NodeInfo{chatHistory: ch, knownNodes: ch}
		fmt.Println(nodeInfo)
	} else {
		fmt.Fprintf(conn, thisAddress+"\n")
		fmt.Println("Connection attempt made!")
	}
}

// Listen for incomming tcp requests
func listen(thisPort string) {
	// Port format: ":8080"
	ln, _ := net.Listen("tcp", thisPort)

	for {
		// Accept connection and grab msg
		conn, _ := ln.Accept()
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		msg = msg[:len(msg)-1] // strips the newline character from input

		fmt.Println("Message recieved from " + msg)
		// Add to list of clients
		// c := Client{conn, true, name}
		// clients = append(clients, c)

		//logAndSend("***New Client Connected : "+c.name+"\n", c)

		//go c.start()
	}
}
