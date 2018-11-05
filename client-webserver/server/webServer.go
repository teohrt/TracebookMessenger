package main

import (
	"fmt"
	"net/http"
)

var history = []string{}

// Prints to the terminal and the web server
func logger(msg string) {
	fmt.Println(string(msg))
	history = append(history, msg)
}

// Displays chat history on web server
func showHistory(w http.ResponseWriter, r *http.Request) {
	for i, _ := range history {
		fmt.Fprintf(w, history[i]+"\n")
	}
}

// route for web server
func runWebServer(webPort string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", showHistory)
	http.ListenAndServe(webPort, mux)
}