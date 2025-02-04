package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	auth_middleware "github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/middlewares"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterOrderRoutes initializes repositories, services and attaches handlers to the router
func RegisterOrderRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteOrders = "/orders"
	jwtSecretKey := []byte(config.JWTSecretKey)

	// Initialize repositories
	orderRepository := repositories.NewOrderRepository(db)
	userRepository := repositories.NewUserRepository(db)
	productRepository := repositories.NewProductRepository(db)
	cartRepository := repositories.NewCartRepository(db)

	// Initialize the publisher
	publisher, err := queue.NewPublisher(config)
	if err != nil {
		errorLogger.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(userRepository)
	productService := services.NewProductService(productRepository)
	cartService := services.NewCartService(cartRepository, *userService, *productService)
	orderService := services.NewOrderService(orderRepository, publisher, userService, cartService)

	// Initialize the handler
	orderHandler := handlers.NewOrderHandler(orderService, logger, errorLogger)

	router.Use(middlewares.ErrorMiddleware) // Attach ErrorMiddleware
	// router.Use(corsMiddleware.Handler)
	authMiddleware := auth_middleware.AuthMiddleware(jwtSecretKey)

	router.Handle(RouteOrders, authMiddleware(http.HandlerFunc(orderHandler.Create))).Methods("POST")
	router.Handle(RouteOrders, authMiddleware(http.HandlerFunc(orderHandler.Read))).Methods("GET")
	router.Handle(RouteOrders+"/all", authMiddleware(http.HandlerFunc(orderHandler.ReadAll))).Methods("GET")
	router.Handle(RouteOrders+"{id}/status", http.HandlerFunc(orderHandler.UpdateOrderStatus)).Methods("PUT")
	router.Handle(RouteOrders+"{id}/address", http.HandlerFunc(orderHandler.UpdateOrderShippingAddress)).Methods("PUT")
	router.Handle(RouteOrders+"{id}/carts/add/{cartID}", http.HandlerFunc(orderHandler.AddCartToOrder)).Methods("PUT")
	router.Handle(RouteOrders+"{id}/items/remove/{cartID}", http.HandlerFunc(orderHandler.RemoveCartFromOrder)).Methods("PUT")
	router.Handle(RouteOrders+"{id}/confirm", http.HandlerFunc(orderHandler.ConfirmOrder)).Methods("PUT")
	router.Handle(RouteOrders+"{id}/cancel", http.HandlerFunc(orderHandler.CancelOrder)).Methods("PUT")
	router.Handle(RouteOrders+"{id}", http.HandlerFunc(orderHandler.Delete)).Methods("DELETE")
}
