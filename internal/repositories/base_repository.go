package repositories

import (
	"context"

	"github.com/akinolaemmanuel49/gocommerce/utils"
	"go.mongodb.org/mongo-driver/bson"
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

// Update updates a document in a collection based on their _id.
func (r *BaseRepository) Update(ctx context.Context, id string, updateFields interface{}) (*mongo.UpdateResult, error) {
	objectID, err := utils.StringToObjectID(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateFields}
	return r.Collection.UpdateOne(ctx, filter, update)
}

// Delete deletes a document from a collection based on their _id.
func (r *BaseRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := utils.StringToObjectID(id)
	if err != nil {
		return nil, err
	}
	filter := map[string]interface{}{"_id": objectID}
	return r.Collection.DeleteOne(ctx, filter)
}
