package repositories

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/internal/models"
	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepository struct {
	*BaseRepository
}

// NewCartRepository creates a new instance of CartRepository
func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{
		BaseRepository: NewBaseRepository(db.Collection("carts")),
	}
}

// FindAll retrieves carts based on filters and implements cursor-based pagination
func (r *CartRepository) FindAll(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Cart, string, error) {
	query := bson.M{}
	filter["isDeleted"] = false
	filter["isOrdered"] = false
	if len(filter) > 0 {
		query = filter
	}

	// If lastID is provided, add it to the filter for pagination
	if lastID != "" {
		objID, err := utils.StringToObjectID(lastID)
		if err != nil {
			return nil, "", err
		}
		query["_id"] = bson.M{"$gt": objID} // Fetch carts with IDs greater than lastID
	}

	options := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"_id": 1}) // Sort by _id in ascending order for cursor-based pagination

	cursor, err := r.Collection.Find(ctx, query, options)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	// Decode carts
	var carts []models.Cart
	if err := cursor.All(ctx, &carts); err != nil {
		return nil, "", err
	}

	// Determine the next cursor (last _id in the result)
	var nextCursor string
	if len(carts) > 0 {
		nextCursor = carts[len(carts)-1].ID // Assuming models.Cart.ID corresponds to the MongoDB `_id`
	}

	return carts, nextCursor, nil
}

// FindByID retrieves a cart by its ID
func (r *CartRepository) FindByID(ctx context.Context, ID string) (*models.Cart, error) {
	objectID, err := utils.StringToObjectID(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "isDeleted": false}
	var cart models.Cart

	if err := r.Collection.FindOne(ctx, filter).Decode(&cart); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &cart, nil
}
