package main

import (
	"encoding/json"
	"net/http"
)

// JSON helper function
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload) // Convert the payload to JSON

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Responding with Error helper function
func respondWithError(w http.ResponseWriter, code int, msg interface{}) {
	respondWithJSON(w, code, msg)
}
