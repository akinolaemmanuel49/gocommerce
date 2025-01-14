package handlers

import (
	"encoding/json"
	"fmt"
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

var _ IOrderHandler = (*OrderHandler)(nil)

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

// func (h *OrderHandler) Read(w http.ResponseWriter, r *http.Request, id string) {
func (h *OrderHandler) Read(w http.ResponseWriter, r *http.Request, id string) {
	// Validate ID
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	// Fetch order by ID
	order, err := h.orderService.RetrieveOrderByID(r.Context(), id)
	switch err {
	case nil:
		// No error
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Order", "ID", id), h.errorLogger)
		return
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

// UpdateOrderStatus handles PUT /orders/:id/status requests :::
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request, id string) {
	// Validate ID
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	var input models.OrderStatusUpdate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update order status
	err := h.orderService.ChangeOrderStatusByID(r.Context(), id, input.Status)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := map[string]string{"message": "Order status update was successful"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// UpdateOrderShippingAddress handles PUT /orders/:id/address requests :::
func (h *OrderHandler) UpdateOrderShippingAddress(w http.ResponseWriter, r *http.Request, id string) {
	// Validate ID
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	var input models.UpdateAddress

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update order address
	err := h.orderService.ChangeOrderShippingAddressByID(r.Context(), id, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := map[string]string{"message": "Order shipping address update was successful"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// AddItemToOrder handles PUT /orders/:id/items/add requests :::
func (h *OrderHandler) AddItemToOrder(w http.ResponseWriter, r *http.Request, id string) {
	// Validate ID
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	var input models.OrderItem

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to add new item to order
	err := h.orderService.AddItemToOrderByID(r.Context(), id, input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := map[string]string{"message": "Add order item was successful"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// RemoveItemFromOrder handles PUT /orders/:id/items/remove/:productID requests :::
func (h *OrderHandler) RemoveItemFromOrder(w http.ResponseWriter, r *http.Request, id string, productID string) {
	// Validate IDs
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}
	if err := utils.ValidateID(productID, "Product"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid product ID"), h.errorLogger)
		return
	}

	// Call service to add new item to order
	err := h.orderService.RemoveItemFromOrderByID(r.Context(), id, productID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := map[string]string{"message": "Remove order item was successful"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// ConfirmOrder handles PUT /orders/:id/confirm requests
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request, id string) {
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to confirm order
	_, err := h.orderService.ConfirmOrderByID(r.Context(), id)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}
	// if !isConfirmed {
	// 	response := map[string]string{"message": fmt.Sprintf("Order with id: %s not confirmed", id)}
	// 	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
	// }

	response := map[string]string{"message": fmt.Sprintf("Order with id: %s confirmed", id)}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// CancelOrder handles PUT /orders/:id/cancel requests
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request, id string) {
	if err := utils.ValidateID(id, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to cancel order
	_, err := h.orderService.CancelOrderByID(r.Context(), id)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}
	// if !isCancelled {
	// 	response := map[string]string{"message": fmt.Sprintf("Order with id: %s cancelled", id)}
	// 	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
	// }
	response := map[string]string{"message": fmt.Sprintf("Order with id: %s cancelled", id)}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Delete handles DELETE /orders/:id/delete requests
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
