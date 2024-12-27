package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterCategoryRoutes(router *http.ServeMux, db *mongo.Database, logger *log.Logger) {
	const RouteCategories = "/categories"

	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
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
