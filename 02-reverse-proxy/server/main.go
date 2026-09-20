package main

import (
	"fmt"
	"log"
	"net/http"
)

// Simple http server listening on localhost
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello World")
		if err != nil {
			log.Fatal(err)
		}
	})

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
