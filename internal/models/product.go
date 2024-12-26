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
