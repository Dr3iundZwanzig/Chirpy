package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dr3iundZwanzig/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	reqStruct := request{}
	err := decoder.Decode(&reqStruct)
	if err != nil {
		log.Printf("Error decoding: %v", err)
		errorRespHelper("Error decoding JSON", w, http.StatusInternalServerError)
		return
	}
	user, err := cfg.db.GetUserByEmail(req.Context(), reqStruct.Email)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		errorRespHelper("Something went wring", w, http.StatusInternalServerError)
		return
	}

	ok, err := auth.CheckPasswordHash(reqStruct.Password, user.HashedPassword)
	if !ok {
		log.Printf("Wrong password: %v", err)
		errorRespHelper("Incorrect email or password", w, http.StatusUnauthorized)
		return
	}
	respHelper(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, w, http.StatusOK)
}
