package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/gorilla/mux"
)

// RegisterHealthRoute initializes repositories, services and attaches handlers to the router
func RegisterHealthRoute(router *mux.Router, logger, errorLogger *log.Logger) {
	const RouteHealth = "/health"

	// Initialize the handler
	healthHandler := handlers.NewHealthHandler(logger, errorLogger)

	router.Handle(RouteHealth, http.HandlerFunc(healthHandler.CheckHealth)).Methods("GET")
}
