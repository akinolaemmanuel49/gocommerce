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

type ProductRepository struct {
	*BaseRepository
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		BaseRepository: NewBaseRepository(db.Collection("products")),
	}
}

// FindAll retrieves products based on filters and implements cursor-based pagination.
func (r *ProductRepository) FindAll(ctx context.Context, filter map[string]interface{}, lastID string, limit int) ([]models.Product, string, error) {
	query := bson.M{}
	filter["isDeleted"] = false
	if len(filter) > 0 {
		query = filter
	}

	// If lastID is provided, add it to the filter for pagination
	if lastID != "" {
		objID, err := utils.StringToObjectID(lastID)
		if err != nil {
			return nil, "", errors.NewValidationError("nextCursor", "must be a valid ObjectID")
		}
		query["_id"] = bson.M{"$gt": objID} // Fetch products with IDs greater than lastID
	}

	options := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"_id": 1}) // Sort by _id in ascending order for cursor-based pagination

	cursor, err := r.Collection.Find(ctx, query, options)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	// Decode products
	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, "", err
	}

	// Determine the next cursor (last _id in the result)
	var nextCursor string
	if len(products) > 0 {
		nextCursor = products[len(products)-1].ID // Assuming models.Product.ID corresponds to the MongoDB `_id`
	}

	return products, nextCursor, nil
}

// FindByID retrieves a product by its ID
func (r *ProductRepository) FindByID(ctx context.Context, ID string) (*models.Product, error) {
	objID, err := utils.StringToObjectID(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "isDeleted": false}
	var product models.Product

	if err := r.Collection.FindOne(ctx, filter).Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Product", "ID", ID)
		}
		return nil, err
	}
	return &product, nil
}
