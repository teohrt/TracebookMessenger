package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/teohrt/TracebookMessenger/client"
	"github.com/teohrt/TracebookMessenger/dtos"
)

// Listen handles incomming TCP connections
func Listen(this *dtos.Node) {
	ln, _ := net.Listen("tcp", this.NodeAddress)

	// Decodes gobs recieved from accepted connections
	for {
		conn, _ := ln.Accept()
		go decode(conn, this)
	}
}

// Decodes incomming gobs and updates nodes accordingly
func decode(conn net.Conn, this *dtos.Node) error {
	tmp := make([]byte, 500)
	if _, err := conn.Read(tmp); err != nil {
		return err
	}
	tmpBuf := bytes.NewBuffer(tmp)
	decodedStruct := new(dtos.Node)
	gobobj := gob.NewDecoder(tmpBuf)
	if err := gobobj.Decode(decodedStruct); err != nil {
		return err
	}

	// Add node if previously unknown
	if !addressIsKnown(decodedStruct.NodeAddress, this) {
		client.Log("New node connected: "+decodedStruct.Name, this)
		this.PeerAddresses = append(this.PeerAddresses, decodedStruct.NodeAddress)
	}

	// Update this node's state
	if len(decodedStruct.PeerAddresses) > len(this.PeerAddresses) {
		this.PeerAddresses = decodedStruct.PeerAddresses
	}
	if len(decodedStruct.ChatHistory) > len(this.ChatHistory) {

		previousLength := len(this.ChatHistory)
		// Update state
		this.ChatHistory = decodedStruct.ChatHistory

		// Print whole history if new node
		if previousLength == 0 {
			printChatHistory(this)
		} else {
			// Print recent chat history update
			fmt.Print(this.ChatHistory[len(this.ChatHistory)-1])
		}
	}

	// Send update to new node, and update every other known node
	if len(decodedStruct.PeerAddresses) < len(this.PeerAddresses) || len(decodedStruct.ChatHistory) < len(this.ChatHistory) {
		client.UpdateNetwork(this)
	}

	return conn.Close()
}

func addressIsKnown(addr string, this *dtos.Node) bool {
	for _, address := range this.PeerAddresses {
		if addr == address {
			return true
		}
	}
	return false

}

func printChatHistory(this *dtos.Node) {
	for i := range this.ChatHistory {
		fmt.Print(this.ChatHistory[i])
	}
}
