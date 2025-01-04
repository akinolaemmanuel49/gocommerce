package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IUserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

type IProductHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

type ICategoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string)
	ReadAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

type IOrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request, id string, ch *amqp.Channel)
	ReadAll(w http.ResponseWriter, r *http.Request)
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id string, ch *amqp.Channel)
	UpdateOrderShippingAddress(w http.ResponseWriter, r *http.Request, id string, ch *amqp.Channel)
	AddItemToOrder(w http.ResponseWriter, r *http.Request, id string)
	RemoveItemFromOrder(w http.ResponseWriter, r *http.Request, id string, productID string)
	ConfirmOrder(w http.ResponseWriter, r *http.Request, id string, ch *amqp.Channel)
	CancelOrder(w http.ResponseWriter, r *http.Request, id string, ch *amqp.Channel)
	Delete(w http.ResponseWriter, r *http.Request, id string, ch *amqp.Channel)
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
