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

func NewCartHandler(cartService *services.CartService, logger, errorLogger *log.Logger) *CartHandler {
	return &CartHandler{cartService: cartService, logger: logger, errorLogger: errorLogger}
}

var _ ICartHandler = (*CartHandler)(nil)

// Create handles POST /carts requests and accepts CreateCart as input
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

// Read handles GET /carts/:id requests
func (h *CartHandler) Read(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	// Validate ID
	if err := utils.ValidateID(id, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid cart ID"), h.errorLogger)
		return
	}

	// Fetch cart by ID
	cart, err := h.cartService.RetrieveCartByID(ctx, id)
	switch err {
	case nil:
		// No error
	case mongo.ErrNoDocuments:
		errors.HandleError(w, r, errors.NewNotFoundError("Cart", "ID", id), h.errorLogger)
		return
	default:
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with the cart
	utils.WriteJSON(w, r, http.StatusOK, cart, h.logger)
}

// ReadAll handles GET /carts requests with optional filters
func (h *CartHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
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

	filter := map[string]interface{}{}

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

// AddProductToCart handles PUT /carts/:id/items/add requests
func (h *CartHandler) AddProductToCart(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	// Validate ID
	if err := utils.ValidateID(id, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid order ID"), h.errorLogger)
		return
	}

	var input models.CartItem

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("", "Invalid request body"), h.errorLogger)
		return
	}

	// Call service to add new item to cart
	_, err := h.cartService.AddProductToCart(ctx, input.CartID, input.Product.ID, input.Quantity)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := map[string]string{"message": "Add cart item was successful"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// RemoveProductFromCart handles PUT /carts/:id/items/remove requests
func (h *CartHandler) RemoveProductFromCart(w http.ResponseWriter, r *http.Request, id string, productID string) {
	ctx := r.Context()
	// Validate IDs
	if err := utils.ValidateID(id, "Cart"); err != nil {
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
	_, err := h.cartService.RemoveProductFromCart(ctx, id, input.ProductID, input.Quantity)
	if err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with success message
	response := map[string]string{"message": "Remove cart item was successful"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}

// Delete handles DELETE /carts/:id requests
func (h *CartHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	// Validate ID
	if err := utils.ValidateID(id, "Cart"); err != nil {
		errors.HandleError(w, r, errors.NewValidationError("id", "Invalid cart ID"), h.errorLogger)
		return
	}

	// Delete cart
	if err := h.cartService.DeleteCartByID(ctx, id); err != nil {
		errors.HandleError(w, r, err, h.errorLogger)
		return
	}

	// Respond with confirmation
	response := map[string]string{"message": "Cart successfully deleted"}
	utils.WriteJSON(w, r, http.StatusOK, response, h.logger)
}
