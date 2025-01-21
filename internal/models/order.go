package models

import "time"

type OrderStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=pending shipped delivered"`
}

type OrderShippingAddressUpdate struct {
	ShippingAddress Address `json:"shippingAddress,omitempty"`
}

type Order struct {
	ID              string  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          string  `bson:"userId,omitempty" json:"userId,omitempty"`
	Carts           []Cart  `bson:"carts,omitempty" json:"carts,omitempty"`
	TotalPrice      float64 `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	Status          string  `bson:"status,omitempty" json:"status,omitempty"` // "pending", "shipped", "delivered"
	ShippingAddress Address `bson:"shippingAddress,omitempty" json:"shippingAddress,omitempty"`
	IsLocked        bool    `bson:"isLocked,omitempty" json:"isLocked,omitempty"`
	IsCancelled     bool    `bson:"isCancelled, omitempty" json:"isCancelled,omitempty"`
	CommonFields    `bson:"inline"`
}

type CreateOrder struct {
	Carts           []Cart  `json:"carts" validate:"required,dive"`
	TotalPrice      float64 `json:"totalPrice" validate:"required,gt=0"`
	Status          string  `json:"status" validate:"required,oneof=pending shipped delivered"`
	ShippingAddress Address `json:"shippingAddress" validate:"required"`
}

type UpdateOrder struct {
	Status          string  `json:"status,omitempty" validate:"omitempty,oneof=pending shipped delivered"`
	ShippingAddress Address `json:"shippingAddress,omitempty"`
}

func NewOrder(newOrder *CreateOrder, user *User) *Order {
	return &Order{
		UserID:          user.ID,
		Carts:           newOrder.Carts,
		TotalPrice:      newOrder.TotalPrice,
		Status:          "pending",
		ShippingAddress: user.Address,
		IsLocked:        false,
		IsCancelled:     false,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
	}
}
