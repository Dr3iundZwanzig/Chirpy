package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	port := ":8080"
	filepath := "."
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))))
	serveMux.HandleFunc("/metrics", cfg.handlerHitCount)
	serveMux.HandleFunc("/healthz", handlerStatus)
	serveMux.HandleFunc("/reset", cfg.handlerReset)

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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
