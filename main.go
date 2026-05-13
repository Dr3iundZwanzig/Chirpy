package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	port := ":8080"
	server := http.Server{
		Addr:    port,
		Handler: serveMux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
