package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type resp struct {
	Valid bool `json:"valid"`
}

type parameters struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(fmt.Sprintf(`
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}

func handlerStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidate(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Something went wrong")
		errorRespHelper("Something went wrong", w, 500)
		return
	}

	if len(param.Body) > 140 {
		log.Printf("Chirp is too long")
		errorRespHelper("Chirp is too long", w, 400)
		return
	}
	respHelper(resp{
		Valid: true,
	}, w, 200)
}
