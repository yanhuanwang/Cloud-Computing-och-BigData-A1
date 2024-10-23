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
	dbConnectionString = "host=postgres user=postgres password=mysecretpassword dbname=expensedb sslmode=disable"
)

type Expense struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
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

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			description VARCHAR(255) NOT NULL,
			amount NUMERIC(10, 2) NOT NULL,
			date TIMESTAMP NOT NULL
		)`)
		if err != nil {
			log.Fatalf("Unable to create expenses table: %v", err)
		}

		break
	}
	if err != nil {
		log.Fatalf("Unable to connect to the database after %d attempts: %v", maxRetries, err)
	}
	log.Println("Connected to the database and ensured expenses table exists successfully!")
	return
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		lrw := &loggingResponseWriter{w, http.StatusOK}
		handler.ServeHTTP(lrw, r)
		log.Printf("Completed %d %s in %v", lrw.statusCode, http.StatusText(lrw.statusCode), time.Since(start))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
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
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	expense.Date = time.Now()
	_, err := db.Exec("INSERT INTO expenses (username, description, amount, date) VALUES ($1, $2, $3, $4)", expense.Username, expense.Description, expense.Amount, expense.Date)
	if err != nil {
		http.Error(w, "Error saving expense to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Expense added successfully!")
}

func GetExpensesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, username, description, amount, date FROM expenses WHERE username=$1", username)
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
	enableCors(&w)
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE expenses SET description=$1, amount=$2, date=$3 WHERE id=$4 AND username=$5", expense.Description, expense.Amount, expense.Date, expense.ID, expense.Username)
	if err != nil {
		http.Error(w, "Error updating expense", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Expense updated successfully!")
}

func DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var expenseID int
	if err := json.NewDecoder(r.Body).Decode(&expenseID); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM expenses WHERE id=$1 AND username=$2", expenseID, username)
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

	// Configure CORS to allow all origins
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Allowed headers
		AllowCredentials: true,
	})

	// Use the CORS handler
	corsHandler := c.Handler(logRequest(http.DefaultServeMux))

	log.Println("Expense service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsHandler))
}
