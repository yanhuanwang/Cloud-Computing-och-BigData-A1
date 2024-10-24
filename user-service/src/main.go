package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	dbConnectionString = "host=postgres user=postgres password=mysecretpassword dbname=userdb sslmode=disable"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

var db *sql.DB

func init() {
	retryInterval := 5 * time.Second
	maxRetries := 10
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dbConnectionString)
		if err != nil {
			log.Fatalf("Unable to connect to the database after %d attempts: %v", maxRetries, err)
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Unable to reach the database (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100) NOT NULL,
			password VARCHAR(100) NOT NULL
		)`) // Ensuring the users table exists
		if err != nil {
			log.Fatalf("Unable to create users table: %v", err)
		}
		break
	}
	if err != nil {
		log.Fatalf("Unable to connect to the database after %d attempts: %v", maxRetries, err)
	}
	log.Println("Connected to the database and ensured users table exists successfully!")
	return
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User created successfully!")
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE users SET email=$1, password=$2 WHERE username=$3", user.Email, user.Password, user.Username)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User updated successfully!")
}

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM users WHERE username=$1", user.Username)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User deleted successfully!")
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}
	if err := db.Ping(); err != nil {
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
func main() {
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/create-user", CreateUserHandler)
	http.HandleFunc("/edit-user", EditUserHandler)
	http.HandleFunc("/list-users", ListUsersHandler)
	http.HandleFunc("/delete-user", DeleteUserHandler)

	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Configure CORS to allow all origins, methods, and headers
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"*"}, // Allow all HTTP methods
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
	})

	// Use the CORS handler
	corsHandler := c.Handler(http.DefaultServeMux)

	log.Println("User service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8082", corsHandler))
}
