package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

const (
	HashCost = 10
)

// NewUserService creates a new instance of UserService
func NewUserService(userRepository *repositories.UserRepository, logger *log.Logger) *UserService {
	return &UserService{userRepository: userRepository, logger: logger}
}

// GetAllUsers retrieves paginated users with optional filters
func (s *UserService) GetAllUsers(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.User, string, error) {
	users, nextCursor, err := s.userRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching users: %w", err)
	}

	s.logger.Printf("Reading users")
	return users, nextCursor, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	s.logger.Printf("Reading user ID %s", user.ID)
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	s.logger.Printf("Reading user ID %s", ID)
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

	s.logger.Printf("Creating user ID %s", user.ID)
	return user, nil
}

// UpdateUser updates an instance of a user and commits it to the database
func (s *UserService) UpdateUser(ctx context.Context, ID string, updatedUser *models.UpdateUser) (*models.User, error) {
	existingUser, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching existing user: %w", err)
	}

	// Transform UpdateUser to User
	user := &models.User{
		ID:        ID,
		FirstName: ifNotEmpty(updatedUser.FirstName, existingUser.FirstName),
		LastName:  ifNotEmpty(updatedUser.LastName, existingUser.LastName),
		Phone:     ifNotEmpty(updatedUser.Phone, existingUser.Phone),
		Address:   mergeAddress(updatedUser.Address, existingUser.Address),
		CommonFields: models.CommonFields{
			UpdatedAt: time.Now(),
		},
	}

	// Update user in database
	_, err = s.userRepository.Update(ctx, ID, user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	// Return the updated user
	updatedRecord, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated user")
	}

	s.logger.Printf("Updating user ID %s", ID)
	return updatedRecord, nil
}
