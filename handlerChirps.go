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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	userId, _, err := cfg.authHelper(w, req)
	if err != nil {
		return
	}

	chirpId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing chirp id")
		errorRespHelper("Error getting chirp id", w, http.StatusBadRequest)
		return
	}
	chirpToDelete, err := cfg.db.GetChirpById(req.Context(), chirpId)
	if err != nil {
		log.Printf("Error getting chirp from database")
		errorRespHelper("Error getting chirp from database", w, http.StatusNotFound)
		return
	}
	if chirpToDelete.UserID != userId {
		errorRespHelper("Forbidden", w, http.StatusForbidden)
		return
	}
	err = cfg.db.DeleteChirp(req.Context(), chirpId)
	if err != nil {
		log.Printf("Error deleting chirp")
		errorRespHelper("Error deleting chirp", w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpIdHeader := req.PathValue("chirpID")
	chirpId, err := uuid.Parse(chirpIdHeader)
	if err != nil {
		log.Printf("Error parsing chirp id")
		errorRespHelper("Error getting chirp id", w, http.StatusBadRequest)
		return
	}
	chirp, err := cfg.db.GetChirpById(req.Context(), chirpId)
	if err != nil {
		log.Printf("Error getting chirp from database")
		errorRespHelper("Error getting chirp from database", w, http.StatusNotFound)
		return
	}

	respHelper(Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}, w, http.StatusOK)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	queryAuthorId := req.URL.Query().Get("author_id")
	sortOption := req.URL.Query().Get("sort")
	if sortOption == "" {
		sortOption = "asc"
	}
	if sortOption != "asc" && sortOption != "desc" {
		log.Printf("Error wrong sort option")
		errorRespHelper("Error wrong sort option", w, http.StatusBadRequest)
		return
	}
	var chirps []database.Chirp
	var err error
	if queryAuthorId != "" {
		authorId, err := uuid.Parse(queryAuthorId)
		if err != nil {
			log.Printf("Error parsing chirp id")
			errorRespHelper("Error getting chirp id", w, http.StatusBadRequest)
			return
		}
		chirps, err = cfg.db.GetChirpsByUserId(req.Context(), database.GetChirpsByUserIdParams{
			UserID:  authorId,
			Column2: sortOption,
		})
	} else {
		chirps, err = cfg.db.GetChirps(req.Context(), sortOption)
	}
	if err != nil {
		log.Printf("Error getting chirps from database: %v", err)
		errorRespHelper("Error getting chirps from database", w, http.StatusInternalServerError)
		return
	}

	responseChirps := []Chirp{}
	for _, chirp := range chirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
	respHelper(responseChirps, w, http.StatusOK)
}

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding request parameters: %v", err)
		errorRespHelper("Error decoding request parameters", w, http.StatusInternalServerError)
		return
	}
	userId, _, err := cfg.authHelper(w, req)
	if err != nil {
		return
	}

	if len(param.Body) > 140 {
		log.Printf("Chirp is too long")
		errorRespHelper("Chirp is too long", w, http.StatusBadRequest)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanBody := cleanBody(param.Body, badWords)

	chirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanBody,
		UserID:    userId,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		errorRespHelper("Something went wring", w, http.StatusInternalServerError)
		return
	}

	respHelper(Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}, w, http.StatusCreated)
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
