package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Atviksord/BlogAggregator/internal/database"
)

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
	feed, err := cfg.userCreateFeedHelper(r, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(feed)

	if err != nil {
		http.Error(w, "Failed to encode feed into json", http.StatusInternalServerError)
	}
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

	response, err := cfg.userCreateHelper(params, r)
	if err != nil {
		http.Error(w, "Failed to create user with helper function %v", http.StatusUnauthorized)

	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to write the response", http.StatusInternalServerError)
	}

}

func (cfg *apiConfig) UserGetHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user into JSON", http.StatusUnauthorized)
	}

}
func (cfg *apiConfig) feedGetHandler(w http.ResponseWriter, r *http.Request) {
	allfeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		http.Error(w, "Failed to get feeds", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(allfeeds)
	if err != nil {
		http.Error(w, "Failed to turn feed struct into json", http.StatusInternalServerError)
	}
}
func (cfg *apiConfig) feedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollow, err := cfg.feedFollowHandlerHelper(r, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) HandlerRegistry(mux *http.ServeMux) {
	fmt.Println("handlers being registered..")
	mux.HandleFunc("GET /v1/healthz", ReadynessHandler)
	mux.HandleFunc("GET /v1/err", ErrorHandler)
	mux.HandleFunc("POST /v1/users", cfg.UserCreateHandler)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.UserGetHandler))
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.FeedCreateHandler))
	mux.HandleFunc("GET /v1/feeds", cfg.feedGetHandler)
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.feedFollowHandler))

}
