package models

type Address struct {
	Street  string `bson:"street,omitempty"`
	City    string `bson:"city,omitempty"`
	State   string `bson:"state,omitempty"`
	Zip     string `bson:"zip,omitempty"`
	Country string `bson:"country,omitempty"`
}

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
