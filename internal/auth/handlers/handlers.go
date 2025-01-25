package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/auth/services"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *services.AuthService
	logger      *log.Logger
	errorLogger *log.Logger
}

// IAuthHandler defines the interface for authentication handlers.
type IAuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}
