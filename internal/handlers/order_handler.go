package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/services"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewOrderHandler(orderService *services.OrderService, logger, errorLogger *log.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, logger: logger, errorLogger: errorLogger}
}

// Compile-time check that OrderHandler implements HandlerInterface
var _ HandlerInterface = (*OrderHandler)(nil)

// Create handles POST /orders requests and accepts CreateOrder as input
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateOrder

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create order
	order, err := h.orderService.CreateOrder(r.Context(), &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the created order
	utils.WriteJSON(w, r, http.StatusCreated, order, h.logger)
}

func (h *OrderHandler) Read(w http.ResponseWriter, r *http.Request, id string) {
	// Validate ID
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	// Fetch order by ID
	order, err := h.orderService.RetrieveOrderByID(r.Context(), id)
	switch err {
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Order", "ID", id), h.errorLogger)
		return
	case nil:
		// No error
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the order
	utils.WriteJSON(w, r, http.StatusOK, order, h.logger)
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

	orders, nextCursor, err := h.orderService.RetrieveAllOrders(ctx, filter, lastID, limit)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	response := map[string]interface{}{
		"data":       orders,
		"nextCursor": nextCursor,
	}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Update handles PATCH /orders/:id requests
func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	panic("unimplemented") // TODO
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	// Validate ID
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	// Delete order
	if err := h.orderService.DeleteOrderByID(r.Context(), id); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with confirmation
	response := map[string]string{"message": "Order successfully deleted"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
