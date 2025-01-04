package models

import (
	"time"
)

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

type CreateProduct struct {
	Name        string   `json:"name,omitempty" validate:"required"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price,omitempty" validate:"required,gt=0"`
	Images      []string `json:"images,omitempty"`
	CategoryID  string   `json:"categoryId,omitempty" validate:"required"`
	Brand       string   `json:"brand,omitempty"`
}

type UpdateProduct struct {
	Name        *string  `json:"name,omitempty" validate:"required"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"required,gt=0"`
	Brand       *string  `json:"brand,omitempty"`
}

func NewProduct(newProduct *CreateProduct) *Product {
	return &Product{
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Images:      newProduct.Images,
		CategoryID:  newProduct.CategoryID,
		Brand:       newProduct.Brand,
		// Rating:       3.0,
		// ReviewsCount: 0,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}
}

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
