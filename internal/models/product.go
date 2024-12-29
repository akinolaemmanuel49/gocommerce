package models

type Product struct {
	ID           string   `bson:"_id,omitempty"`
	Name         string   `bson:"name,omitempty"`
	Description  string   `bson:"description,omitempty"`
	Price        float64  `bson:"price,omitempty"`
	Images       []string `bson:"images,omitempty"`
	CategoryID   string   `bson:"categoryId,omitempty"`
	Brand        string   `bson:"brand,omitempty"`
	Rating       float64  `bson:"rating,omitempty"`
	ReviewsCount int      `bson:"reviewsCount,omitempty"`
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
