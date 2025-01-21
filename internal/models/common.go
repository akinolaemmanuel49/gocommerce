package models

import "time"

// CommonFields adds common metadata to all models.
type CommonFields struct {
	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	IsDeleted bool      `bson:"isDeleted" json:"isDeleted"`
}

// Address model
type Address struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	Country string `json:"country" validate:"required"`
}

type UpdateAddress struct {
	Street  *string `json:"street" validate:"required"`
	City    *string `json:"city" validate:"required"`
	State   *string `json:"state" validate:"required"`
	Zip     *string `json:"zip" validate:"required"`
	Country *string `json:"country" validate:"required"`
}

type MultipleEntityClientResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"nextCursor"`
}

type ClientResponse struct {
	Message string `json:"message"`
}
