package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/middlewares"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserRoutes(router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteUsers = "/users"

	// Initialize the repository
	userRepository := repositories.NewUserRepository(db)

	// Initialize the service
	userService := services.NewUserService(userRepository)

	// Initialize the handler
	userHandler := handlers.NewUserHandler(userService, logger, errorLogger)

	router.Use(middlewares.ErrorMiddleware) // Attach ErrorMiddleware

	router.HandleFunc(RouteUsers, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userHandler.ReadAll(w, r)
		case "POST":
			userHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteUsers+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "GET":
			userHandler.Read(w, r, id)
		case "PUT":
			userHandler.Update(w, r, id)
		case "DELETE":
			userHandler.Delete(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
