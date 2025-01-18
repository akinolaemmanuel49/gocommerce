package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/auth/services"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *log.Logger
	errorLogger *log.Logger
}

type IAuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}
