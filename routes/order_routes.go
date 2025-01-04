package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterOrderRoutes(router *mux.Router, db *mongo.Database, logger, errorLogger *log.Logger, ch *amqp.Channel) {
	const RouteOrders = "/orders"

	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository)
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
			// orderHandler.Read(w, r, id)
			orderHandler.Read(w, r, id, ch)
			// case "PATCH":
			// 	orderHandler.Update(w, r, id)
		}
	})

	router.HandleFunc(RouteOrders+"/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"] // Extract the `id` path parameter

		switch r.Method {
		case "PATCH":
			orderHandler.Delete(w, r, id, ch)
		}
	})
}
