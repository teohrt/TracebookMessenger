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

	} else {
		// Data for testing purposes
		first := []string{"testing", "one", "two"}
		second := []string{"three", "lol", "gobs are cool"}
		nodeInfo := Node{ChatHistory: first, KnownNodes: second}

		//fmt.Fprintf(conn, thisAddress+"\n")
		binBuf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binBuf)
		gobobj.Encode(nodeInfo)
		conn.Write(binBuf.Bytes())

		fmt.Println("Gob encoded and sent!")
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
func decode(c net.Conn) {
	tmp := make([]byte, 500)
	_, _ = c.Read(tmp)
	tmpBuf := bytes.NewBuffer(tmp)
	tmpStruct := new(Node)
	gobobj := gob.NewDecoder(tmpBuf)
	gobobj.Decode(tmpStruct)

	fmt.Println(tmpStruct)
	c.Close()
}
