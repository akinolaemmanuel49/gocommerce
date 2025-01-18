package services

import (
	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

type AuthService struct {
	userProvider services.UserService
	config       *configs.Config
}

type JWTService struct{}
