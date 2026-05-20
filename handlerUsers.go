package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Dr3iundZwanzig/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerPostUsers(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Email string `json:"email"`
	}
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	reqStruct := request{}
	err := decoder.Decode(&reqStruct)
	if err != nil {
		log.Printf("Error decoding: %v", err)
		errorRespHelper("Something went wring", w, http.StatusInternalServerError)
		return
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     reqStruct.Email,
	}
	user, err := cfg.db.CreateUser(req.Context(), userParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		errorRespHelper("Something went wring", w, http.StatusInternalServerError)
		return
	}
	respHelper(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, w, http.StatusCreated)
}
