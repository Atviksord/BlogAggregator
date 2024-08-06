package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Atviksord/BlogAggregator/internal/database"

	"github.com/google/uuid"
)

// struct to create user
type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}
type Parameters struct {
	Name string `json:"name"`
}

// CUSTOM TYPE FOR HANDLERS THAT REQUIRE AUTH
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func ReadynessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	respondWithError(w, 500, map[string]string{"error": "Internal server Error"})
}

// / DO CHECKS IF AUTHOR IS AUTHORIZED HERE, PROTECTED AUTHORIZED END POINT
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := cfg.extractAPIKey(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		user, err := cfg.userGetHelper(apiKey, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Pass control from AUTH middleware to main handler
		handler(w, r, user)

	}
}

// CREATE FEED
func (cfg *apiConfig) FeedCreateHandler(w http.ResponseWriter, r *http.Request, user database.User) {

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
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		fmt.Printf("Couldnt unmarshal r body into user struct %v", err)
		return

	}
	// Invoke helper function for creating user

	response, err := cfg.userCreateHelper(params, w, r)
	if err != nil {
		fmt.Printf("Failed to create user with helper function %v", err)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to write the response", http.StatusInternalServerError)
	}

}

// API Authorization Key in Header
func (cfg *apiConfig) UserGetHandler(w http.ResponseWriter, r *http.Request) {

	apiKey, err := cfg.extractAPIKey(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := cfg.userGetHelper(apiKey, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user into JSON", http.StatusUnauthorized)
	}

}

func (cfg *apiConfig) HandlerRegistry(mux *http.ServeMux) {
	fmt.Println("handlers being registered..")
	mux.HandleFunc("GET /v1/healthz", ReadynessHandler)
	mux.HandleFunc("GET /v1/err", ErrorHandler)
	mux.HandleFunc("POST /v1/users", cfg.UserCreateHandler)
	mux.HandleFunc("GET /v1/users", cfg.UserGetHandler)
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.FeedCreateHandler))

}
