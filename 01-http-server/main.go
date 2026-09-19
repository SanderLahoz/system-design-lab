package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)

func main() {
	// A mux is a request multiplexer which allows us to control
	// the traffic to different endpoints via handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /users", createUser)

	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// Response writer is the status code for example 200, 404 etc...
// Request contains everything else: the body ; headers ; url
func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "hello World")
	if err != nil {
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "User name is required", http.StatusBadRequest)
	}

	userCache[len(userCache)+1] = user
	
	w.WriteHeader(http.StatusCreated)
}
