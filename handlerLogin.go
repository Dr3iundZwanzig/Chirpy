package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Dr3iundZwanzig/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type LoginUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
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
	expirationTime := 3600
	if reqStruct.ExpiresInSeconds != nil && *reqStruct.ExpiresInSeconds < 3600 {
		expirationTime = *reqStruct.ExpiresInSeconds
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
	userToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expirationTime)*time.Second)
	if err != nil {
		log.Printf("Error getting user token: %v", err)
		errorRespHelper("Something went wring", w, http.StatusInternalServerError)
		return
	}
	respHelper(LoginUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     userToken,
	}, w, http.StatusOK)
}
