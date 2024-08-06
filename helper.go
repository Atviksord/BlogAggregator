package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
func (cfg *apiConfig) extractAPIKey(r *http.Request) (string, error) {
	requestHeader := r.Header.Get("Authorization")
	if requestHeader == "" {
		return "", fmt.Errorf("missing Authorization Header")

	}
	parts := strings.Split(requestHeader, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", fmt.Errorf("invalid Authorization format")
	}
	return parts[1], nil
}

func (cfg *apiConfig) userGetHelper(apiKey string, r *http.Request) (database.User, error) {
	USER, err := cfg.DB.GetApi(r.Context(), apiKey)
	if err != nil {
		return database.User{}, fmt.Errorf("error getting user by API key In helper")
	}

	return USER, nil
}
