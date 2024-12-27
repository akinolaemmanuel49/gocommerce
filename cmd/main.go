package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	l "github.com/akinolaemmanuel49/gocommerce/log"
	"github.com/akinolaemmanuel49/gocommerce/routes"
	"github.com/gorilla/mux"

	"github.com/akinolaemmanuel49/gocommerce/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	RouteHealth = "/health"
)

func main() {
	// Setup logger
	logger := l.SetupLogger("service.log")

	// Load config
	config, err := configs.LoadConfig(".")
	if err != nil {
		logger.Fatalf("Error reading config file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB
	db, err := database.ConnectMongoDB(config.MongoDBURI, config.MongoDBName)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = db.Client().Disconnect(ctx); err != nil {
			logger.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to ensure a successful connection
	if err := db.Client().Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		logger.Fatalf("Failed to ping MongoDB: %v", err)
	}
	logger.Println("Connected to MongoDB successfully")

	// Initialize RabbitMQ
	conn, ch, err := queue.ConnectRabbitMQ(&config, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			logger.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}()
	defer func() {
		if err = ch.Close(); err != nil {
			logger.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}()

	// Start consuming messages from RabbitMQ
	go queue.ConsumeOrderNotifications(&config, ch)

	// Setup HTTP routes
	router := mux.NewRouter()
	setupRoutes(router, db, logger)

	// Start the HTTP server with graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router,
	}
	go func() {
		logger.Printf("Server is running on http://localhost:%s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

func setupRoutes(router *mux.Router, db *mongo.Database, logger *log.Logger) {
	// Health check
	router.HandleFunc(RouteHealth, healthHandler)

	// User routes
	routes.RegisterUserRoutes(router, db, logger)

	// Product routes
	routes.RegisterProductRoutes(router, db, logger)

	// Order routes
	routes.RegisterOrderRoutes(router, db, logger)

	// Category routes
	routes.RegisterCategoryRoutes(router, db, logger)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped cleanly.")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
