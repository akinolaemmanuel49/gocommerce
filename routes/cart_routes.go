package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterCartRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteCarts = "/carts"

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

	router.HandleFunc(RouteCarts, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			cartHandler.ReadAll(w, r)
		case "POST":
			cartHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteCarts+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "GET":
			cartHandler.Read(w, r, id)
		case "DELETE":
			cartHandler.Delete(w, r, id)
		}
	})

	router.HandleFunc(RouteCarts+"/{id}/items/add", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PUT":
			cartHandler.AddProductToCart(w, r, id)
		}
	})

	router.HandleFunc(RouteCarts+"/{id}/items/remove/{productId}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]               // Extract the `id` path parameter
		productId := mux.Vars(r)["productId"] // Extract the `productId` path parameter

		switch r.Method {
		case "PUT":
			cartHandler.RemoveProductFromCart(w, r, id, productId)
		}
	})
}
