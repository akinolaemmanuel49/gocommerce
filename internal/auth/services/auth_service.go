package services

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuthService(userProvider services.UserService, config *configs.Config) *AuthService {
	return &AuthService{userProvider: userProvider, config: config}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	// Find user by email
	user, err := s.userProvider.RetrieveUserByEmail(ctx, email)
	if err == mongo.ErrNoDocuments {
		return "", errors.NewNotFoundError("User", "Email", email)
	}
	if err != nil {
		return "", err
	}

	// Verify password
	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", errors.NewUnauthorizedError()
	}

	// Generate token
	token, err := utils.GenerateJWT([]byte(s.config.JWTSecretKey), user.ID, user.Role)
	if err != nil {
		return "", errors.NewUnauthorizedError()
	}

	return token, nil
}
