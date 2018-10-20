package main

func main() {

	// Chat server port
	chatPort := ":8080"
	// Web server port
	webPort := ":3000"

	welcomeMessage(chatPort, webPort)

	// Run the web server
	go runWebServer(webPort)
	
	// Handle outgoing server messages from console
	go handleOutgoing()

	runChatServer(chatPort)
}