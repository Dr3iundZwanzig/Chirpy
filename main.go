package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Dr3iundZwanzig/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}
	dburl, ok := os.LookupEnv("DB_URL")
	if !ok {
		log.Fatalln("dburl empty")
	}
	db, err := sql.Open("postgres", dburl)
	platform := os.Getenv("PLATFORM")
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}
	dbQueries := database.New(db)
	port := ":8080"
	filepath := "."
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))))
	serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serveMux.HandleFunc("GET /api/healthz", handlerStatus)
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serveMux.HandleFunc("POST /api/users", cfg.handlerUsers)
	serveMux.HandleFunc("POST /api/chirps", cfg.handlerChirps)

	server := http.Server{
		Addr:    port,
		Handler: serveMux,
	}
	fmt.Println("Server started at port" + port)
	err = server.ListenAndServe()
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
