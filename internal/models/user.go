package models

// Address model
type Address struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	Country string `json:"country" validate:"required"`
}

// Database model with bson tags
type User struct {
	ID           string  `bson:"_id,omitempty"`
	Email        string  `bson:"email,omitempty"`
	PasswordHash string  `bson:"passwordHash,omitempty"`
	FirstName    string  `bson:"firstName,omitempty"`
	LastName     string  `bson:"lastName,omitempty"`
	Address      Address `bson:"address,omitempty"`
	Phone        string  `bson:"phone,omitempty"`
	Role         string  `bson:"role,omitempty"` // "customer" or "admin"
	CommonFields `bson:"inline"`
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
	FirstName string  `json:"firstName,omitempty"`
	LastName  string  `json:"lastName,omitempty"`
	Phone     string  `json:"phone,omitempty"`
	Address   Address `json:"address,omitempty"`
}
