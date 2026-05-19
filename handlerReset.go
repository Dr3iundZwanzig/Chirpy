package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		errorRespHelper("Forbidden", w, 403)
		return
	}
	err := cfg.db.DeleteUser(req.Context())
	if err != nil {
		errorRespHelper("Error deleting users", w, 500)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
