package models

type Category struct {
	ID           string   `bson:"_id,omitempty"`         // Unique identifier for the category
	Name         string   `bson:"name,omitempty"`        // Name of the category
	Description  string   `bson:"description,omitempty"` // Description of the category
	ParentID     *string  `bson:"parentId,omitempty"`    // Optional parent category ID for hierarchical categorization
	Image        string   `bson:"image,omitempty"`       // URL or path to an image representing the category
	Products     []string `bson:"products,omitempty"`    // List of associated product IDs
	CreatedAt    int64    `bson:"createdAt,omitempty"`   // Timestamp of creation
	UpdatedAt    int64    `bson:"updatedAt,omitempty"`   // Timestamp of last update
	CommonFields `bson:"inline"`
}

type CreateCategory struct {
	Name        string  `bson:"name,omitempty" validate:"required"` // Name is required
	Description string  `bson:"description,omitempty"`              // Optional description
	ParentID    *string `bson:"parentId,omitempty"`                 // Optional parent category
	Image       string  `bson:"image,omitempty"`                    // Optional category image
}
