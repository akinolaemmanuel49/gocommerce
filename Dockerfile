FROM golang:1.23-alpine

# Copy source code into container and get packages
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . .

# Build binary within container
RUN go build -o app ./cmd/main.go

# Expose port 8000
EXPOSE 8000

# Run binary
CMD ["./app"]
