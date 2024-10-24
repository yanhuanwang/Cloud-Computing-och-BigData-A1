# Expense Service API

This Expense Service API is built using Go (Golang) and PostgreSQL. It provides endpoints for creating, updating, listing, and deleting expense records. The service is accessible via HTTP REST API and supports CORS to allow access from different origins.

## Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- Docker (optional, for containerized deployment)

## Setup

### Step 1: Database Setup

Ensure you have a running PostgreSQL instance. Create a database named `expensedb`.

You can use the following SQL commands to create the database:

```sql
CREATE DATABASE expensedb;
```

The Go application will automatically create the `expenses` table if it does not exist.

### Step 2: Environment Configuration

Update the `dbConnectionString` constant in the Go code to reflect your database connection details:

```go
const (
    dbConnectionString = "host=postgres user=postgres password=mysecretpassword dbname=expensedb sslmode=disable"
)
```

### Step 3: Run the Application

You can run the application using the following command:

```sh
go run expense_service.go
```

The service will be available at `http://localhost:8081`.

## API Endpoints

### 1. Add Expense

**Endpoint:** `/add-expense`
**Method:** `POST`
**Request Body:**

```json
{
  "username": "string",
  "description": "string",
  "amount": "number"
}
```

**Response:**
- `200 OK`: Expense added successfully
- `400 Bad Request`: Invalid request payload
- `500 Internal Server Error`: Error saving expense to database

### 2. Update Expense

**Endpoint:** `/update-expense`
**Method:** `PUT`
**Request Body:**

```json
{
  "id": "integer",
  "username": "string",
  "description": "string",
  "amount": "number",
  "date": "string"
}
```

**Response:**
- `200 OK`: Expense updated successfully
- `400 Bad Request`: Invalid request payload
- `500 Internal Server Error`: Error updating expense

### 3. Get Expenses

**Endpoint:** `/get-expenses`
**Method:** `GET`
**Query Parameter:**

- `username`: The username for which to retrieve expenses

**Response:**
- `200 OK`: A list of expenses
- `400 Bad Request`: Username is required
- `500 Internal Server Error`: Error fetching expenses

### 4. Delete Expense

**Endpoint:** `/delete-expense`
**Method:** `DELETE`
**Request Body:**

```json
{
  "id": "integer",
  "username": "string"
}
```

**Response:**
- `200 OK`: Expense deleted successfully
- `400 Bad Request`: Invalid request payload
- `500 Internal Server Error`: Error deleting expense

## CORS Configuration

The service is configured to allow all origins, methods, and headers by default. You can modify the CORS settings in the `main` function if needed.

```go
c := cors.New(cors.Options{
    AllowedOrigins:   []string{"*"}, // Allow all origins
    AllowedMethods:   []string{"*"}, // Allow all HTTP methods
    AllowedHeaders:   []string{"*"}, // Allow all headers
    AllowCredentials: true,
})
```

## Running with Docker

You can also run the service using Docker. Create a `Dockerfile` with the following content:

```dockerfile
FROM golang:1.23
WORKDIR /app
COPY . .
RUN go build -o expense_service
CMD ["./expense_service"]
EXPOSE 8081
```

Build and run the Docker container:

```sh
docker build -t expense-service .
docker run -p 8081:8081 expense-service
```

## License

This project is licensed under the MIT License.

