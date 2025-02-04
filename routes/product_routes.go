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

// RegisterProductRoutes initializes repositories, services and attaches handlers to the router
func RegisterProductRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteProducts = "/products"
	jwtSecretKey := []byte(config.JWTSecretKey)

	// Initialize the repository
	productRepository := repositories.NewProductRepository(db)

	// Initialize the service
	productService := services.NewProductService(productRepository)

	// Initialize the handler
	productHandler := handlers.NewProductHandler(productService, logger, errorLogger)

	router.Use(middlewares.ErrorMiddleware) // Attach ErrorMiddleware
	// router.Use(corsMiddleware.Handler)
	authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey)

	// Attach routes
	router.Handle(RouteProducts, authMiddleware(http.HandlerFunc(productHandler.Create))).Methods("POST")
	router.Handle(RouteProducts+"/all", http.HandlerFunc(productHandler.ReadAll)).Methods("GET")
	router.Handle(RouteProducts+"/{id}", http.HandlerFunc(productHandler.Read)).Methods("GET")
	router.Handle(RouteProducts+"/{id}", authMiddleware(http.HandlerFunc(productHandler.Update))).Methods("PUT")
	router.Handle(RouteProducts+"/{id}", authMiddleware(http.HandlerFunc(productHandler.Delete))).Methods("DELETE")
}
