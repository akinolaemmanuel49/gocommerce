package services

import (
	"context"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewOrderService creates a new instance of OrderService
func NewOrderService(orderRepository *repositories.OrderRepository, publisher *queue.Publisher, userService *UserService, cartService *CartService) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		publisher:       publisher,
	}
}

// CreateOrder creates a new instance of an order and commits it to the database
func (s *OrderService) CreateOrder(ctx context.Context, userID string, newOrder *models.CreateOrder) (*models.Order, error) {
	// Check for valid user
	user, err := s.userRepository.FindByID(ctx, userID)
	if err == mongo.ErrNoDocuments {
		return nil, errors.NewNotFoundError("User", "ID", userID)
	}
	if err != nil {
		return nil, err
	}

	// Transform CreateOrder to Order
	order := models.NewOrder(newOrder, user)

	// Insert order into the database
	result, err := s.orderRepository.Insert(ctx, order)
	if err != nil {
		return nil, err
	}

	// Convert InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	order.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	// // Publish message
	message := queue.OrderMessage{
		OrderID:          order.ID,
		UserID:           order.UserID,
		EventType:        "OrderCreated",
		Status:           order.Status,
		Message:          "Order created successfully",
		NotificationTime: time.Now(),
	}

	if err := s.publisher.Publish(ctx, message); err != nil {
		s.errorLogger.Printf("Failed to publish message: %v", err)
	}

	return order, nil
}

