FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod, go.sum and source code files
COPY src/go.mod src/go.sum src/main.go ./
COPY swagger-ui ./static/
COPY listExpense.html ./static/

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o expense-service

# Expose port 8081
EXPOSE 8081

# Command to run the executable
CMD ["./expense-service"]
