package models

import (
	"time"

	"github.com/akinolaemmanuel49/gocommerce/internal/auth/utils"
)

// Database model with bson tags
type User struct {
	ID           string  `bson:"_id,omitempty" json:"id,omitempty"`
	Email        string  `bson:"email,omitempty" json:"email,omitempty"`
	PasswordHash string  `bson:"passwordHash,omitempty" json:"passwordHash,omitempty" swaggerignore:"true"`
	FirstName    string  `bson:"firstName,omitempty" json:"firstName,omitempty"`
	LastName     string  `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Address      Address `bson:"address,omitempty" json:"address,omitempty"`
	Phone        string  `bson:"phone,omitempty" json:"phone,omitempty"`
	Role         string  `bson:"role,omitempty" json:"role,omitempty"`
	CommonFields `bson:"inline"`
}

func ResponseUser(user *User) (*User, error) {
	responseUser := &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Address:   user.Address,
		Phone:     user.Phone,
		Role:      user.Role,
		CommonFields: CommonFields{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsDeleted: false,
		},
	}

	return responseUser, nil
}

// Request DTO for creating a user
type CreateUser struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=customer admin"`
}

// Request DTO for updating a user
type UpdateUser struct {
	FirstName *string        `json:"firstName,omitempty"`
	LastName  *string        `json:"lastName,omitempty"`
	Phone     *string        `json:"phone,omitempty"`
	Address   *UpdateAddress `json:"address,omitempty"`
}

func NewUser(newUser *CreateUser, HashCost int) (*User, error) {
	// Hash the password
	passwordHashString, err := utils.HashPassword(newUser.Password, HashCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        newUser.Email,
		PasswordHash: passwordHashString,
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
		Role:         newUser.Role,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}

	return user, nil
}

func UserUpdate(updatedUser *UpdateUser, existingUser *User) *User {
	return &User{
		FirstName: IfNotNil(updatedUser.FirstName, existingUser.FirstName),
		LastName:  IfNotNil(updatedUser.LastName, existingUser.LastName),
		Phone:     IfNotNil(updatedUser.Phone, existingUser.Phone),
		Address:   MergeAddress(updatedUser.Address, existingUser.Address),
		CommonFields: CommonFields{
			UpdatedAt: time.Now(),
		},
	}
}
