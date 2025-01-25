package services

import (
	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

// AuthService handles business logic for authentication
type AuthService struct {
	userProvider services.UserService
	config       *configs.Config
}
