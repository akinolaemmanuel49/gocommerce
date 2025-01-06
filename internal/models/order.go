package models

import "time"

type OrderItem struct {
	ProductID string  `bson:"productId,omitempty"`
	Quantity  int     `bson:"quantity,omitempty"`
	Price     float64 `bson:"price,omitempty"`
}

type OrderItemCreate struct {
	ProductID string  `json:"productId" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type OrderStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=pending shipped delivered"`
}

type OrderShippingAddressUpdate struct {
	ShippingAddress Address `json:"shippingAddress,omitempty"`
}

type Order struct {
	ID              string      `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          string      `bson:"userId,omitempty" json:"userId,omitempty"`
	Items           []OrderItem `bson:"items,omitempty" json:"orderItem,omitempty"`
	TotalPrice      float64     `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	Status          string      `bson:"status,omitempty" json:"status,omitempty"` // "pending", "shipped", "delivered"
	ShippingAddress Address     `bson:"shippingAddress,omitempty" json:"shippingAddress,omitempty"`
	IsLocked        bool        `bson:"isLocked,omitempty" json:"isLocked,omitempty"`
	IsCancelled     bool        `bson:"isCancelled, omitempty" json:"isCancelled,omitempty"`
	CommonFields    `bson:"inline"`
}

type CreateOrder struct {
	UserID          string      `json:"userId" validate:"required"`
	Items           []OrderItem `json:"items" validate:"required,dive"`
	TotalPrice      float64     `json:"totalPrice" validate:"required,gt=0"`
	Status          string      `json:"status" validate:"required,oneof=pending shipped delivered"`
	ShippingAddress Address     `json:"shippingAddress" validate:"required"`
}

type UpdateOrder struct {
	Status          string  `json:"status,omitempty" validate:"omitempty,oneof=pending shipped delivered"`
	ShippingAddress Address `json:"shippingAddress,omitempty"`
}

func NewOrder(newOrder *CreateOrder) *Order {
	return &Order{
		UserID:          newOrder.UserID,
		Items:           newOrder.Items,
		TotalPrice:      newOrder.TotalPrice,
		Status:          "pending",
		ShippingAddress: newOrder.ShippingAddress,
		IsLocked:        false,
		IsCancelled:     false,
		CommonFields: CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: true,
		},
	}
}
