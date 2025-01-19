package services

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
}
type ProductService struct {
	CategoryService   // Embedding CategoryService to allow ProductService to access category-related functionalities
	productRepository *repositories.ProductRepository
}
type CartService struct {
	UserService    // Embedding UserService to allow CartService to access user-related functionalities
	ProductService // Embedding ProductService to allow CartService to access product-related functionalities
	cartRepository *repositories.CartRepository
}
type OrderService struct {
	UserService     // Embedding UserService to allow OrderService to access user-related functionalities
	CartService     // Embedding CartService to allow OrderService to access cart-related functionalities
	orderRepository *repositories.OrderRepository
	publisher       *queue.Publisher
	errorLogger     *log.Logger
}
