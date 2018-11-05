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
	peers       []string
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
	State = Node{ChatHistory: []string{}, peers: []string{}, NodeAddress: thisAddress, Name: name}
	State.peers = append(State.peers, thisAddress)

	fmt.Println("--------------------------------------------")
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
func initialConnection(peerAddress string, thisAddress string) {
	conn, err := net.Dial("tcp", peerAddress)

	if err != nil {
		fmt.Println("Listening for connections...")

	} else {
		State.peers = append(State.peers, peerAddress)

		binBuf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binBuf)
		gobobj.Encode(State)
		conn.Write(binBuf.Bytes())
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
		publicServiceAnouncement("New node connected: " + decodedStruct.Name)
		State.peers = append(State.peers, decodedStruct.NodeAddress)
	}

	// Update this node's state
	if len(decodedStruct.peers) > len(State.peers) {
		State.peers = decodedStruct.peers
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
	if len(decodedStruct.peers) < len(State.peers) || len(decodedStruct.ChatHistory) < len(State.ChatHistory) {
		updateNetwork()
	}

	conn.Close()
}

// Sends update to every known node
func updateNetwork() {
	for _, address := range State.peers {
		// Don't try to update yourself
		if address != State.NodeAddress {
			updateSingleNode(address)
		}
	}
}

// Sends update to single peer
func updateSingleNode(address string) {
	conn, err := net.Dial("tcp", address)

	// If connection was made
	if err == nil {
		binBuf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(binBuf)
		gobobj.Encode(State)
		conn.Write(binBuf.Bytes())
		conn.Close()
	} else {
		//fmt.Println("Could not contact: " + address)
	}
}

// Returns true if argument address is in the state's slice of peers
func addressIsKnown(a string) bool {
	for _, address := range State.peers {
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
		// Update this node's chat history and update peers
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

// Add "server" messages to chat history and send to peers
func publicServiceAnouncement(msg string) {
	fmt.Println(msg)
	State.ChatHistory = append(State.ChatHistory, msg+"\n")
	updateNetwork()
}
