package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq" // Import the PostgreSQL driver anonymously
	"golang.org/x/crypto/bcrypt"
)

// Todo represents our database model and JSON structure
type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

// User represents user authentication model
type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"` // Used for registration
	Password string `json:"password"`
}

type contextKey string

const userIDKey = contextKey("user_id")

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the origin of the request (e.g., http://localhost:5173)
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// 2. CRITICAL: Allow cookies to be sent back and forth!
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, pass the request down to the actual router
		next.ServeHTTP(w, r)
	})
}

// Add this helper function to main.go
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized |"+err.Error(), http.StatusUnauthorized)
			return
		}

		// --- NEW LOGIC: Extract User ID and add to Context ---
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID := int(claims["user_id"].(float64))

		// Create a copy of the request with the new context attached
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		reqWithContext := r.WithContext(ctx)

		// Pass the new request down the chain!
		next.ServeHTTP(w, reqWithContext)
	}
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

	// Authentication endpoints
	mux.HandleFunc("POST /api/register", func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)

		// 1. Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 14)
		if err != nil {
			http.Error(w, "Server error |"+err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Insert user into database
		_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
			creds.Username, creds.Email, string(hashedPassword))

		if err != nil {
			http.Error(w, "Failed to create user (maybe username/email exists?) |"+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "User registered successfully"}`))
	})

	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		json.NewDecoder(r.Body).Decode(&creds)

		// 1. Fetch the user's hashed password and ID from the DB
		var storedHash string
		var userID int
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE username=$1", creds.Username).Scan(&userID, &storedHash)
		if err != nil {
			http.Error(w, "Invalid credentials |"+err.Error(), http.StatusUnauthorized)
			return
		}

		// 2. Compare the stored hash with the provided password
		if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid credentials |"+err.Error(), http.StatusUnauthorized)
			return
		}

		// 3. Create the JWT Token (The VIP Wristband)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
		})

		jwtSecret := os.Getenv("JWT_SECRET")
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		// 4. Set the token in an HTTP-Only Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: true,  // JavaScript cannot read this!
			Secure:   false, // Set to true in production (requires HTTPS)
			Path:     "/",
			SameSite: http.SameSiteLaxMode, // Important for CORS/SvelteKit
		})

		w.Write([]byte(`{"message": "Logged in successfully"}`))
	})

	mux.HandleFunc("POST /api/logout", func(w http.ResponseWriter, r *http.Request) {
		// To log out, we just overwrite the cookie with an expired one
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})
		w.Write([]byte(`{"message": "Logged out"}`))
	})

	// --- GET CURRENT USER (/me) ---
	mux.HandleFunc("GET /api/me", func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Not logged in |"+err.Error(), http.StatusUnauthorized)
			return
		}

		// 2. Parse the token
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid session |"+err.Error(), http.StatusUnauthorized)
			return
		}

		// 3. Extract the user_id from the token payload (claims)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// JWT stores numbers as float64, so we cast it to an int
		userID := int(claims["user_id"].(float64))

		// 4. Fetch the username from the database
		var username string
		err = db.QueryRow("SELECT username FROM users WHERE id=$1", userID).Scan(&username)
		if err != nil {
			http.Error(w, "User not found |"+err.Error(), http.StatusNotFound)
			return
		}

		// 5. Return the username to SvelteKit
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"username": username,
		})
	})

	// App endpoints
	// GET: Fetch all To-Dos
	mux.HandleFunc("GET /api/todos", requireAuth(
		func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(userIDKey).(int)
			rows, err := db.Query("SELECT id, title, description, is_completed FROM todos WHERE user_id=$1", userID)

			if err != nil {
				http.Error(w, "Failed to query database |"+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var todos []Todo
			for rows.Next() {
				var t Todo
				// We use sql.NullString for description in case it's empty in the DB
				var desc sql.NullString
				if err := rows.Scan(&t.ID, &t.Title, &desc, &t.IsCompleted); err != nil {
					http.Error(w, "Failed to parse data |"+err.Error(), http.StatusInternalServerError)
					return
				}
				if desc.Valid {
					t.Description = desc.String
				}
				todos = append(todos, t)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos)
		}))

	// POST: Create a new To-Do
	mux.HandleFunc("POST /api/todos", requireAuth(
		func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(userIDKey).(int)

			var t Todo
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, "Invalid request payload |"+err.Error(), http.StatusBadRequest)
				return
			}

			// Insert into DB and return the new ID
			err := db.QueryRow(
				"INSERT INTO todos (user_id, title, description) VALUES ($1, $2, $3) RETURNING id",
				userID, t.Title, t.Description,
			).Scan(&t.ID)

			if err != nil {
				http.Error(w, "Failed to create To-Do |"+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(t)
		}))

	// DELETE: Remove a To-Do
	// DELETE: Remove a To-Do
	mux.HandleFunc("DELETE /api/todos/{id}",
		requireAuth(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extract the User ID from the context (Security check!)
			userID := r.Context().Value(userIDKey).(int)

			// 2. Get the To-Do ID from the URL path
			todoID := r.PathValue("id")
			if todoID == "" {
				http.Error(w, "Missing To-Do ID", http.StatusBadRequest)
				return
			}
			
			var todoTitle string
			err = db.QueryRow("SELECT title FROM todos WHERE id = $1", todoID).Scan(&todoTitle)
			if err != nil {
				todoTitle = "Unknown |" + err.Error() // Fallback title if we can't fetch it (e.g., already deleted)
			}

			// 3. Execute the delete query
			// CRITICAL: We include 'AND user_id = $2' to ensure they can only delete their own tasks!
			result, err := db.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", todoID, userID)
			if err != nil {
				http.Error(w, "Failed to delete To-Do |"+err.Error(), http.StatusInternalServerError)
				return
			}

			// 4. Check if a row was actually deleted
			rowsAffected, err := result.RowsAffected()
			if err != nil || rowsAffected == 0 {
				// If 0 rows were affected, either the ID doesn't exist, or it belongs to another user
				http.Error(w, "To-Do not found or unauthorized", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "To-Do '` + todoTitle + `' deleted successfully"}`))
		}))

	// 3. Start the Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}
