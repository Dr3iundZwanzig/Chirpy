package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dr3iundZwanzig/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authHelper(w http.ResponseWriter, req *http.Request) (uuid.UUID, string, error) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		errorRespHelper("Unauthorized", w, http.StatusUnauthorized)
		return uuid.UUID{}, "", fmt.Errorf("Error getting token")
	}
	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		errorRespHelper("Unauthorized", w, http.StatusUnauthorized)
		return uuid.UUID{}, "", fmt.Errorf("Error validating token")
	}
	return userId, token, nil
}
