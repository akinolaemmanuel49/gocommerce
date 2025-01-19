package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/middlewares"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/services"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	provider "github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAuthRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	jwtSecretKey := []byte(config.JWTSecretKey)
	// Initialize AuthMiddleware with the JWT secret key
	authMiddleware := middlewares.AuthMiddleware(jwtSecretKey)
	const RouteAuth = "/auth"

	// Initialize the repositories
	userRepository := repositories.NewUserRepository(db)

	// Initialize the services
	userProvider := provider.NewUserService(userRepository)
	authService := services.NewAuthService(*userProvider, config)

	// Initialize the handler
	authHandler := handlers.NewAuthHandler(authService, logger, errorLogger)

	router.HandleFunc(RouteAuth+"/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			authHandler.Login(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Public routes
	router.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OKAY DESU!"))
	})

	// Protected routes
	router.Handle("/protected-endpoint", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaims(r.Context())
		if claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Write([]byte("Welcome, user ID: " + claims.UserID + ", Role: " + claims.Role))
	})))
}
