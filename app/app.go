package app

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/teohrt/TracebookMessenger/client"
	"github.com/teohrt/TracebookMessenger/dtos"
	"github.com/teohrt/TracebookMessenger/server"
)

// Start is the entry point of the application responsible for initialization
func Start(node *dtos.Node) {
	if err := validateNode(node); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := initNode(node); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(node.FirstPeer) > 1 {
		if err := client.InitialConnection(node); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("--------------------------------------------")
	fmt.Println("Node hosted at : " + node.NodeAddress)
	fmt.Println("--------------------------------------------")
	fmt.Println("|      Welcome to TracebookMessenger!      |")
	fmt.Println("|        Listening for connections.        |")
	fmt.Println("|            Feel free to chat!            |")
	fmt.Println("--------------------------------------------")

	go server.Listen(node)

	client.SendMessage(node)
}

func initNode(node *dtos.Node) error {
	self, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err
	}
	localAddr := self.LocalAddr().(*net.UDPAddr)
	thisAddress := localAddr.IP.String() + ":" + node.Port
	node.NodeAddress = thisAddress
	node.PeerAddresses = append(node.PeerAddresses, node.FirstPeer)

	return nil
}

func validateNode(node *dtos.Node) error {
	if node.Port == "" {
		return errors.New("must provide a port for the node to run on")
	}

	if len(node.Name) < 1 {
		return errors.New("must specify a chat name")
	}

	return nil
}
