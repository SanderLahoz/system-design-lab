package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const originServer = "127.0.0.1:8080"

func main() {
	// Start TCP listener on port 3000
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
	// Close connection when main exits
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("Failed to close listener: ", err)
		}
	}(listener)

	fmt.Println("Proxy listening on port 3000...")

	// Infinite event loop that waits for client connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection: ", err)
			continue
		}

		// Handle each client separately in their own goroutine
		go handleConnection(clientConn)
	}
}

// A function to handle a client connection
func handleConnection(clientConn net.Conn) {
	// Close connection when function exits
	defer func(clientConn net.Conn) {
		err := clientConn.Close()
		if err != nil {
			log.Println("Failed to close connection: ", err)
		}
	}(clientConn)

	// Establish connection to origin server (8080)
	originConn, err := net.Dial("tcp", originServer)
	if err != nil {
		log.Println("Failed to connect to originServer: ", err)
		return
	}
	// Close connection when function exits
	defer func(originConn net.Conn) {
		err := originConn.Close()
		if err != nil {
			log.Println("Failed to close originServer: ", err)
		}
	}(originConn)

	// Send client request -> origin server (in a goroutine)
	go func() {
		_, _ = io.Copy(originConn, clientConn)
	}()

	// Send origin response -> client
	_, _ = io.Copy(clientConn, originConn)
}
