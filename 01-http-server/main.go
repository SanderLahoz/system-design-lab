package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)

// Blocks all reading and writing
// whenever this mutex gets locked
var cacheMutex sync.RWMutex

func main() {
	// A mux is a request multiplexer which allows us to control
	// the traffic to different endpoints via handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)

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

	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve user id from path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve user from cache using id
	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Configure the response type to let the
	// client know that they will receive valid json
	w.Header().Set("Content-Type", "application/json")

	// Marshal user to json format
	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write back the json response (User)
	_, err = w.Write(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write back the successful status code
	w.WriteHeader(http.StatusOK)
}
