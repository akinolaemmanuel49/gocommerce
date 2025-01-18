package services

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

func NewAuthService(userProvider services.UserService, config *configs.Config) *AuthService {
	return &AuthService{userProvider: userProvider, config: config}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	// Find user by email
	user, err := s.userProvider.RetrieveUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Verify password
	if utils.VerifyPassword(password, user.PasswordHash) {
		return nil, errors.NewUnauthorizedError()
	}

	// Generate token
	panic("")
}
