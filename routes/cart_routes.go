package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	auth_middleware "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/middlewares"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterCartRoutes initializes repositories, services and attaches handlers to the router
func RegisterCartRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteCarts = "/carts"
	jwtSecretKey := []byte(config.JWTSecretKey)

	// Initialize repositories
	cartRepository := repositories.NewCartRepository(db)
	userRepository := repositories.NewUserRepository(db)
	productRepository := repositories.NewProductRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepository)
	productService := services.NewProductService(productRepository)
	cartService := services.NewCartService(cartRepository, *userService, *productService)

	// Initialize the handler
	cartHandler := handlers.NewCartHandler(cartService, logger, errorLogger)

	router.Use(middlewares.ErrorMiddleware) // Attach ErrorMiddleware
	authMiddleware := auth_middleware.AuthMiddleware(jwtSecretKey)

	router.Handle(RouteCarts, authMiddleware(http.HandlerFunc(cartHandler.Create))).Methods("POST")
	router.Handle(RouteCarts, authMiddleware(http.HandlerFunc(cartHandler.Read))).Methods("GET")
	router.Handle(RouteCarts+"/all", authMiddleware(http.HandlerFunc(cartHandler.ReadAll))).Methods("GET")
	router.Handle(RouteCarts+"{id}/items/add", http.HandlerFunc(cartHandler.AddProductToCart)).Methods("PUT")
	router.Handle(RouteCarts+"{id}/items/remove", http.HandlerFunc(cartHandler.RemoveProductFromCart)).Methods("PUT")
	router.Handle(RouteCarts+"{id}", http.HandlerFunc(cartHandler.Delete)).Methods("DELETE")

}
