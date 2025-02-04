package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/configs"
	auth_routes "github.com/akinolaemmanuel49/gocommerce/internal/auth/routes"
	l "github.com/akinolaemmanuel49/gocommerce/log"
	"github.com/akinolaemmanuel49/gocommerce/routes"
	"github.com/gorilla/mux"

	_ "github.com/akinolaemmanuel49/gocommerce/docs"
	"github.com/akinolaemmanuel49/gocommerce/internal/database"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @title GoCommerce API
// @version 1.0
// @description This is an e-commerce server GoCommerce server.
// @termsOfService https://gocommerce-h1v5.onrender.com/terms/

// @contact.name Akinola Abiodun E.
// @contact.email biteatertest@gmail.com

// @license.name MIT
// @license.url https://gocommerce-h1v5.onrender.com/license/

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization

// @host https://gocommerce-h1v5.onrender.com

// Config variable contains project config details
var Config configs.Config

func main() {
	// Setup flags
	envFile := flag.String("env-file", ".", "Path to the .env file")
	flag.Parse()

	// Setup logger
	logger, err := l.SetupLogger("service.log", "INFO")
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Fatalf("%s", fmt.Sprintf("%-7s: Error setting up error logger: %v", "ERROR", err))
	}
	errorLogger, err := l.SetupLogger("service.log", "ERROR")
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime)
		log.Fatalf("%s", fmt.Sprintf("%-7s: Error setting up error logger: %v", "ERROR", err))
	}

	// Load config
	config, err := configs.LoadConfig(*envFile, logger, errorLogger)
	if err != nil {
		errorLogger.Fatalf("Error reading config file: %v", err)
	}
	Config = config

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_db := database.NewDatabase()
	// Initialize MongoDB
	db, err := _db.Connect(config.MongoDBURI, config.MongoDBName)
	if err != nil {
		errorLogger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = db.Client().Disconnect(ctx); err != nil {
			errorLogger.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to ensure a successful connection
	if err := db.Client().Database("admin").RunCommand(ctx, bson.M{"ping": 1}).Err(); err != nil {
		errorLogger.Fatalf("Failed to ping MongoDB: %v", err)
	}
	logger.Println("Connected to MongoDB successfully")

	// Initialize RabbitMQ
	conn, ch, err := queue.ConnectRabbitMQ(&config, logger, errorLogger, true)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			errorLogger.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}()
	defer func() {
		if err = ch.Close(); err != nil {
			errorLogger.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}()

	// Start consuming messages from RabbitMQ
	go queue.ConsumeOrderNotifications(&config, ch, logger, errorLogger)

	// Setup HTTP routes
	router := mux.NewRouter()
	setupRoutes(&config, router, db, logger, errorLogger)

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
	gracefulShutdown(server, logger, errorLogger)
}

func setupRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	// Health check
	routes.RegisterHealthRoute(router, logger, errorLogger)

	// Legal routes
	routes.RegisterLegalRoutes(router)

	// User routes
	routes.RegisterUserRoutes(config, router, db, logger, errorLogger)

	// Category routes
	routes.RegisterCategoryRoutes(config, router, db, logger, errorLogger)

	// Product routes
	routes.RegisterProductRoutes(config, router, db, logger, errorLogger)

	// Cart Routes
	routes.RegisterCartRoutes(config, router, db, logger, errorLogger)

	// Order routes
	routes.RegisterOrderRoutes(config, router, db, logger, errorLogger)

	// Auth routes
	auth_routes.RegisterAuthRoutes(config, router, db, logger, errorLogger)

	// Swagger documentation route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Catch-all for unmatched routes
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errors.HandleError(w, r, errors.NewNotFoundError("route", "path", r.URL.Path), errorLogger)
	})

	// MethodNotAllowedHandler for unmatched HTTP methods
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errors.HandleError(w, r, errors.NewMethodNotAllowedError(r.Method), errorLogger)
	})

}

func gracefulShutdown(server *http.Server, logger, errorLogger *log.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		errorLogger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Println("Server stopped cleanly.")
}
