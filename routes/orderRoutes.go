package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterOrderRoutes(router *mux.Router, db *mongo.Database, logger *log.Logger) {
	const RouteOrders = "/orders"

	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository, logger)
	orderHandler := handlers.NewOrderHandler(orderService, logger)

	router.HandleFunc(RouteOrders, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			orderHandler.GetAllOrders(w, r)
		case "POST":
			orderHandler.CreateOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
