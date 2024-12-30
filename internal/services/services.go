package services

import (
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
}

type CategoryService struct {
	categoryRepository *repositories.CategoryRepository
}
