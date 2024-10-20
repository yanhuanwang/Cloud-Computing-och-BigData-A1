package main

import (
	"database/sql"
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
	dbConnectionString = "host=postgres user=postgres password=mysecretpassword dbname=expensedb sslmode=disable"
)

type Expense struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
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
			log.Printf("Unable to connect to the database after %d attempts: %v", maxRetries, err)
			time.Sleep(retryInterval)
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Unable to reach the database (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}
		break
	}
	if err !=nil{
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

func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	expense.Username = claims.Username
	_, err = db.Exec("INSERT INTO expenses (username, description, amount, date) VALUES ($1, $2, $3, $4)", expense.Username, expense.Description, expense.Amount, expense.Date)
	if err != nil {
		http.Error(w, "Error saving expense to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Expense added successfully!")
}

func GetExpensesHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query("SELECT id, username, description, amount, date FROM expenses WHERE username=$1", claims.Username)
	if err != nil {
		http.Error(w, "Error fetching expenses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.ID, &expense.Username, &expense.Description, &expense.Amount, &expense.Date); err != nil {
			http.Error(w, "Error scanning expense", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE expenses SET description=$1, amount=$2, date=$3 WHERE id=$4 AND username=$5", expense.Description, expense.Amount, expense.Date, expense.ID, claims.Username)
	if err != nil {
		http.Error(w, "Error updating expense", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Expense updated successfully!")
}

func DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var expenseID int
	if err := json.NewDecoder(r.Body).Decode(&expenseID); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM expenses WHERE id=$1 AND username=$2", expenseID, claims.Username)
	if err != nil {
		http.Error(w, "Error deleting expense", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Expense deleted successfully!")
}

func main() {
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/add-expense", AddExpenseHandler)
	http.HandleFunc("/get-expenses", GetExpensesHandler)
	http.HandleFunc("/update-expense", UpdateExpenseHandler)
	http.HandleFunc("/delete-expense", DeleteExpenseHandler)

	log.Println("Expense service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
