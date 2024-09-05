package main

import (
	"fmt"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"Received request: %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(
				w,
				"Server does not support hijacking",
				http.StatusInternalServerError,
			)
			return
		}

		next.ServeHTTP(w, r)

		log.Printf("Responded to %s with status: %d", r.RemoteAddr, http.StatusOK)

		conn, _, err := hj.Hijack()
		if err != nil {
			log.Printf("Failed to hijack connection: %v", err)
			return
		}
		defer conn.Close()

		log.Printf("Connection to %s closed", r.RemoteAddr)
	})
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorldHandler)
	loggedMux := loggingMiddleware(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
