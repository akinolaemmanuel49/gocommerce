package routes

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/configs"
	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterHealthRoute(config *configs.Config, router *mux.Router, logger, errorLogger *log.Logger) {
	const RouteHealth = "/health"

	// Initialize the handler
	healthHandler := handlers.NewHealthHandler(logger, errorLogger)

	router.Handle(RouteHealth, http.HandlerFunc(healthHandler.CheckHealth)).Methods("GET")
}
