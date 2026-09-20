package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const originServer = "127.0.0.1:8080"

func main() {

	// Start tcp server
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
	fmt.Println("Listening on port 3000")

	// Accept connections
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Failed to accept connection", err)
	}

	// Wrap connection with buffered IO
	// we can now read the data directly from this reader
	reader := bufio.NewReader(conn)

	// We read data from the reader into this buffer in memory
	buffer := make([]byte, 4096)
	_, err = reader.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection to buffer", err)
		return
	}

	// Contact origin server to establish connection
	OriginConn, err := net.Dial("tcp", originServer)
	if err != nil {
		log.Fatal("Failed to connect to originServer", err)
	}

	// Write the buffer to the origin server
	// (Forward message to the server through this proxy)
	_, err = OriginConn.Write(buffer)
	if err != nil {
		log.Fatal("Failed to write to originServer", err)
		return
	}

}
