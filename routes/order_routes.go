package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterOrderRoutes(config *configs.Config, router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger) {
	const RouteOrders = "/orders"

	// Initialize repositories
	orderRepository := repositories.NewOrderRepository(db)
	userRepository := repositories.NewUserRepository(db)

	// Initialize the publisher
	publisher, err := queue.NewPublisher(config)
	if err != nil {
		errorLogger.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(userRepository)
	orderService := services.NewOrderService(orderRepository, publisher, userService)

	// Initialize the handler
	orderHandler := handlers.NewOrderHandler(orderService, logger, errorLogger)

	router.HandleFunc(RouteOrders, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			orderHandler.ReadAll(w, r)
		case "POST":
			orderHandler.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "GET":
			orderHandler.Read(w, r, id)
			// orderHandler.Read(w, r, id, ch)
			// case "PATCH":
			// orderHandler.Update(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.Delete(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/address", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.UpdateOrderShippingAddress(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.UpdateOrderStatus(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/items/add", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.AddItemToOrder(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/items/{productId}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]               // Extract the `id` path parameter
		productId := mux.Vars(r)["productId"] // Extract the `productId` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.RemoveItemFromOrder(w, r, id, productId)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/confirm", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.ConfirmOrder(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.CancelOrder(w, r, id)
		}
	})
}
