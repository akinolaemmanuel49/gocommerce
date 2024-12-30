package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/utils"
)

const (
	HashCost = 10
)

// NewUserService creates a new instance of UserService
func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// CreateUser creates a new instance of a user and commits it to the database
func (s *UserService) CreateUser(ctx context.Context, newUser *models.CreateUser) (*models.User, error) {
	// Check for existing user
	existingUser, err := s.userRepository.FindByEmail(ctx, newUser.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.NewConflictError("User", "email", newUser.Email)
	}

	// Hash password
	passwordInByte := []byte(newUser.Password)
	passwordHashInByte, err := bcrypt.GenerateFromPassword(passwordInByte, HashCost)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Convert result.InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	user.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return user, nil
}

// RetrieveUserByID retrieves a user by ID
func (s *UserService) RetrieveUserByID(ctx context.Context, ID string) (*models.User, error) {
	user, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// RetrieveUserByEmail retrieves a user by email
func (s *UserService) RetrieveUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// RetrieveAllUsers retrieves paginated users with optional filters
func (s *UserService) RetrieveAllUsers(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.User, string, error) {
	users, nextCursor, err := s.userRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", err
	}

	return users, nextCursor, nil
}

// UpdateUserByID updates an instance of a user and commits it to the database
func (s *UserService) UpdateUserByID(ctx context.Context, ID string, updatedUser *models.UpdateUser) (*models.User, error) {
	// Check for existing user
	existingUser, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	// Transform UpdateUser to User
	user := &models.User{
		ID:        ID,
		FirstName: utils.IfNotNil(updatedUser.FirstName, existingUser.FirstName),
		LastName:  utils.IfNotNil(updatedUser.LastName, existingUser.LastName),
		Phone:     utils.IfNotNil(updatedUser.Phone, existingUser.Phone),
		Address:   utils.MergeAddress(updatedUser.Address, existingUser.Address),
		CommonFields: models.CommonFields{
			UpdatedAt: time.Now(),
		},
	}

	// Update user in database
	_, err = s.userRepository.Update(ctx, ID, user)
	if err != nil {
		return nil, err
	}

	// Return the updated user
	updatedRecord, err := s.userRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return updatedRecord, nil
}

// DeleteUserByID sets the IsDeleted flag for a user instance to true (performs a soft-delete)
func (s *UserService) DeleteUserByID(ctx context.Context, ID string) error {
	// Check for existing user
	existingUser, err := s.userRepository.FindByID(ctx, ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}

	// Set existingUser IsDeleted field to true
	existingUser.IsDeleted = true
	existingUser.UpdatedAt = time.Now()
	_, err = s.userRepository.Update(ctx, ID, existingUser)
	if err != nil {
		return err
	}

	return nil
}
