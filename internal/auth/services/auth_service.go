package services

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewAuthService creates a new instance of AuthService
func NewAuthService(userProvider services.UserService, config *configs.Config) *AuthService {
	return &AuthService{userProvider: userProvider, config: config}
}

// Login method for AuthService implements business logic for the login process
func (s *AuthService) Login(ctx context.Context, email string, password string) (*models.Token, error) {
	// Find user by email
	user, err := s.userProvider.RetrieveUserByEmail(ctx, email)
	if err == mongo.ErrNoDocuments {
		return nil, errors.NewNotFoundError("User", "Email", email)
	}
	if err != nil {
		return nil, err
	}

	// Verify password
	if !utils.VerifyPassword(password, user.PasswordHash) {
		return nil, errors.NewAuthorizationError("Invalid password")
	}

	// Generate token
	stringToken, err := utils.GenerateJWT([]byte(s.config.JWTSecretKey), user.ID, user.Role)
	if err != nil {
		return nil, errors.NewAuthorizationError("Failed to generate token")
	}

	token := &models.Token{Token: stringToken}

	return token, nil
}
