package services

import (
	"context"
	"time"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewOrderService creates a new instance of OrderService
func NewOrderService(orderRepository *repositories.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
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
func (s *OrderService) RetrieveOrderByID(ctx context.Context, ID string) (*models.Order, error) {
	order, err := s.orderRepository.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

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

func (s *OrderService) ChangeOrderStatus(ctx context.Context, ID string, status string) error {
	if status != "pending" && status != "shipped" && status != "delivered" {
		return errors.NewValidationError("Status", "Invalid status value")
	}

	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
	_, err := s.orderRepository.Update(ctx, ID, update)
	return err
}

func (s *OrderService) UpdateShippingAddress(ctx context.Context, ID string, address models.Address) error {
	update := bson.M{"$set": bson.M{"shippingAddress": address, "updatedAt": time.Now()}}
	_, err := s.orderRepository.Update(ctx, ID, update)
	return err
}

func (s *OrderService) AddItemToOrder(ctx context.Context, ID string, item models.OrderItem) error {
	update := bson.M{
		"$push": bson.M{"items": item},
		"$inc":  bson.M{"totalPrice": item.Price * float64(item.Quantity)},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err := s.orderRepository.Update(ctx, ID, update)
	return err
}

func (s *OrderService) RemoveItemFromOrder(ctx context.Context, ID string, productID string) error {
	order, err := s.orderRepository.FindByID(ctx, ID)
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
