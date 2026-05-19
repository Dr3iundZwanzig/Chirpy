package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Dr3iundZwanzig/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding request parameters")
		errorRespHelper("Something went wrong", w, 500)
		return
	}

	if len(param.Body) > 140 {
		log.Printf("Chirp is too long")
		errorRespHelper("Chirp is too long", w, 400)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanBody := cleanBody(param.Body, badWords)

	chirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanBody,
		UserID:    param.UserId,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		errorRespHelper("Something went wring", w, 500)
		return
	}

	respHelper(Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}, w, 201)
}

func cleanBody(body string, badWords []string) string {
	splitBody := strings.Split(body, " ")
	for i, word := range splitBody {
		for _, badWord := range badWords {
			if strings.Contains(strings.ToLower(word), badWord) {
				splitBody[i] = "****"
			}
		}
	}
	return strings.Join(splitBody, " ")
}