// RetrieveOrderByID retrieves an order by its ID
func (s *OrderService) RetrieveOrderByID(ctx context.Context, userID string, ID string) (*models.Order, error) {
	order, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if userID != order.UserID {
		return nil, errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	return order, nil
}

// RetrieveAllOrders retrieves paginated orders with optional filters
func (s *OrderService) RetrieveAllOrders(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Order, string, error) {
	orders, nextCursor, err := s.orderRepository.FindAll(ctx, filter, lastID, limit)
	if err != nil {
		return nil, "", err
	}

	return orders, nextCursor, nil
}

// ChangeOrderStatusByID alters the status field for an order
func (s *OrderService) ChangeOrderStatusByID(ctx context.Context, ID string, status string) error {
	// Check if the status is valid
	if status != "pending" && status != "shipped" && status != "delivered" {
		return errors.NewValidationError("Status", "Invalid status value")
	}

	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	if existingOrder.IsCancelled {
		return errors.NewValidationError("Status", "Unable to alter this resource, because it has been cancelled")
	}

	if !existingOrder.IsLocked {
		return errors.NewValidationError("Status", "Cannot modify this resource, because it is not locked")
	}

	// Skip update if the status has no change
	if existingOrder.Status == status {
		return nil
	}

	// Ensure status field can only ever be progressively updated
	if existingOrder.Status == "shipped" && status == "delivered" || existingOrder.Status == "pending" && status == "shipped" {
		update := bson.M{"status": status, "updatedAt": time.Now()}
		_, err = s.orderRepository.Update(ctx, ID, update)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.NewConflictError("Order", "Status", "Status can only be progressively updated pending -> shipped -> delivered")
}

// ChangeOrderShippingAddressByID changes the shipping address for an order, it checks if the resource is locked first
func (s *OrderService) ChangeOrderShippingAddressByID(ctx context.Context, ID string, userID string, newAddress *models.UpdateAddress) error {
	// Check whether the resource is available for modification
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	if userID != existingOrder.UserID {
		return errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	if existingOrder.IsLocked {
		return errors.NewValidationError("ShippingAddress", "Unable to alter this resource, because is locked")
	}
	if existingOrder.IsCancelled {
		return errors.NewValidationError("ShippingAddress", "Unable to alter this resource, because it has been cancelled")
	}

	// Tranform Address
	shippingAddress :=
		&models.Address{
			Street:  models.IfNotNil(newAddress.Street, existingOrder.ShippingAddress.Street),
			City:    models.IfNotNil(newAddress.City, existingOrder.ShippingAddress.City),
			State:   models.IfNotNil(newAddress.State, existingOrder.ShippingAddress.State),
			Zip:     models.IfNotNil(newAddress.Zip, existingOrder.ShippingAddress.Zip),
			Country: models.IfNotNil(newAddress.Country, existingOrder.ShippingAddress.Country),
		}

	update := bson.M{
		"$set": bson.M{
			"shippingAddress": shippingAddress,
			"updatedAt":       time.Now()},
	}
	_, err = s.orderRepository.Update(ctx, ID, update)
	return err
}

// AddCartToOrderByID adds a cart to an order, it checks if the resource is locked first
func (s *OrderService) AddCartToOrderByID(ctx context.Context, ID string, userID string, cartID string) error {
	// Retrieve the order by ID
	order, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	if userID != order.UserID {
		return errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	// Check if the order can be modified
	if order.IsLocked {
		return errors.NewValidationError("Cart", "Unable to alter this resource, because it is locked")
	}
	if order.IsCancelled {
		return errors.NewValidationError("Cart", "Unable to alter this resource, because it has been cancelled")
	}

	// Check if cart is valid
	cart, err := s.cartRepository.FindByID(ctx, cartID)
	if err != nil {
		return errors.NewNotFoundError("Cart", "ID", cartID)
	}

	// Check for cart in order, if present return conflict error
	for _, c := range order.Carts {
		if c.ID == cartID {
			return errors.NewConflictError("Cart", "ID", cartID)
		}
	}

	// Add cart to carts and increment order total price
	order.Carts = append(order.Carts, *cart)
	// Set the total price for orders
	order.TotalPrice += cart.TotalPrice

	// Prepare update fields for the database
	updateFields := bson.M{
		"$set": bson.M{
			"carts":      order.Carts,
			"totalPrice": order.TotalPrice,
			"updatedAt":  time.Now(),
		},
	}

	// Update the order in the database
	_, err = s.orderRepository.Update(ctx, ID, updateFields)
	return err
}

// RemoveCartFromOrderByID removes a cart from an order, it checks if the resource is locked first
func (s *OrderService) RemoveCartFromOrderByID(ctx context.Context, ID string, userID string, cartID string) error {
	// Check whether the resource is available for modification
	order, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	if userID != order.UserID {
		return errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	if order.IsLocked {
		return errors.NewValidationError("Carts", "Unable to alter this resource because it is locked")
	}
	if order.IsCancelled {
		return errors.NewValidationError("Carts", "Unable to alter this resource because it has been cancelled")
	}

	// Check if cart with matching id is present in order carts
	_, err = s.cartRepository.FindByID(ctx, cartID)
	if err != nil {
		return errors.NewNotFoundError("Cart", "ID", cartID)
	}

	// Remove cart with matching id from order carts
	var updatedCart []models.Cart
	totalPriceAdjustment := 0.0
	for _, c := range order.Carts {
		if c.ID == cartID {
			totalPriceAdjustment -= c.TotalPrice
			continue
		}
		updatedCart = append(updatedCart, c)
	}

	// If there are no remaining carts, the total price should be 0
	newTotalPrice := order.TotalPrice + totalPriceAdjustment
	if len(updatedCart) == 0 {
		newTotalPrice = 0
	}

	// Prepare the update fields
	update := bson.M{
		"$set": bson.M{
			"carts":      updatedCart,
			"totalPrice": newTotalPrice,
			"updatedAt":  time.Now(),
		},
	}

	// Update the order in the database
	_, err = s.orderRepository.Update(ctx, ID, update)
	return err
}

// ConfirmOrderByID sets the IsLocked flag for an order instance to true
func (s *OrderService) ConfirmOrderByID(ctx context.Context, ID string, userID string) (bool, error) {
	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}

	if userID != existingOrder.UserID {
		return false, errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	// Check if the order is already locked
	if existingOrder.IsLocked {
		return true, nil
	}
	// Build query
	lock := bson.M{"$set": bson.M{"isLocked": true}}
	// Execute update
	_, err = s.orderRepository.Update(ctx, ID, lock)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CancelOrderByID sets the IsCancelled flag for an order instance to true
func (s *OrderService) CancelOrderByID(ctx context.Context, ID string, userID string) (bool, error) {
	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}

	if userID != existingOrder.UserID {
		return false, errors.NewForbiddenError("You do not have permission to modify this resource")
	}
	// Check if the order is already cancelled
	if existingOrder.IsCancelled {
		return true, nil
	}
	// Build query
	cancel := bson.M{"isCancelled": true, "isLocked": true}
	// Execute update
	_, err = s.orderRepository.Update(ctx, ID, cancel)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteOrderByID sets the IsDeleted flag for an order instance to true (performs a soft-delete)
func (s *OrderService) DeleteOrderByID(ctx context.Context, ID string, userID string) error {
	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if err == mongo.ErrNoDocuments {
		return nil
	}

	if userID != existingOrder.UserID {
		return errors.NewForbiddenError("You do not have permission to modify this resource")
	}

	if existingOrder != nil {
		// Apply transformation, set order IsDeleted field to true
		order := &models.Order{
			CommonFields: models.CommonFields{
				IsDeleted: true,
				UpdatedAt: time.Now(),
			},
		}
		// Check if order status is pending
		if existingOrder.Status == "pending" && !existingOrder.IsCancelled {
			order.IsCancelled = true
			order.IsLocked = true
		}

		deleted := bson.M{"$set": order}

		_, err = s.orderRepository.Update(ctx, ID, deleted)
		if err != nil {
			return err
		}
	}

	return nil
}
