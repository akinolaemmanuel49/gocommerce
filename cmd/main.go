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
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"

	"github.com/akinolaemmanuel49/gocommerce/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Load config
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Initialize MongoDB
	db, err := database.ConnectMongoDB(config.MongoDBURI, config.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = db.Client().Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to ensure a successful connection
	if err := db.Client().Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB successfully!")

	// Initialize RabbitMQ
	conn, ch, err := queue.ConnectRabbitMQ(&config)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}()
	defer func() {
		if err = ch.Close(); err != nil {
			log.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}()

	// Start consuming messages from RabbitMQ
	go queue.ConsumeOrderNotifications(&config, ch)

	// Setup HTTP routes
	router := http.NewServeMux()
	setupRoutes(router, db)

	// Start the HTTP server with graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		log.Println("Server is running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

// setupRoutes configures the HTTP routes
func setupRoutes(router *http.ServeMux, db *mongo.Database) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	// Routes: Register them with the passed router
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/users", userHandler.GetAllUsers)
	// http.HandleFunc("/users/email", userHandler.GetUserByEmail) // Example: /users/email?email=example@example.com
	router.HandleFunc("/products", productHandler.GetAllProducts)
	// http.HandleFunc("/products/id", productHandler.GetProductByID) // Example: /products/id?id=12345
}

// gracefulShutdown handles cleanup and ensures the server stops gracefully
func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped cleanly.")
}

// healthHandler is a simple handler to check if the server is running.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
