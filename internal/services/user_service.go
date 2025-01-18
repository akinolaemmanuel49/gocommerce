package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	// Ensure required fields are present and valid
	if newUser.Email == "" {
		return nil, errors.NewValidationError("User", "Email is required")
	}

	if !utils.ValidateEmail(newUser.Email) {
		return nil, errors.NewValidationError("User", "Email is invalid")
	}

	if newUser.Password == "" {
		return nil, errors.NewValidationError("User", "Password is required")
	}

	if len(newUser.Password) < 8 {
		return nil, errors.NewValidationError("User", "Password must be at least 8 characters long")
	}

	if newUser.FirstName == "" {
		return nil, errors.NewValidationError("User", "First name is required")
	}

	if newUser.LastName == "" {
		return nil, errors.NewValidationError("User", "Last name is required")
	}

	if newUser.Role == "" {
		return nil, errors.NewValidationError("User", "Role is required")
	}

	if newUser.Role != "admin" && newUser.Role != "customer" {
		return nil, errors.NewValidationError("User", "Role must be either 'admin' or 'customer'")
	}

	// Transform CreateUser to User
	user, err := models.NewUser(newUser, HashCost)
	if err != nil {
		return nil, err
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

// VerifyUser accepts plain password and hashed password, compares them then, returns true if it is a match, false otherwise
func (s *UserService) VerifyUser(ctx context.Context, plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
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
	user := models.UserUpdate(updatedUser, existingUser)

	update := bson.M{
		"$set": user,
	}

	// Update user in database
	_, err = s.userRepository.Update(ctx, ID, update)
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

	if existingUser != nil {
		// Apply transformation, set user IsDeleted field to true
		user := &models.User{
			CommonFields: models.CommonFields{
				IsDeleted: true,
				UpdatedAt: time.Now(),
			},
		}

		deleted := bson.M{"$set": user}

		_, err = s.userRepository.Update(ctx, ID, deleted)
		if err != nil {
			return err
		}
	}

	return nil
}
