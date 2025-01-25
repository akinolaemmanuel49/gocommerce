package services

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

// UserService represents business logic for user services
type UserService struct {
	userRepository *repositories.UserRepository
}

// CategoryService represents business logic for category services
type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
}

// ProductService represents business logic for product services
type ProductService struct {
	CategoryService   // Embedding CategoryService to allow ProductService to access category-related functionalities
	productRepository *repositories.ProductRepository
}

// CartService represents business logic for cart services
type CartService struct {
	UserService    // Embedding UserService to allow CartService to access user-related functionalities
	ProductService // Embedding ProductService to allow CartService to access product-related functionalities
	cartRepository *repositories.CartRepository
}

// OrderService represents business logic for order services
type OrderService struct {
	UserService     // Embedding UserService to allow OrderService to access user-related functionalities
	CartService     // Embedding CartService to allow OrderService to access cart-related functionalities
	orderRepository *repositories.OrderRepository
	publisher       *queue.Publisher
	errorLogger     *log.Logger
}
