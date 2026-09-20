package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// createServer returns an http.Server configured for a specific port/appID
func createServer(port string) *http.Server {
	mux := http.NewServeMux()

	// Endpoints definition
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve '/' strictly for root path matches
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_, err := fmt.Fprintf(w, "appid: %s home page: says hello!", port)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	})

	mux.HandleFunc("/app1", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "appid: %s app1 page: says hello!", port)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	})

	mux.HandleFunc("/app2", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "appid: %s app2 page: says hello!", port)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "appid: %s ADMIN page: very few people should see this", port)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func main() {
	ports := []string{"1111", "2222", "3333", "4444"}
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			server := createServer(p)
			fmt.Printf("Starting service on http://localhost:%s\n", p)
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Error starting server on port %s: %v", p, err)
			}
		}(port)
	}

	// Keep the main goroutine running until all servers exit
	wg.Wait()
}
