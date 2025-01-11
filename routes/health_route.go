package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHealthRoute(router *mux.Router, logger, errorLogger *log.Logger) {
	// func RegisterHealthRoute(router *mux.Router, logger, errorLogger *log.Logger, config *configs.Config) {
	const RouteHealth = "/health"
	// healthHandler handles GET /health requests
	router.HandleFunc(RouteHealth, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Log request using the info logger
			// logger.Println("Listening on Port: ", config.Port)
			logger.Printf("%s %d %s [User-Agent: %s]", r.Method, http.StatusOK, r.URL.Path, r.UserAgent())
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}
	})
}
