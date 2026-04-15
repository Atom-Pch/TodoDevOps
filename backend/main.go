package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver anonymously
)

// Todo represents our database model and JSON structure
type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

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

	// GET: Fetch all To-Dos
	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, description, is_completed FROM todos")
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var t Todo
			// We use sql.NullString for description in case it's empty in the DB
			var desc sql.NullString
			if err := rows.Scan(&t.ID, &t.Title, &desc, &t.IsCompleted); err != nil {
				http.Error(w, "Failed to parse data", http.StatusInternalServerError)
				return
			}
			if desc.Valid {
				t.Description = desc.String
			}
			todos = append(todos, t)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	})

	// POST: Create a new To-Do
	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var t Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insert into DB and return the new ID
		err := db.QueryRow(
			"INSERT INTO todos (title, description) VALUES ($1, $2) RETURNING id",
			t.Title, t.Description,
		).Scan(&t.ID)

		if err != nil {
			http.Error(w, "Failed to create To-Do", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	})

	// 3. Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
