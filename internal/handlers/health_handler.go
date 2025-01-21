package handlers

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/utils"
)

func NewHealthHandler(logger, errorLogger *log.Logger) *HealthHandler {
	return &HealthHandler{logger: logger, errorLogger: errorLogger}
}

// healthHandler handles GET /health requests
// @Summary Health check
// @Description Returns the health status of the API
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /health [get]
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, r, http.StatusOK, "OK", h.logger)

}
