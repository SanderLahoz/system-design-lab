package main

import (
	"fmt"
	"net/http"
)

func main() {
	// A mux is a request multiplexer which allows us to control
	// the traffic to different endpoints via handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

// Response writer is the status code for example 200, 404 etc...
// Request contains everything else: the body ; headers ; url
func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello World")
	if err != nil {
		return
	}
}
