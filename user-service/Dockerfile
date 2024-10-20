FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod, go.sum and source code files
COPY src/go.mod src/go.sum src/main.go ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o user-service

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./user-service"]
