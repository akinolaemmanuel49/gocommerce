package models

import "time"

type Cart struct {
	ID           string     `json:"id" bson:"_id"`
	UserID       string     `json:"userId" bson:"userId"`
	Items        []CartItem `json:"items" bson:"items"`
	TotalPrice   float64    `json:"totalPrice" bson:"totalPrice"`
	IsOrdered    bool       `json:"isOrdered" bson:"isOrdered"`
	Name         string     `json:"name,omitempty" bson:"name,omitempty"`
	CommonFields `bson:"inline"`
}

type CreateCart struct {
	UserID     string  `json:"userId" validate:"required"`
	TotalPrice float64 `json:"totalPrice" validate:"required,gt=0"`
	Name       string  `json:"name,omitempty"`
}

type CartItem struct {
	ID       string  `json:"id" bson:"_id"`
	CartID   string  `json:"cartId" bson:"cartId"`
	Product  Product `json:"product" bson:"product"`
	Quantity int     `json:"quantity" bson:"quantity"`
}

func NewCart(newCart *CreateCart) *Cart {
	return &Cart{
		UserID:     newCart.UserID,
		TotalPrice: 0,
		IsOrdered:  false,
		Name:       newCart.Name,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}
}
