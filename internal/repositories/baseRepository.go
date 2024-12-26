package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository struct {
	Collection *mongo.Collection
}

func NewBaseRepository(collection *mongo.Collection) *BaseRepository {
	return &BaseRepository{Collection: collection}
}

// Insert adds a new document to a collection.
func (r *BaseRepository) Insert(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, document)
}

// FindByID retrieves documents based on their _id.
func (r *BaseRepository) FindByID(ctx context.Context, id string, result interface{}) error {
	filter := map[string]interface{}{"_id": id}
	return r.Collection.FindOne(ctx, filter).Decode(result)

}

// Update updates a document in a collection based on their _id.
func (r *BaseRepository) Update(ctx context.Context, id string, update interface{}) (*mongo.UpdateResult, error) {
	filter := map[string]interface{}{"_id": id}
	return r.Collection.UpdateOne(ctx, filter, update)
}

// Delete deletes a document from a collection based on their _id.
func (r *BaseRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	filter := map[string]interface{}{"_id": id}
	return r.Collection.DeleteOne(ctx, filter)
}
