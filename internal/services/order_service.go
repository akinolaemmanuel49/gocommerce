package services

import (
	"context"
	"fmt"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/queue"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewOrderService creates a new instance of OrderService
func NewOrderService(orderRepository *repositories.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

// CreateOrder creates a new instance of an order and commits it to the database
func (s *OrderService) CreateOrder(ctx context.Context, newOrder *models.CreateOrder) (*models.Order, error) {
	// Check for valid user
	_, err := s.userService.RetrieveUserByID(ctx, newOrder.UserID)
	if err != nil {
		return nil, err
	}

	// Transform CreateOrder to Order
	order := models.NewOrder(newOrder)

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

	return order, nil
}

// RetrieveOrderByID retrieves an order by its ID
// func (s *OrderService) RetrieveOrderByID(ctx context.Context, ID string) (*models.Order, error) {
func (s *OrderService) RetrieveOrderByID(ctx context.Context, ID string, ch *amqp.Channel) (*models.Order, error) {
	order, err := s.orderRepository.FindByID(ctx, ID)
	fmt.Printf("DEBUG: ORDER => %v\n", order)
	fmt.Printf("DEBUG: CHANNEL => %v\n", ch)
	fmt.Println("DEBUG: ", err)
	if err != nil {
		return nil, err
	}

	fmt.Println("DEBUG: RABBITMQ ", s.config.RabbitMQURI)
	queue.PublishOrderNotification(s.config, ch, order.ID, order.UserID, "sms", "DEBUG: CHECK FOR ERRORS", s.logger, s.errorLogger)
	return order, nil
}

// RetrieveAllOrders retrieves paginated orders with optional filters
func (s *OrderService) RetrieveAllOrders(ctx context.Context, filter map[string]interface{}, lastId string, limit int) ([]models.Order, string, error) {
	orders, nextCursor, err := s.orderRepository.FindAll(ctx, filter, lastId, limit)
	if err != nil {
		return nil, "", err
	}

	return orders, nextCursor, nil
}

// ChangeOrderStatusByID alters the status field for an order
func (s *OrderService) ChangeOrderStatusByID(ctx context.Context, ID string, status string, ch *amqp.Channel) error {
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

	if !existingOrder.IsLocked && status == "shipped" || !existingOrder.IsLocked && status == "delivered" {
		updateWithLock := bson.M{"$set": bson.M{"status": status, "isLocked": true, "updatedAt": time.Now()}}
		_, err = s.orderRepository.Update(ctx, ID, updateWithLock)
		if err != nil {
			return err
		}
	}

	// Skip update if the status has no change
	if existingOrder.Status == status {
		queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s status is already set to %s", ID, status), s.logger, s.errorLogger)
		return nil
	}

	// Ensure status field can only ever be progressively updated
	if existingOrder.Status == "shipped" && status == "delivered" || existingOrder.Status == "pending" && status == "shipped" {
		update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
		_, err = s.orderRepository.Update(ctx, ID, update)
		if err != nil {
			return err
		}
		queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s status was changed from %s to %s", ID, existingOrder.Status, status), s.logger, s.errorLogger)
		return nil
	}

	return errors.NewConflictError("Order", "Status", "Status can only be progressively updated pending -> shipped -> delivered")
}

// ChangeOrderShippingAddressByID changes the shipping address for an order, it checks if the resource is locked first
func (s *OrderService) ChangeOrderShippingAddressByID(ctx context.Context, ID string, newAddress *models.UpdateAddress, ch *amqp.Channel) error {
	// Check whether the resource is available for modification
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
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
			Street:  utils.IfNotNil(newAddress.Street, existingOrder.ShippingAddress.Street),
			City:    utils.IfNotNil(newAddress.City, existingOrder.ShippingAddress.City),
			State:   utils.IfNotNil(newAddress.State, existingOrder.ShippingAddress.State),
			Zip:     utils.IfNotNil(newAddress.Zip, existingOrder.ShippingAddress.Zip),
			Country: utils.IfNotNil(newAddress.Country, existingOrder.ShippingAddress.Country),
		}

	update := bson.M{"$set": bson.M{"shippingAddress": shippingAddress, "updatedAt": time.Now()}}
	queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s shipping address was changed", ID), s.logger, s.errorLogger)
	_, err = s.orderRepository.Update(ctx, ID, update)
	return err
}

// AddItemToOrderByID adds a product item to the items field
func (s *OrderService) AddItemToOrderByID(ctx context.Context, ID string, item models.OrderItem) error {
	// Check whether the resource is available for modification
	order, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}
	if order.IsLocked {
		return errors.NewValidationError("Items", "Unable to alter this resource, because is locked")
	}
	if order.IsCancelled {
		return errors.NewValidationError("Items", "Unable to alter this resource, because it has been cancelled")
	}

	update := bson.M{
		"$push": bson.M{"items": item},
		"$inc":  bson.M{"totalPrice": item.Price * float64(item.Quantity)},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err = s.orderRepository.Update(ctx, ID, update)
	return err
}

// RemoveItemFromOrder using the order id removes a product using its product id from the items field
func (s *OrderService) RemoveItemFromOrderByID(ctx context.Context, ID string, productID string) error {
	// Check whether the resource is available for modification
	order, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}
	if order.IsLocked {
		return errors.NewValidationError("Items", "Unable to alter this resource, because is locked")
	}
	if order.IsCancelled {
		return errors.NewValidationError("Items", "Unable to alter this resource, because it has been cancelled")
	}

	order, err = s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	var updatedItems []models.OrderItem
	totalPriceAdjustment := 0.0
	for _, item := range order.Items {
		if item.ProductID == productID {
			totalPriceAdjustment -= item.Price * float64(item.Quantity)
			continue
		}
		updatedItems = append(updatedItems, item)
	}

	update := bson.M{
		"$set": bson.M{"items": updatedItems, "totalPrice": order.TotalPrice + totalPriceAdjustment, "updatedAt": time.Now()},
	}
	_, err = s.orderRepository.Update(ctx, ID, update)
	return err
}

