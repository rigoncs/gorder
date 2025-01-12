package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	log.Println("Listening :8082")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.RequestURI)
		_, _ = io.WriteString(w, "<h1>Welcome to Home Page</h1>")
	})
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.RequestURI)
		_, _ = io.WriteString(w, "pong")
	})
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
