package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	auth_middlewares "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/middlewares"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUserRoutes initializes repositories, services and attaches handlers to the router
func RegisterUserRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteUsers = "/users"
	jwtSecretKey := []byte(config.JWTSecretKey)

	// Initialize the repository
	userRepository := repositories.NewUserRepository(db)

	// Initialize the service
	userService := services.NewUserService(userRepository)

	// Initialize the handler
	userHandler := handlers.NewUserHandler(userService, logger, errorLogger)

	router.Use(middlewares.ErrorMiddleware) // Attach ErrorMiddleware
	router.Use(corsMiddleware.Handler)
	authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey) // Attach AuthMiddleware

	// Attach routes
	router.Handle(RouteUsers, http.HandlerFunc(userHandler.Create)).Methods("POST")
	router.Handle(RouteUsers, authMiddleware(http.HandlerFunc(userHandler.Read))).Methods("GET")
	router.Handle(RouteUsers+"/all", authMiddleware(http.HandlerFunc(userHandler.ReadAll))).Methods("GET")
	router.Handle(RouteUsers, authMiddleware(http.HandlerFunc(userHandler.Update))).Methods("PUT")
	router.Handle(RouteUsers, authMiddleware(http.HandlerFunc(userHandler.Delete))).Methods("DELETE")
}
