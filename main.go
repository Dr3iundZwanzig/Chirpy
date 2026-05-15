package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":8080"
	filepath := "."

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepath))))
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))

	})

	server := http.Server{
		Addr:    port,
		Handler: serveMux,
	}
	fmt.Println("Server started at port" + port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
