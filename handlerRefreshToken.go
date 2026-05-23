package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/Dr3iundZwanzig/Chirpy/internal/auth"
	"github.com/Dr3iundZwanzig/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	w.Header().Set("Content-Type", "application/json")
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %v", err)
		errorRespHelper("Unauthorized", w, http.StatusUnauthorized)
		return
	}
	token, err := cfg.db.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		log.Printf("Error refresh token does not exist: %v", err)
		errorRespHelper("Refresh token does not exist", w, http.StatusUnauthorized)
		return
	}
	if token.RevokedAt.Valid || time.Now().After(token.ExpiresAt) {
		log.Printf("Token expired: %v", err)
		errorRespHelper("Token expired", w, http.StatusUnauthorized)
		return
	}
	expirationTime := time.Duration(3600) * time.Second
	newToken, err := auth.MakeJWT(token.UserID, cfg.secret, expirationTime)
	if err != nil {
		log.Printf("Error getting user token: %v", err)
		errorRespHelper("Error getting user token", w, http.StatusInternalServerError)
		return
	}
	respHelper(response{
		Token: newToken,
	}, w, http.StatusOK)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %v", err)
		errorRespHelper("Unauthorized", w, http.StatusUnauthorized)
		return
	}
	token, err := cfg.db.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		log.Printf("Error refresh token does not exist: %v", err)
		errorRespHelper("Refresh token does not exist", w, http.StatusUnauthorized)
		return
	}
	err = cfg.db.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		Token:     token.Token,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		errorRespHelper("Error revoking refresh token", w, http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
