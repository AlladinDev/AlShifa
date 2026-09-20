package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseUserID(userID any) (primitive.ObjectID, error) {
	if userID == "" || userID == nil {
		return primitive.NilObjectID, errors.New("userID missing")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid UserID failed to parse it into string")
	}

	userMongoDBID, mongoErr := primitive.ObjectIDFromHex(userIDStr)
	if mongoErr != nil {
		return primitive.NilObjectID, errors.New("invalid UserID type expected mongodbID")
	}

	return userMongoDBID, nil
}
