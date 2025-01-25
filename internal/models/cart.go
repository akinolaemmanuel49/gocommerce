package models

import (
	"time"
)

// Cart model with bson and json tags
type Cart struct {
	ID           string     `json:"id" bson:"_id"`
	UserID       string     `json:"userId" bson:"userId"`
	Items        []CartItem `json:"items" bson:"items"`
	TotalPrice   float64    `json:"totalPrice" bson:"totalPrice"`
	Name         string     `json:"name,omitempty" bson:"name,omitempty"`
	CommonFields `bson:"inline"`
}

// CreateCart provides an interface to create a new cart
type CreateCart struct {
	Name string `json:"name,omitempty"`
}

// CartItem model with json tags
type CartItem struct {
	ID       string  `json:"id" bson:"_id"`
	CartID   string  `json:"cartId" bson:"cartId"`
	Product  Product `json:"product" bson:"product"`
	Quantity int     `json:"quantity" bson:"quantity"`
}

// CartItemCreate provides an interface to create a cart item
type CartItemCreate struct {
	ProductID string `json:"id"`
	Quantity  int    `json:"quantity"`
}

// CartItemUpdate provides an interface to update a cart item
type CartItemUpdate struct {
	ProductID string `json:"id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
}

// NewCart returns a pointer to the new cart
func NewCart(newCart *CreateCart, userID string) *Cart {
	return &Cart{
		UserID:     userID,
		TotalPrice: 0,
		Name:       newCart.Name,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}
}
