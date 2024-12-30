package repositories

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	*BaseRepository
}

// NewOrderRepository creates a new instance of OrderRepository
func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(db.Collection("orders")),
	}
}

// FindAll retrieves orders based on filters and implements cursor-based pagination
func (r *OrderRepository) FindAll(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Order, string, error) {
	query := bson.M{}
	if len(filter) > 0 {
		query = filter
	}

	// If lastID is provided, add it to the filter for pagination
	if lastID != "" {
		objID, err := utils.StringToObjectID(lastID)
		if err != nil {
			return nil, "", errors.NewValidationError("nextCursor", "must be a valid ObjectID")
		}
		query["_id"] = bson.M{"$gt": objID} // Fetch orders with IDs greater than lastID
	}

	options := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"_id": 1}) // Sort by _id in ascending order for cursor-based pagination

	cursor, err := r.Collection.Find(ctx, query, options)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	// Decode orders
	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, "", err
	}

	// Determine the next cursor (last _id in the result)
	var nextCursor string
	if len(orders) > 0 {
		nextCursor = orders[len(orders)-1].ID // Assuming models.Order.ID corresponds to the MongoDB `_id`
	}

	return orders, nextCursor, nil
}

// FindByID retrieves an order by its ID
func (r *OrderRepository) FindByID(ctx context.Context, ID string) (*models.Order, error) {
	var order models.Order
	objectID, err := utils.StringToObjectID(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	if err := r.Collection.FindOne(ctx, filter).Decode(&order); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Order", "ID", ID)
		}
		return nil, err
	}
	return &order, nil
}
