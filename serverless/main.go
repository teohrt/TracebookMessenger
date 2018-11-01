package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Node struct {
	ChatHistory []string
	KnownNodes  []string
	NodeAddress string
}

var State = Node{}

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

	// Initialize node state
	thisAddress = localAddr.IP.String() + broadcastPort
	State = Node{ChatHistory: []string{}, KnownNodes: []string{}, NodeAddress: thisAddress}
	State.KnownNodes = append(State.KnownNodes, thisAddress)

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
		fmt.Print("Known nodes: ")
		fmt.Println(State.KnownNodes)

	} else {
		State.KnownNodes = append(State.KnownNodes, otherNodeAddress)

		binBuf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binBuf)
		gobobj.Encode(State)
		conn.Write(binBuf.Bytes())

		fmt.Println("Gob encoded and sent!")

		fmt.Print("Known nodes: ")
		fmt.Println(State.KnownNodes)

		conn.Close()
	}
}

// Listen and handle incomming tcp connections
func listen(thisPort string) {
	// Port format: ":8080"
	ln, _ := net.Listen("tcp", thisPort)

	// Decodes gobs recieved from every accepted connection
	for {
		conn, _ := ln.Accept()
		go decode(conn)
	}
}

// Decodes and prints gobs
func decode(conn net.Conn) {
	tmp := make([]byte, 500)
	_, _ = conn.Read(tmp)
	tmpBuf := bytes.NewBuffer(tmp)
	decodedStruct := new(Node)
	gobobj := gob.NewDecoder(tmpBuf)
	gobobj.Decode(decodedStruct)

	fmt.Println("Gob recieved from: " + decodedStruct.NodeAddress)

	// Update this node's state
	if len(decodedStruct.KnownNodes) > len(State.KnownNodes) {
		State.KnownNodes = decodedStruct.KnownNodes
		fmt.Println("My KnownNodes have been updated!")
		fmt.Print("Known nodes: ")
		fmt.Println(State.KnownNodes)
	}
	if len(decodedStruct.ChatHistory) > len(State.ChatHistory) {
		State.ChatHistory = decodedStruct.ChatHistory
		fmt.Println("My ChatHistory has been updated!")
	}

	// Send updates to new node
	if len(decodedStruct.KnownNodes) == 0 {
		sendUpdate(decodedStruct.NodeAddress)

		// Update previously existing nodes with new node's address
		// TODO
	}

	conn.Close()
}

// Change this to update every known node
// TODO
func sendUpdate(nodeAddress string) {
	conn, _ := net.Dial("tcp", nodeAddress)

	binBuf := new(bytes.Buffer)
	gobobj := gob.NewEncoder(binBuf)
	gobobj.Encode(State)
	conn.Write(binBuf.Bytes())

	fmt.Println("Update sent back to: " + nodeAddress)
	conn.Close()
}
