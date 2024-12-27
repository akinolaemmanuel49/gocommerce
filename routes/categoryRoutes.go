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

func RegisterCategoryRoutes(router *mux.Router, db *mongo.Database, logger *log.Logger) {
	const RouteCategories = "/categories"

	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)

	router.HandleFunc(RouteCategories, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			categoryHandler.GetAllCategories(w, r)
		case "POST":
			categoryHandler.CreateCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
