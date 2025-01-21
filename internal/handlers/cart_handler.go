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

func NewCartHandler(cartService *services.CartService, logger, errorLogger *log.Logger) *CartHandler {
	return &CartHandler{cartService: cartService, logger: logger, errorLogger: errorLogger}
}

var _ ICartHandler = (*CartHandler)(nil)

// Create handles POST /carts requests and accepts CreateCart as input [CUSTOMER]
// @Security BearerAuth
// @Summary Create a new cart.
// @Description This endpoint creates a new cart.
// @Tags Carts
// @Accept json
// @Produce json
// @Param cart body models.CreateCart true "Cart Details"
// @Success 201 {object} models.Category "Created cart"
// @Failure 400 "Invalid Request Body"
// @Failure 401 "Not Found"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /carts [post]
func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateCart
	ctx := r.Context()

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to create cart
	cart, err := h.cartService.CreateCart(ctx, &input)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the created cart
	utils.WriteJSON(w, r, http.StatusCreated, cart, h.logger)
}

// Read handles GET /carts/:id requests [CUSTOMER]
// @Security BearerAuth
// @Summary Read a cart
// @Description This endpoint fetches a single cart.
// @Tags Carts
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} models.Cart "Returned cart"
// @Failure 400 "Invalid Cart ID"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /carts [get]
func (h *CartHandler) Read(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Cart"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to get cart by ID
	cart, err := h.cartService.RetrieveCartByID(ctx, ID)
	switch err {
	case nil:
		// No error continue execution
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Cart", "ID", ID), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, cart, h.logger)
}

// ReadAll handles GET /carts requests with optional filters [CUSTOMER]
// @Summary Read all carts
// @Description This endpoint fetches a list of carts with cursor based pagination, optionally filtered by name.
// @Tags Cart
// @Accept json
// @Produce json
// @Param name query string false "Filter carts by name"
// @Param lastID query string false "Last cart id in a page"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.MultipleEntityClientResponse "Returned products and next cursor"
// @Failure 500 "Internal Server Error"
// @Router /carts/all [get]
func (h *CartHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
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

	filter := map[string]interface{}{}
	if name := query.Get("name"); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	carts, nextCursor, err := h.cartService.RetrieveAllCarts(ctx, filter, lastID, limit)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	response := map[string]interface{}{
		"data":       carts,
		"nextCursor": nextCursor,
	}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// AddProductToCart handles PUT /carts/:id/items/add requests [CUSTOMER]
// @Security BearerAuth
// @Summary Add product to cart
// @Description This endpoint adds a product to a cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param cartItem body models.CartItemCreate true "Cart Item"
// @Success 200 {object} models.ClientResponse "Response Message"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /carts/{id}/items/add [put]
func (h *CartHandler) AddProductToCart(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("ID", "Invalid order ID"), h.errorLogger)
		return
	}

	var input models.CartItemCreate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to add new item to cart
	_, err := h.cartService.AddProductToCart(ctx, ID, input.ProductID, input.Quantity)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: "Add cart item was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// RemoveProductFromCart handles PUT /carts/:id/items/remove requests [CUSTOMER]
// @Security BearerAuth
// @Summary Removes product from cart
// @Description This endpoint removes a product from a cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param cartItemDelete body models.CartItemUpdate true "Cart Item Delete"
// @Success 200 {object} models.ClientResponse "Response Message"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /carts/{id}/items/remove [put]
func (h *CartHandler) RemoveProductFromCart(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid cart ID"), h.errorLogger)
		return
	}

	var input models.CartItemUpdate

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to add new item to cart
	_, err := h.cartService.RemoveProductFromCart(ctx, ID, input.ProductID, input.Quantity)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := models.ClientResponse{
		Message: "Remove cart item was successful",
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Delete handles DELETE /carts/:id requests [CUSTOMER]
// @Security BearerAuth
// @Summary Delete cart
// @Description This endpoint deletes a single cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} models.ClientResponse "Response Message"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 500 "Internal Server Error"
// @Router /carts [delete]
func (h *CartHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Initialize context
	ctx := r.Context()

	// Get ID from URL
	ID := utils.GetIDFromURL(r, "id")

	// Validate the ID
	if err := utils.ValidateID(ID, "Cart"); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Call service to delete cart
	if err := h.cartService.DeleteCartByID(ctx, ID); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Build response map
	response := map[string]interface{}{
		"message": fmt.Sprintf("Cart with ID: %s was successfully deleted", ID),
	}

	// Write response to client
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
