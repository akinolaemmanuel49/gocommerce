package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/services"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	provider "github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAuthRoutes initializes repositories, services and attaches handlers to the router
func RegisterAuthRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	// Initialize AuthMiddleware with the JWT secret key
	const RouteAuth = "/auth"

	// Initialize the repositories
	userRepository := repositories.NewUserRepository(db)

	// Initialize the services
	userProvider := provider.NewUserService(userRepository)
	authService := services.NewAuthService(*userProvider, config)

	// Initialize the handler
	authHandler := handlers.NewAuthHandler(authService, logger, errorLogger)

	router.Handle(RouteAuth+"/login", http.HandlerFunc(authHandler.Login)).Methods("POST")
}
