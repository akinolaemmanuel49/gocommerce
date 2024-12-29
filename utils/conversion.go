package utils

import (
	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObjectID(ID string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return primitive.NilObjectID, errors.NewValidationError("ID", "invalid ID format")
	}
	return objectID, nil
}
