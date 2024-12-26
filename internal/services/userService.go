package services

import (
	"context"
	"fmt"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// GetAllUsers retrieves paginated users with optional filters
func (s *UserService) GetAllUsers(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.User, string, error) {
	users, nextCursor, err := s.userRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching users: %w", err)
	}

	return users, nextCursor, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}
