package main

import (
	"fmt"
	"net/http"
)

func ReadynessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	respondWithError(w, 500, map[string]string{"error": "Internal server Error"})
}
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	respondWithError(w, 500, map[string]string{"error": "Internal server Error"})
}

func HandlerRegistry(mux *http.ServeMux) {
	fmt.Println("handlers being registered..")
	// TEST HANDLER
	mux.HandleFunc("GET /v1/healthz", ReadynessHandler)
	mux.HandleFunc("GET /v1/err", ErrorHandler)
	mux.HandleFunc("POST /v1/users", UserCreateHandler)
}
