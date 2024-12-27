package services

import (
	"context"
	"fmt"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepository *repositories.OrderRepository
}

// NewOrderService creates a new instance of OrderService
func NewOrderService(orderRepository *repositories.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

// GetAllOrders retrieves paginated orders with optional filters
func (s *OrderService) GetAllOrders(ctx context.Context, filter map[string]interface{}, lastId string, limit int) ([]models.Order, string, error) {
	orders, nextCursor, err := s.orderRepository.FindAll(ctx, filter, lastId, limit)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching orders: %w", err)
	}

	return orders, nextCursor, nil
}

// CreateOrder creates a new instance of an order and commits it to the database
func (s *OrderService) CreateOrder(ctx context.Context, newOrder *models.CreateOrder) (*models.Order, error) {
	// Transform CreateOrder to Order
	order := &models.Order{
		UserID:          newOrder.UserID,
		Items:           newOrder.Items,
		TotalPrice:      newOrder.TotalPrice,
		Status:          newOrder.Status,
		ShippingAddress: newOrder.ShippingAddress,
		CommonFields: models.CommonFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert order into the database
	result, err := s.orderRepository.Insert(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("error creating new order: %w", err)
	}

	// Convert InsertedID to string if it's an ObjectID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	order.ID = objectID.Hex() // Set the ID to the string representation of the ObjectID

	return order, nil
}
