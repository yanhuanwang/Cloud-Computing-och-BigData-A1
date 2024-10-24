# User Service API

This User Service API is built using Go (Golang) and PostgreSQL. It provides endpoints for creating, editing, listing, and deleting users. The service is accessible via HTTP REST API and supports CORS to allow access from different origins.

## Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- Docker (optional, for containerized deployment)

## Setup

### Step 1: Database Setup

Ensure you have a running PostgreSQL instance. Create a database named `userdb`.

You can use the following SQL commands to create the database:

```sql
CREATE DATABASE userdb;
```

The Go application will automatically create the `users` table if it does not exist.

### Step 2: Environment Configuration

Update the `dbConnectionString` constant in the Go code to reflect your database connection details:

```go
const (
    dbConnectionString = "host=postgres user=postgres password=mysecretpassword dbname=userdb sslmode=disable"
)
```

### Step 3: Run the Application

You can run the application using the following command:

```sh
go run user_service.go
```

The service will be available at `http://localhost:8080`.

## API Endpoints

### 1. Create User

**Endpoint:** `/create-user`
**Method:** `POST`
**Request Body:**

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**
- `200 OK`: User created successfully
- `400 Bad Request`: Invalid request payload
- `500 Internal Server Error`: Error creating user

### 2. Edit User

**Endpoint:** `/edit-user`
**Method:** `PUT`
**Request Body:**

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**
- `200 OK`: User updated successfully
- `400 Bad Request`: Invalid request payload
- `500 Internal Server Error`: Error updating user

### 3. List Users

**Endpoint:** `/list-users`
**Method:** `GET`

**Response:**
- `200 OK`: A list of users
- `500 Internal Server Error`: Error fetching users

### 4. Delete User

**Endpoint:** `/delete-user`
**Method:** `DELETE`
**Request Body:**

```json
{
  "username": "string"
}
```

**Response:**
- `200 OK`: User deleted successfully
- `400 Bad Request`: Invalid request payload
- `500 Internal Server Error`: Error deleting user

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
RUN go build -o user_service
CMD ["./user_service"]
EXPOSE 8082
```

Build and run the Docker container:

```sh
docker build -t user-service .
docker run -p 8082:8082 user-service
```

## License

This project is licensed under the MIT License.

