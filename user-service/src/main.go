package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("my_secret_key")

const (
	dbConnectionString = "host=postgres user=postgres password=mysecretpassword dbname=userdb sslmode=disable"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
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
		break
	}
	if err != nil{
		log.Fatalf("Unable to connect to the database after %d attempts: %v", maxRetries, err)
	}
	log.Println("Connected to the database successfully!")
	return
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error saving user to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s registered successfully!", creds.Username)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=$1", creds.Username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	fmt.Fprintf(w, "Login successful!")
}

func main() {
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	log.Println("User service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
