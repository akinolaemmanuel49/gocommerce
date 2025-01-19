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

func RegisterCategoryRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteCategories = "/categories"
	jwtSecretKey := []byte(config.JWTSecretKey)

	// Initialize the repository
	categoryRepository := repositories.NewCategoryRepository(db)

	// Initialize the service
	categoryService := services.NewCategoryService(categoryRepository)

	// Initialize the handler
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger, errorLogger)

	router.Use(middlewares.ErrorMiddleware)                         // Attach ErrorMiddleware
	authMiddleware := auth_middlewares.AuthMiddleware(jwtSecretKey) // Attach AuthMiddleware

	// Attach routes
	router.Handle(RouteCategories, authMiddleware(http.HandlerFunc(categoryHandler.Create))).Methods("POST")
	router.Handle(RouteCategories+"/all", http.HandlerFunc(categoryHandler.ReadAll)).Methods("GET")
	router.Handle(RouteCategories+"/{id}", http.HandlerFunc(categoryHandler.Read)).Methods("GET")
	router.Handle(RouteCategories+"/{id}", authMiddleware(http.HandlerFunc(categoryHandler.Update))).Methods("PUT")
	router.Handle(RouteCategories+"/{id}", authMiddleware(http.HandlerFunc(categoryHandler.Delete))).Methods("DELETE")
}
