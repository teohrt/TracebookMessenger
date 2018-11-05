package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type Node struct {
	ChatHistory []string
	KnownNodes  []string
	NodeAddress string
	Name        string
}

var State = Node{}

func main() {
	var broadcastPort string
	var thisAddress string
	var initConnAddress string
	var name string

	// Store this node's IP
	self, _ := net.Dial("udp", "8.8.8.8:80")
	localAddr := self.LocalAddr().(*net.UDPAddr)

	fmt.Print("Assign this node a port: ")
	fmt.Scanln(&broadcastPort)
	broadcastPort = ":" + broadcastPort
	//broadcastPort := ":8080"

	fmt.Print("Enter address of peer: ")
	fmt.Scanln(&initConnAddress)
	fmt.Print("What is your name? : ")
	fmt.Scanln(&name)

	// Initialize node state
	thisAddress = localAddr.IP.String() + broadcastPort
	State = Node{ChatHistory: []string{}, KnownNodes: []string{}, NodeAddress: thisAddress, Name: name}
	State.KnownNodes = append(State.KnownNodes, thisAddress)

	fmt.Println("Node hosted at : " + thisAddress)
	fmt.Println("--------------------------------------------")
	fmt.Println("|      Welcome to TracebookMessenger!      |")
	fmt.Println("|        Listening for connections.        |")
	fmt.Println("|            Feel free to chat!            |")
	fmt.Println("--------------------------------------------")

	initialConnection(initConnAddress, thisAddress)

	go listen(broadcastPort)

	sendMessage()
}

// Connects to a node to get updated
// chat history and node network
func initialConnection(otherNodeAddress string, thisAddress string) {
	conn, err := net.Dial("tcp", otherNodeAddress)

	if err != nil {
		fmt.Println("Listening for connections...")
		//fmt.Print("Known nodes: ")
		//fmt.Println(State.KnownNodes)

	} else {
		State.KnownNodes = append(State.KnownNodes, otherNodeAddress)

		binBuf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binBuf)
		gobobj.Encode(State)
		conn.Write(binBuf.Bytes())

		//fmt.Println("Gob encoded and sent!")

		//fmt.Print("Known nodes: ")
		//fmt.Println(State.KnownNodes)

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

	// Add node if previously unknown
	if !addressIsKnown(decodedStruct.NodeAddress) {
		fmt.Println("New node connected: " + decodedStruct.Name)
		State.KnownNodes = append(State.KnownNodes, decodedStruct.NodeAddress)
	}
	//fmt.Println("Gob recieved from: " + decodedStruct.NodeAddress)

	// Update this node's state
	if len(decodedStruct.KnownNodes) > len(State.KnownNodes) {
		State.KnownNodes = decodedStruct.KnownNodes
		//fmt.Println("My KnownNodes have been updated!")
	}
	if len(decodedStruct.ChatHistory) > len(State.ChatHistory) {

		previousLength := len(State.ChatHistory)
		// Update state
		State.ChatHistory = decodedStruct.ChatHistory

		// Print whole history if new node
		if previousLength == 0 {
			printChatHistory()
		} else {
			// Print recent update
			fmt.Print(State.ChatHistory[len(State.ChatHistory)-1])
		}

	}

	// Send update to new node, and update every other known node
	if len(decodedStruct.KnownNodes) < len(State.KnownNodes) || len(decodedStruct.ChatHistory) < len(State.ChatHistory) {
		updateNetwork()
	}

	//fmt.Print("Known nodes: ")
	//fmt.Println(State.KnownNodes)

	conn.Close()
}

// Sends update to every known node
func updateNetwork() {
	for _, address := range State.KnownNodes {
		// Don't try to update yourself
		if address != State.NodeAddress {
			updateSingleNode(address)
		}
	}
}

// Sends update to single node
func updateSingleNode(address string) {
	conn, err := net.Dial("tcp", address)

	// If connection was made
	if err == nil {
		binBuf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binBuf)
		gobobj.Encode(State)
		conn.Write(binBuf.Bytes())

		//fmt.Println("Update sent back to: " + address)
		conn.Close()
	} else {
		//fmt.Println("Could not contact: " + address)
	}
}

// Returns true if argument address is in the state's slice of known nodes
func addressIsKnown(a string) bool {
	for _, address := range State.KnownNodes {
		if a == address {
			return true
		}
	}
	return false
}

func sendMessage() {
	for {
		// Grab user input for message
		input := bufio.NewReader(os.Stdin)
		msg, _ := input.ReadString('\n')
		// Update this node's chat history and update known nodes
		msg = ("( " + State.Name + " ) : " + msg)
		State.ChatHistory = append(State.ChatHistory, msg)
		updateNetwork()
	}
}

func printChatHistory() {
	for i := range State.ChatHistory {
		fmt.Print(State.ChatHistory[i])
	}
}
