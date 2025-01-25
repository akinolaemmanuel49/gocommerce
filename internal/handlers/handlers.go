package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

// IUserHandler defines the interface for user handlers
type IUserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// ICategoryHandler defines the interface for category handlers
type ICategoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// IProductHandler defines the interface for product handlers
type IProductHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// ICartHandler defines the interface for cart handlers
type ICartHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	AddProductToCart(w http.ResponseWriter, r *http.Request)
	RemoveProductFromCart(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// IOrderHandler defines the interface for order handlers
type IOrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request)
	UpdateOrderShippingAddress(w http.ResponseWriter, r *http.Request)
	AddCartToOrder(w http.ResponseWriter, r *http.Request)
	RemoveCartFromOrder(w http.ResponseWriter, r *http.Request)
	ConfirmOrder(w http.ResponseWriter, r *http.Request)
	CancelOrder(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// HealthHandler handles health-related HTTP requests
type HealthHandler struct {
	logger      *log.Logger
	errorLogger *log.Logger
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
	logger      *log.Logger
	errorLogger *log.Logger
}

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryService *services.CategoryService
	logger          *log.Logger
	errorLogger     *log.Logger
}

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productService *services.ProductService
	logger         *log.Logger
	errorLogger    *log.Logger
}

// CartHandler handles cart-related HTTP requests
type CartHandler struct {
	cartService *services.CartService
	logger      *log.Logger
	errorLogger *log.Logger
}

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderService *services.OrderService
	logger       *log.Logger
	errorLogger  *log.Logger
}
