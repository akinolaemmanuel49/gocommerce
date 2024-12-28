package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

type HandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

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

type CategoryHandler struct {
	categoryService *services.CategoryService
	logger          *log.Logger
}
