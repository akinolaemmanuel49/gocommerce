package services

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

type ProductService struct {
	productRepository *repositories.ProductRepository
}

type OrderService struct {
	orderRepository *repositories.OrderRepository
	userService     UserService
	cartService     CartService
	publisher       *queue.Publisher
	errorLogger     *log.Logger
}

type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
}

type CartService struct {
	cartRepository *repositories.CartRepository
	userService    UserService
	productService ProductService
}
