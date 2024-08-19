package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Atviksord/BlogAggregator/internal/database"
	"github.com/google/uuid"
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
		return
	}
	err = cfg.autoFollowFeed(feed, user, r)
	if err != nil {
		http.Error(w, "Failed to auto-follow feed", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(feed)
	if err != nil {
		http.Error(w, "Failed to encode feed into json", http.StatusInternalServerError)
		return
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
		http.Error(w, "Failed to encode user into JSON", http.StatusInternalServerError)
	}

}

// Get specific feed in detail
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
	_, err := cfg.feedFollowHandlerHelper(r, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// get a list of all feeds followed by authenticated user
func (cfg *apiConfig) getFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	followList, err := cfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to get Follow Feeds from DB", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(followList)
	if err != nil {
		http.Error(w, "Failed to turn followfeed struct into json", http.StatusInternalServerError)
	}

}
func (cfg *apiConfig) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedId := r.PathValue("feed_id")
	if feedId == "" {
		http.Error(w, "No feed ID", http.StatusBadRequest)

	}
	trueFeedID, err := uuid.Parse(feedId)
	if err != nil {
		fmt.Printf("Failed to convert string to UUID %v", err)
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: trueFeedID,
		UserID: user.ID,
	})
	if err != nil {
		http.Error(w, "Failed to DeleteFeed %v", http.StatusInternalServerError)
	}
}
func (cfg *apiConfig) GetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) error {
	limitStr := r.URL.Query().Get("limit")
	limit := 5 //default if no query params

	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Error converting limit string to int", http.StatusInternalServerError)
			return err
		}
		posts, err := cfg.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
			UserID: user.ID,
			Limit:  int32(limitInt),
		})
		if err != nil {
			http.Error(w, "Error getting post by user", http.StatusInternalServerError)
		}

		w.WriteHeader(200)
		w.Write(posts)
		return nil
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
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.getFeedFollowHandler))
	mux.HandleFunc("DELETE /v1/feed_follows/{feed_id}", cfg.middlewareAuth(cfg.deleteFeedFollowHandler))
	mux.HandleFunc("GET /v1/posts", cfg.middlewareAuth(cfg.GetPostsByUser))

}
