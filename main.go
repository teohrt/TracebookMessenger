package main

import (
	"flag"

	"github.com/teohrt/TracebookMessenger/app"
	"github.com/teohrt/TracebookMessenger/dtos"
)

func main() {
	chatName := flag.String("name", "anon", "The name you will be seen as in chat")
	port := flag.String("port", "1234", "The port exposed for chat, Defaults to 1234")
	peerAddress := flag.String("peerAddress", "", "The peer address that gets you into the party")
	flag.Parse()

	app.Start(&dtos.Node{
		Name:      *chatName,
		Port:      *port,
		FirstPeer: *peerAddress,
	})
}
