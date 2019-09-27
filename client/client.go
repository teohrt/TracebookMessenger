package client

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/teohrt/TracebookMessenger/dtos"
)

// InitialConnection connects to a peer to update chat history and peer network
func InitialConnection(node *dtos.Node) error {
	conn, err := net.Dial("tcp", node.FirstPeer)
	if err != nil {
		return err
	}

	binBuf := new(bytes.Buffer)
	gobobj := gob.NewEncoder(binBuf)

	if err := gobobj.Encode(&node); err != nil {
		return err
	}
	if _, err := conn.Write(binBuf.Bytes()); err != nil {
		return err
	}
	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func SendMessage(this *dtos.Node) {
	for {
		// Grab user input for message
		input := bufio.NewReader(os.Stdin)
		msg, _ := input.ReadString('\n')
		// Update this node's chat history and update peers
		msg = ("( " + this.Name + " ) : " + msg)
		this.ChatHistory = append(this.ChatHistory, msg)
		UpdateNetwork(this)
	}
}

// Sends update to every known node
func UpdateNetwork(this *dtos.Node) {
	for _, address := range this.PeerAddresses {
		// Don't try to update yourself
		if address != this.NodeAddress {
			updateSingleNode(address, this)
		}
	}
}

// Sends update to single peer
func updateSingleNode(address string, this *dtos.Node) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	// If connection was made
	binBuf := new(bytes.Buffer)
	gobobj := gob.NewEncoder(binBuf)
	if err := gobobj.Encode(this); err != nil {
		return err
	}
	conn.Write(binBuf.Bytes())
	return conn.Close()
}

// Add "server" messages to chat history and send to peers
func Log(msg string, node *dtos.Node) {
	fmt.Println(msg)
	node.ChatHistory = append(node.ChatHistory, msg+"\n")
	UpdateNetwork(node)
}
