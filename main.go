package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Atviksord/BlogAggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load Database
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error loading Postgres database")
	}
	dbQueries := database.New(db)
	// Load environment variables

	err = godotenv.Load()
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
