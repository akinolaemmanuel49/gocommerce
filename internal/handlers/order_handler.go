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
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
		errorLogger:  errorLogger,
	}
}

var _ IOrderHandler = (*OrderHandler)(nil)

// Create handles POST /orders requests and accepts CreateOrder as req [CUSTOMER]
// @Security BearerAuth
// @Summary Create a new order.
// @Description This endpoint creates a new order.
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrder true "Order Details"
// @Success 201 {object} models.Order "Created order"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Initialize request body
	var req models.CreateOrder

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create order
	order, err := h.orderService.CreateOrder(ctx, claims.UserID, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusCreated, order, h.logger)
}

// Read handles GET /orders/:id requests [CUSTOMER]
// @Security BearerAuth
// @Summary Retrieve an order by ID.
// @Description This endpoint retrieves a specific order by its ID. Only authorized users can access their own orders.
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order "Order details"
// @Failure 400 "Invalid order ID"
// @Failure 401 "Unauthorized"
// @Failure 404 "Order not found"
// @Failure 500 "Internal server error"
// @Router /orders/{id} [get]
func (h *OrderHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to read order by ID
	order, err := h.orderService.RetrieveOrderByID(ctx, claims.UserID, ID)

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

// ReadAll handles GET /orders/all requests with optional filters
// @Security BearerAuth
// @Summary Retrieve a list of orders with optional filters and pagination.
// @Description This endpoint retrieves all orders for the authenticated user. Filters and pagination can be applied.
// @Tags Orders
// @Accept json
// @Produce json
// @Param lastID query string false "Cursor for pagination (last ID from previous result)"
// @Param limit query int false "Number of records to retrieve (default: 10)"
// @Param status query string false "Filter by order status (e.g., pending, completed)"
// @Param isCancelled query boolean false "Filter by cancellation status"
// @Param isLocked query boolean false "Filter by locked status"
// @Param dateStart query string false "Start date for filtering orders (format: YYYY-MM-DD)"
// @Param dateEnd query string false "End date for filtering orders (format: YYYY-MM-DD)"
// @Success 200 {object} models.MultipleEntityClientResponse "List of orders and pagination cursor"
// @Failure 400 "Invalid query parameters"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /orders/all [get]
func (h *OrderHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

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
	// Convert string to object id
	objectID, err := utils.StringToObjectID(claims.UserID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}
	filter["userId"] = objectID
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
	response := models.MultipleEntityClientResponse{
		Data:       orders,
		NextCursor: nextCursor,
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// UpdateOrderStatus handles PUT /orders/:id/status requests
// @Summary Update order status
// @Description Update the status of an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body models.OrderStatusUpdate true "New order status"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Get claims from context and check if user is an admin
	_, err := utils.IsAdmin(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Initialize request body
	var req models.OrderStatusUpdate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update order status
	err = h.orderService.ChangeOrderStatusByID(ctx, ID, req.Status)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: "Order status update was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// UpdateOrderShippingAddress handles PUT /orders/:id/address requests
// @Summary Update order shipping address
// @Description Update the shipping address of an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param address body models.UpdateAddress true "New shipping address"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id}/address [put]
func (h *OrderHandler) UpdateOrderShippingAddress(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Initialize request body
	var req models.UpdateAddress

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to update order address
	err = h.orderService.ChangeOrderShippingAddressByID(ctx, ID, claims.UserID, &req)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: "Order shipping address update was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// AddCartToOrder handles PUT /orders/:id/carts/add requests
// @Summary Add a cart to an order
// @Description Add a cart to an existing order by IDs
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param cartID path string true "Cart ID"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id}/carts/add/{cartID} [put]
func (h *OrderHandler) AddCartToOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

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
	err = h.orderService.AddCartToOrderByID(ctx, ID, claims.UserID, cartID)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: "Add cart to order was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// RemoveCartFromOrder handles PUT /orders/:id/carts/remove requests
// @Summary Remove a cart from an order
// @Description Remove a cart from an existing order by IDs
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param cartID path string true "Cart ID"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id}/carts/remove/{cartID} [put]
func (h *OrderHandler) RemoveCartFromOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

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
	err = h.orderService.RemoveCartFromOrderByID(ctx, ID, claims.UserID, cartID)
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
// @Summary Confirm an order
// @Description Confirm an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id}/confirm [put]
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to confirm order
	_, err = h.orderService.ConfirmOrderByID(ctx, ID, claims.UserID)
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
// @Summary Cancel an order
// @Description Cancel an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to cancel order
	_, err = h.orderService.CancelOrderByID(ctx, ID, claims.UserID)
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
// @Summary Delete an order
// @Description Delete an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.ClientResponse
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /orders/{id} [delete]
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get claims from context
	claims, err := utils.IsAuthorized(ctx)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate ID
	if err := utils.ValidateID(ID, "Order"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	// Call service to delete order
	if err := h.orderService.DeleteOrderByID(ctx, ID, claims.UserID); err != nil {
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
