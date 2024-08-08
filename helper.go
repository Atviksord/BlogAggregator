package main

import (
	"encoding/json"
	"fmt"
	"io"
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
func (cfg *apiConfig) userCreateHelper(params Parameters, r *http.Request) (database.User, error) {
	fmt.Printf("Inserting user %s", params.Name)
	USER, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		return USER, fmt.Errorf("failed to create user")

	}

	return USER, nil

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
func (cfg *apiConfig) userCreateFeedHelper(r *http.Request, user database.User) (database.Feed, error) {
	type FeedRequest struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	var feedrequest FeedRequest
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return database.Feed{}, fmt.Errorf("unable to read from body in createfeed")
	}
	err = json.Unmarshal(requestBody, &feedrequest)

	if err != nil {
		return database.Feed{}, fmt.Errorf("unable to marshal body, creating feed")
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedrequest.Name,
		Url:       feedrequest.URL,
		UserID:    user.ID,
	})
	if err != nil {
		return database.Feed{}, fmt.Errorf("unable to create feed")
	}
	return feed, nil
}
func (cfg *apiConfig) feedFollowHandlerHelper(r *http.Request, user database.User) (database.FeedsFollow, error) {
	type FeedParams struct {
		Feed_id uuid.UUID `json:"feed_id"`
	}
	feedParameters := FeedParams{}
	requestBody, err := io.ReadAll(r.Body)
	err = json.Unmarshal(requestBody, &feedParameters)
	if err != nil {
		return database.FeedsFollow{}, fmt.Errorf("Unable to read request body", http.StatusInternalServerError)

	}

	followFeeds, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedParameters.Feed_id,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return database.FeedsFollow{}, fmt.Errorf("Unable to CreateFeedFollow in handler", http.StatusInternalServerError)
	}
	return followFeeds, nil

}
