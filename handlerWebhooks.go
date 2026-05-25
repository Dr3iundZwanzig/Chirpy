package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dr3iundZwanzig/Chirpy/internal/auth"
	"github.com/Dr3iundZwanzig/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		errorRespHelper("Error no api key in header", w, http.StatusUnauthorized)
		return
	}
	if apiKey != cfg.polkaApiKey {
		errorRespHelper("Error wrong api key", w, http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(req.Body)
	reqStruct := request{}
	err = decoder.Decode(&reqStruct)
	if err != nil {
		log.Printf("Error decoding: %v", err)
		errorRespHelper("Error decoding JSON", w, http.StatusInternalServerError)
		return
	}
	if reqStruct.Event != "user.upgraded" {
		errorRespHelper("Error wrong event", w, http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(reqStruct.Data.UserID)
	if err != nil {
		log.Printf("Error decoding: %v", err)
		errorRespHelper("Error user id has wrong format", w, http.StatusNotFound)
		return
	}
	err = cfg.db.UpdateChirpyRed(req.Context(), database.UpdateChirpyRedParams{
		ID:          userID,
		IsChirpyRed: true,
	})
	if err != nil {
		log.Printf("Error decoding: %v", err)
		errorRespHelper("Error user id does not exist", w, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
