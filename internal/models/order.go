package models

type OrderItem struct {
	ProductID string  `bson:"productId,omitempty"`
	Quantity  int     `bson:"quantity,omitempty"`
	Price     float64 `bson:"price,omitempty"`
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
