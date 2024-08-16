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

	// Load environment variables

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load env variable %v", err)
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment variables")
	}

	// Attempt to connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to confirm the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// initialize database queries(SQL)
	dbQueries := database.New(db)
	cfg := &apiConfig{
		DB: dbQueries,
	}

	// Create server and launch it
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment variables")
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	// register endpoints
	cfg.HandlerRegistry(mux)
	// start webscraper worker go routine
	go cfg.FeedFetchWorker(10)
	log.Printf("Server is starting on port %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
