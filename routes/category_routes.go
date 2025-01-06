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

func RegisterCategoryRoutes(router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteCategories = "/categories"

	// Initialize the repository
	categoryRepository := repositories.NewCategoryRepository(db)

	// Initialize the service
	categoryService := services.NewCategoryService(categoryRepository)

	// Initialize the handler
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger, errorLogger)

	router.HandleFunc(RouteCategories, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			categoryHandler.ReadAll(w, r)
		case "POST":
			categoryHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteCategories+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "GET":
			categoryHandler.Read(w, r, id)
		case "PATCH":
			categoryHandler.Update(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteCategories+"/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			categoryHandler.Delete(w, r, id)
		}
	})
}
