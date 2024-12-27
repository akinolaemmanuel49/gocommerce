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

func RegisterUserRoutes(router *mux.Router, db *mongo.Database, logger *log.Logger) {
	const RouteUsers = "/users"

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, logger)
	userHandler := handlers.NewUserHandler(userService, logger)

	router.HandleFunc(RouteUsers, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userHandler.GetAllUsers(w, r)
		case "POST":
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteUsers+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "GET":
			userHandler.ReadUser(w, r, id)
		case "PATCH":
			userHandler.UpdateUser(w, r, id)
		case "DELETE":
			// TODO
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
