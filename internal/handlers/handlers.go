package handlers

import (
	"log"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	logger      *log.Logger
}

type ProductHandler struct {
	productService *services.ProductService
	logger         *log.Logger
}

type OrderHandler struct {
	orderService *services.OrderService
	logger       *log.Logger
}
