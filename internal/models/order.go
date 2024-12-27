package models

type OrderItem struct {
	ProductID string  `json:"productId" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type Order struct {
	ID              string      `bson:"_id,omitempty"`
	UserID          string      `bson:"userId,omitempty"`
	Items           []OrderItem `bson:"items,omitempty"`
	TotalPrice      float64     `bson:"totalPrice,omitempty"`
	Status          string      `bson:"status,omitempty"` // "pending", "shipped", "delivered"
	ShippingAddress Address     `bson:"shippingAddress,omitempty"`
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
