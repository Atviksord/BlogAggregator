package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// struct to create user
type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}
type Parameters struct {
	Name string `json:"name"`
}

func ReadynessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	respondWithError(w, 500, map[string]string{"error": "Internal server Error"})
}
func (cfg *apiConfig) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	params := Parameters{}
	err = json.Unmarshal(requestBody, &params)
	if err != nil {
		fmt.Errorf("Couldnt unmarshal r body into user struct %v", err)
	}
	// Invoke helper function for creating user

	err, response := cfg.userCreateHelper(params, w, r)
	if err != nil {
		fmt.Printf("Failed to create user with helper function")
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to write the response", http.StatusInternalServerError)
	}

}

func (cfg *apiConfig) HandlerRegistry(mux *http.ServeMux) {
	fmt.Println("handlers being registered..")
	// TEST HANDLER
	mux.HandleFunc("GET /v1/healthz", ReadynessHandler)
	mux.HandleFunc("GET /v1/err", ErrorHandler)
	mux.HandleFunc("POST /v1/users", cfg.UserCreateHandler)
}
