# Step 1: Use a Go image as a base
FROM golang:1.23-alpine

# Step 2: Set the working directory in the container
WORKDIR /app

# Step 3: Copy the Go modules and vendor files
COPY go.mod go.sum ./

# Step 4: Download dependencies (or run `go mod vendor` if you have vendor dependencies)
RUN go mod download

# Step 5: Copy the entire source code
COPY . .

# Step 6: Build the Go application (assuming the entry point is in cmd/main.go)
RUN go build -o app ./cmd/main.go

# Step 7: Expose the application port (e.g., port 8000)
EXPOSE 8000

# Step 8: Set the entry point to run the compiled Go application
CMD ["./app"]
