package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

type IUserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IProductHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type ICategoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IOrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string)
	ReadAll(w http.ResponseWriter, r *http.Request)
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id string)
	UpdateOrderShippingAddress(w http.ResponseWriter, r *http.Request, id string)
	AddCartToOrder(w http.ResponseWriter, r *http.Request, id string, cartID string)
	RemoveCartFromOrder(w http.ResponseWriter, r *http.Request, id string, cartID string)
	ConfirmOrder(w http.ResponseWriter, r *http.Request, id string)
	CancelOrder(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

type ICartHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string)
	ReadAll(w http.ResponseWriter, r *http.Request)
	AddProductToCart(w http.ResponseWriter, r *http.Request, id string)
	RemoveProductFromCart(w http.ResponseWriter, r *http.Request, id string, productID string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}
type UserHandler struct {
	userService *services.UserService
	logger      *log.Logger
	errorLogger *log.Logger
}

type ProductHandler struct {
	productService *services.ProductService
	logger         *log.Logger
	errorLogger    *log.Logger
}

type OrderHandler struct {
	orderService *services.OrderService
	logger       *log.Logger
	errorLogger  *log.Logger
}

type CategoryHandler struct {
	categoryService *services.CategoryService
	logger          *log.Logger
	errorLogger     *log.Logger
}

type CartHandler struct {
	cartService *services.CartService
	logger      *log.Logger
	errorLogger *log.Logger
}
