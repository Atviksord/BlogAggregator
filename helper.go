package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Atviksord/BlogAggregator/internal/database"
	"github.com/google/uuid"
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

// Helps the Create User Handler with creation logic
func (cfg *apiConfig) userCreateHelper(params Parameters, w http.ResponseWriter, r *http.Request) (createUserResponse, error) {
	fmt.Printf("Inserting user %s", params.Name)
	USER, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return createUserResponse{}, err
	}
	response := createUserResponse{
		ID:        USER.ID,
		CreatedAt: USER.CreatedAt,
		UpdatedAt: USER.UpdatedAt,
		Name:      USER.Name,
	}
	return response, nil

}
func (cfg *apiConfig) userGetHelper(apiKey string, w http.ResponseWriter, r *http.Request) (createUserResponse, error) {
	USER, err := cfg.DB.GetApi(r.Context(), apiKey)
	if err != nil {
		http.Error(w, "Error getting user by API key")
	}
	response := createUserResponse{
		ID:        USER.ID,
		CreatedAt: USER.CreatedAt,
		UpdatedAt: USER.UpdatedAt,
		Name:      USER.Name,
		ApiKey:    USER.ApiKey,
	}
	return response, nil
}
