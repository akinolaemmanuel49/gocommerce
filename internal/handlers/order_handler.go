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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewOrderHandler(orderService *services.OrderService, logger, errorLogger *log.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, logger: logger, errorLogger: errorLogger}
}

var _ IOrderHandler = (*OrderHandler)(nil)

// Create handles POST /orders requests and accepts CreateOrder as input
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Initialize request body
	var input models.CreateOrder

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create order
	order, err := h.orderService.CreateOrder(ctx, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusCreated, order, h.logger)
}

func (h *OrderHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to read order by ID
	order, err := h.orderService.RetrieveOrderByID(ctx, ID)

	switch err {
	case nil:
		// No error
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Order", "ID", ID), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, order, h.logger)
}

// ReadAll handles GET /orders requests with optional filters
func (h *OrderHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Initialize context
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

	// Build filter map from query parameters
	filter := map[string]interface{}{}
	if status := query.Get("status"); status != "" {
		filter["status"] = status
	}
	if isCancelled := query.Get("isCancelled"); isCancelled != "" {
		if b, err := strconv.ParseBool(isCancelled); err == nil {
			filter["isCancelled"] = b
		}
	}
	if isLocked := query.Get("isLocked"); isLocked != "" {
		if b, err := strconv.ParseBool(isLocked); err == nil {
			filter["isLocked"] = b
		}
	}
	dateStart := query.Get("dateStart")
	dateEnd := query.Get("dateEnd")
	if dateStart != "" && dateEnd != "" {
		filter["createdAt"] = bson.M{"$gte": dateStart, "$lte": dateEnd}
	}

	// Call service to read all orders
	orders, nextCursor, err := h.orderService.RetrieveAllOrders(ctx, filter, lastID, limit)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"data":       orders,
		"nextCursor": nextCursor,
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// UpdateOrderStatus handles PUT /orders/:id/status requests
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Initialize request body
	var input models.OrderStatusUpdate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update order status
	err := h.orderService.ChangeOrderStatusByID(ctx, ID, input.Status)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": "Order status update was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// UpdateOrderShippingAddress handles PUT /orders/:id/address requests
func (h *OrderHandler) UpdateOrderShippingAddress(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Initialize request body
	var input models.UpdateAddress

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update order address
	err := h.orderService.ChangeOrderShippingAddressByID(ctx, ID, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": "Order shipping address update was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// AddCartToOrder handles PUT /orders/:id/carts/add requests
func (h *OrderHandler) AddCartToOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")
	cartID := utils.GetIDFromURL(r, "cartID")

	// Validate the IDs
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}
	if err := utils.ValidateID(cartID, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("cartID", "Invalid cart ID"), h.errorLogger)
		return
	}

	// Call service to add new cart to order
	err := h.orderService.AddCartToOrderByID(ctx, ID, cartID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": "Add cart to order was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// RemoveCartFromOrder handles PUT /orders/:id/carts/remove requests
func (h *OrderHandler) RemoveCartFromOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")
	cartID := utils.GetIDFromURL(r, "cartID")

	// Validate the IDs
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}
	if err := utils.ValidateID(cartID, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("cartID", "Invalid cart ID"), h.errorLogger)
		return
	}

	// Call service to add new cart to order
	err := h.orderService.RemoveCartFromOrderByID(ctx, ID, cartID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := models.ClientResponse{
		Message: "Remove cart from order was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// ConfirmOrder handles PUT /orders/:id/confirm requests
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to confirm order
	_, err := h.orderService.ConfirmOrderByID(ctx, ID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: fmt.Sprintf("Order with id: %s confirmed", ID),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// CancelOrder handles PUT /orders/:id/cancel requests
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to cancel order
	_, err := h.orderService.CancelOrderByID(ctx, ID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: fmt.Sprintf("Order with id: %s cancelled", ID),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Delete handles DELETE /orders/:id requests
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to delete order
	if err := h.orderService.DeleteOrderByID(ctx, ID); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: fmt.Sprintf("Order with ID: %s was successfully deleted", ID),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
