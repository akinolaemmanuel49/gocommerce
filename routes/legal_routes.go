package routes

import (
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/internal/handlers"
	"github.com/gorilla/mux"
)

// RegisterLegalRoute initializes repositories, services and attaches handlers to the router
func RegisterLegalRoutes(router *mux.Router) {
	// Initialize the handler
	legalHandler := handlers.NewLegalHandler()

	router.Handle("/license/", http.HandlerFunc(legalHandler.GetLicense)).Methods("GET")
	router.Handle("/terms/", http.HandlerFunc(legalHandler.GetTerms)).Methods("GET")
}
