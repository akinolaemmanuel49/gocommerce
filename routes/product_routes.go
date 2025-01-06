package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterProductRoutes(router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteProducts = "/products"

	// Initialize the repository
	productRepository := repositories.NewProductRepository(db)

	// Initialize the service
	productService := services.NewProductService(productRepository)

	// Initialize the handler
	productHandler := handlers.NewProductHandler(productService, logger, errorLogger)

	router.HandleFunc(RouteProducts, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			productHandler.ReadAll(w, r)
		case "POST":
			productHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteProducts+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "GET":
			productHandler.Read(w, r, id)
		case "PATCH":
			productHandler.Update(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteProducts+"/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			productHandler.Delete(w, r, id)
		}
	})
}
