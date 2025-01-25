package models

import (
	"time"
)

// Product model with bson tags
type Product struct {
	ID          string   `bson:"_id,omitempty"`
	Name        string   `bson:"name,omitempty"`
	Description string   `bson:"description,omitempty"`
	Price       float64  `bson:"price,omitempty"`
	Images      []string `bson:"images,omitempty"`
	CategoryID  string   `bson:"categoryId,omitempty"`
	Brand       string   `bson:"brand,omitempty"`
	// Rating       float64  `bson:"rating,omitempty"`
	// ReviewsCount int      `bson:"reviewsCount,omitempty"`
	CommonFields `bson:"inline"`
}

// CreateProduct provides an interface to create new products
type CreateProduct struct {
	Name        string   `json:"name,omitempty" validate:"required"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price,omitempty" validate:"required,gt=0"`
	Images      []string `json:"images,omitempty"`
	CategoryID  string   `json:"categoryId,omitempty" validate:"required"`
	Brand       string   `json:"brand,omitempty"`
}

// UpdateProduct provides an interface to update products
type UpdateProduct struct {
	Name        *string  `json:"name,omitempty" validate:"required"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"required,gt=0"`
	Brand       *string  `json:"brand,omitempty"`
}

// NewProduct returns a pointer to the new product
func NewProduct(newProduct *CreateProduct) *Product {
	return &Product{
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Images:      newProduct.Images,
		CategoryID:  newProduct.CategoryID,
		Brand:       newProduct.Brand,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}
}

// ProductUpdate returns a pointer to the updated product
func ProductUpdate(updatedProduct *UpdateProduct, existingProduct *Product) *Product {
	return &Product{
		Name:        IfNotNil(updatedProduct.Name, existingProduct.Name),
		Description: IfNotNil(updatedProduct.Description, existingProduct.Description),
		Price:       IfNotNil(updatedProduct.Price, existingProduct.Price),
		Brand:       IfNotNil(updatedProduct.Brand, existingProduct.Brand),
		CommonFields: CommonFields{
			UpdatedAt: time.Now(),
		},
	}
}
