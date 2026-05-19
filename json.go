package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResp struct {
	Error string `json:"error"`
}

func errorRespHelper(errorMsg string, w http.ResponseWriter, errCode int) {
	respBody := errorResp{
		Error: errorMsg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Fatalf("Error mashalling JSON: %v", err)
	}
	w.WriteHeader(errCode)
	w.Write(dat)
}

func respHelper(p interface{}, w http.ResponseWriter, code int) {
	dat, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error mashalling JSON: %v", err)
	}
	w.WriteHeader(code)
	w.Write(dat)
}
