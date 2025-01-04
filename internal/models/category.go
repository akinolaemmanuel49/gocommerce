package models

import "time"

type Category struct {
	ID          string  `bson:"_id,omitempty"`         // Unique identifier for the category
	Name        string  `bson:"name,omitempty"`        // Name of the category
	Description string  `bson:"description,omitempty"` // Description of the category
	ParentID    *string `bson:"parentId,omitempty"`    // Optional parent category ID for hierarchical categorization
	Image       string  `bson:"image,omitempty"`       // URL or path to an image representing the category
	// Products     []string `bson:"products,omitempty"`    // List of associated product IDs
	CommonFields `bson:"inline"`
}

type CreateCategory struct {
	Name        string  `json:"name,omitempty" validate:"required"` // Name is required
	Description string  `json:"description,omitempty"`              // Optional description
	ParentID    *string `json:"parentId,omitempty"`                 // Optional parent category
	Image       string  `json:"image,omitempty"`                    // Optional category image
}

type UpdateCategory struct {
	Name        *string `json:"name,omitempty" validate:"required"` // Name is required
	Description *string `json:"description,omitempty"`              // Optional description
	Image       *string `json:"image,omitempty"`
}

func NewCategory(newCategory *CreateCategory) *Category {
	return &Category{
		Name:        newCategory.Name,
		Description: newCategory.Description,
		ParentID:    newCategory.ParentID,
		Image:       newCategory.Image,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}
}
