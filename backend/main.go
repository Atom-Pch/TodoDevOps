package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver anonymously
)

func main() {
	// 1. Database Connection Configuration
	// We read these from the environment. Later, Docker and Terraform will inject these!
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	// Default to localhost if DB_HOST isn't set (useful for quick local testing)
	if dbHost == "" {
		dbHost = "localhost"
	}

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure our credentials actually work
	if err := db.Ping(); err != nil {
		log.Printf("Warning: Database ping failed (Is Postgres running?): %v", err)
	} else {
		fmt.Println("Successfully connected to PostgreSQL!")
	}

	// 2. Setup API Routes
	// We use the standard HTTP multiplexer (router)
	mux := http.NewServeMux()

	// Health check endpoint - This is CRITICAL for DevOps!
	// AWS Load Balancers and Docker will use this to check if your app is alive.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy and running!"))
	})

	// Placeholder routes for our To-Do features
	mux.HandleFunc("GET /todo-list", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This will return a list of To-Dos"))
	})

	// 3. Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
