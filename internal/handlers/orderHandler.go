package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
)

func NewOrderHandler(orderService *services.OrderService, logger *log.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, logger: logger}
}

// Compile-time check that OrderHandler implements HandlerInterface
var _ HandlerInterface = (*OrderHandler)(nil)

// Create handles POST /orders requests and accepts CreateOrder as input
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateOrder

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call service to create order
	order, err := h.orderService.CreateOrder(r.Context(), &input)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Respond with the created order
	writeJSON(w, http.StatusCreated, order)
}

// Read handles GET /orders/:id requests
func (h *OrderHandler) Read(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented") // TODO
}

// ReadAll handles GET /orders requests with optional filters
func (h *OrderHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Log to stdout
	h.logger.Printf("%v %v", r.Method, r.URL.Path)

	ctx := r.Context()

	// Parse query parameters for filters and pagination
	query := r.URL.Query()
	lastID := query.Get("lastID")
	limitStr := query.Get("limit")

	limit := 10 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Build filter map
	filter := map[string]interface{}{}
	if category := query.Get("category"); category != "" {
		filter["category"] = category
	}
	if priceMin := query.Get("priceMin"); priceMin != "" {
		if min, err := strconv.ParseFloat(priceMin, 64); err == nil {
			filter["price"] = map[string]interface{}{"$gte": min}
		}
	}
	if priceMax := query.Get("priceMax"); priceMax != "" {
		if max, err := strconv.ParseFloat(priceMax, 64); err == nil {
			if priceFilter, exists := filter["price"].(map[string]interface{}); exists {
				priceFilter["$lte"] = max
			} else {
				filter["price"] = map[string]interface{}{"$lte": max}
			}
		}
	}

	orders, nextCursor, err := h.orderService.GetAllOrders(ctx, nil, lastID, limit)
	if err != nil {
		http.Error(w, "Error fetching orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       orders,
		"nextCursor": nextCursor,
	}
	writeJSON(w, http.StatusOK, response)
}

// Update handles PATCH /orders/:id requests
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented") // TODO
}

// DELETE handles PATCH /orders/:id/delete requests
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented") // TODO
}
