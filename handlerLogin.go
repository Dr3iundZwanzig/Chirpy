package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Dr3iundZwanzig/Chirpy/internal/auth"
	"github.com/Dr3iundZwanzig/Chirpy/internal/database"
	"github.com/google/uuid"
)

type LoginUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

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
	expirationTime := time.Duration(3600) * time.Second
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
	userToken, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		log.Printf("Error getting user token: %v", err)
		errorRespHelper("Error getting user token", w, http.StatusInternalServerError)
		return
	}
	refreshToken, err := cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(60*24) * time.Hour),
	})
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		errorRespHelper("Error creating refresh token", w, http.StatusInternalServerError)
		return
	}
	respHelper(LoginUser{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        userToken,
		RefreshToken: refreshToken.Token,
	}, w, http.StatusOK)
}
