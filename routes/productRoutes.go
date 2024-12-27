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

func RegisterProductRoutes(router *mux.Router, db *mongo.Database, logger *log.Logger) {
	const RouteProducts = "/products"

	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository, logger)
	productHandler := handlers.NewProductHandler(productService, logger)

	router.HandleFunc(RouteProducts, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			productHandler.GetAllProducts(w, r)
		case "POST":
			productHandler.CreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