// ConfirmOrderByID sets the IsLocked flag for an order instance to true
func (s *OrderService) ConfirmOrderByID(ctx context.Context, ID string, ch *amqp.Channel) (bool, error) {
	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	// Check if the order is already locked
	if existingOrder.IsLocked {
		queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s has already been confirmed", ID), s.logger, s.errorLogger)
		return true, nil
	}
	// Build query
	lock := bson.M{"isLocked": true}
	// Execute update
	_, err = s.orderRepository.Update(ctx, ID, lock)
	if err != nil {
		return false, err
	}

	queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s was confirmed", ID), s.logger, s.errorLogger)
	return true, nil
}

// CancelOrderByID sets the IsCancelled flag for an order instance to true
func (s *OrderService) CancelOrderByID(ctx context.Context, ID string, ch *amqp.Channel) (bool, error) {
	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	// Check if the order is already cancelled
	if existingOrder.IsCancelled {
		queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s has already been cancelled", ID), s.logger, s.errorLogger)
		return true, nil
	}
	// Build query
	cancel := bson.M{"isCancelled": true, "isLocked": true}
	// Execute update
	_, err = s.orderRepository.Update(ctx, ID, cancel)
	if err != nil {
		return false, err
	}
	queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s was cancelled", ID), s.logger, s.errorLogger)
	return true, nil
}

// DeleteOrderByID sets the IsDeleted flag for an order instance to true (performs a soft-delete)
func (s *OrderService) DeleteOrderByID(ctx context.Context, ID string, ch *amqp.Channel) error {
	// Check for existing order
	existingOrder, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return nil
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

		_, err = s.orderRepository.Update(ctx, ID, order)
		if err != nil {
			return err
		}
	}

	queue.PublishOrderNotification(s.config, ch, existingOrder.ID, existingOrder.UserID, "debug", fmt.Sprintf("Order with ID %s was deleted", ID), s.logger, s.errorLogger)
	return nil
}
