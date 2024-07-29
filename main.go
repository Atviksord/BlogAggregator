package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load env variable %v", err)
	}
	port := os.Getenv("PORT")
	// Create server and launch it
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	HandlerRegistry(mux)
	log.Printf("Server is starting on port %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
