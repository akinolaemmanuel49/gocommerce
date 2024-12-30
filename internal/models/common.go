package models

import "time"

// CommonFields adds common metadata to all models.
type CommonFields struct {
	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	IsDeleted bool      `bson:"isDeleted" json:"isDeleted"`
}
