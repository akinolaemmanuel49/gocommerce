package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

const (
	HashCost = 10
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

// CreateUser creates a new instance of a user and commits it to the database
func (s *UserService) CreateUser(ctx context.Context, newUser *models.CreateUser) (*models.User, error) {
	// Hash password
	passwordInByte := []byte(newUser.Password)
	passwordHashInByte, err := bcrypt.GenerateFromPassword(passwordInByte, HashCost)
	if err != nil {
		return nil, fmt.Errorf("error creating new user: %w", err)
	}

	passwordHashString := string(passwordHashInByte)

	// Transform CreateUser to User
	user := &models.User{
		Email:        newUser.Email,
		PasswordHash: passwordHashString,
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
		Role:         newUser.Role,
		CommonFields: models.CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert user into the database
	result, err := s.userRepository.Insert(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating new user: %w", err)
	}

	// Convert result.InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	user.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return user, nil
}
