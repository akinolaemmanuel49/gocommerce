package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserRoutes(router *http.ServeMux, db *mongo.Database, logger *log.Logger) {
	const RouteUsers = "/users"

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
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
}
