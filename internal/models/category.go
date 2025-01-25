package models

import "time"

// Category model with bson tags
type Category struct {
	ID           string  `bson:"_id,omitempty"`         // Unique identifier for the category
	Name         string  `bson:"name,omitempty"`        // Name of the category
	Description  string  `bson:"description,omitempty"` // Description of the category
	ParentID     *string `bson:"parentId,omitempty"`    // Optional parent category ID for hierarchical categorization
	Image        string  `bson:"image,omitempty"`       // URL or path to an image representing the category
	CommonFields `bson:"inline"`
}

// CreateCategory provides an interface to create a new category
type CreateCategory struct {
	Name        string  `json:"name,omitempty" validate:"required"` // Name is required
	Description string  `json:"description,omitempty"`              // Optional description
	ParentID    *string `json:"parentId,omitempty"`                 // Optional parent category
	Image       string  `json:"image,omitempty"`                    // Optional category image
}

// UpdateCategory provides an interface to update a category
type UpdateCategory struct {
	Name        *string `json:"name,omitempty" validate:"required"` // Name is required
	Description *string `json:"description,omitempty"`              // Optional description
	Image       *string `json:"image,omitempty"`
}

// NewCategory returns a pointer to the new category
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

// CategoryUpdate returns a pointer to the updated category
func CategoryUpdate(updatedCategory *UpdateCategory, existingCategory *Category) *Category {
	return &Category{
		Name:        IfNotNil(updatedCategory.Name, existingCategory.Name),
		Description: IfNotNil(updatedCategory.Description, existingCategory.Description),
		Image:       IfNotNil(updatedCategory.Image, existingCategory.Image),
		CommonFields: CommonFields{
			UpdatedAt: time.Now(),
		},
	}
}
