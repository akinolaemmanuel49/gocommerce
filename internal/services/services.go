package services

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
	logger         *log.Logger
}

type ProductService struct {
	productRepository *repositories.ProductRepository
	logger            *log.Logger
}

type OrderService struct {
	orderRepository *repositories.OrderRepository
	logger          *log.Logger
}

type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
	logger             *log.Logger
}
