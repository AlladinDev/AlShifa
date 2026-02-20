package utils

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseUserID(userID any) (error, primitive.ObjectID) {
	if userID == "" {
		return errors.New("userID missing"), primitive.NilObjectID
	}

	fmt.Println("userid receuved is", userID)
	userIDStr, ok := userID.(string)
	if !ok {
		return errors.New("Invalid UserID failed to parse it into string"), primitive.NilObjectID
	}

	userMongoDBID, mongoErr := primitive.ObjectIDFromHex(userIDStr)
	if mongoErr != nil {
		return errors.New("Invalid UserID type expected mongodbID"), primitive.NilObjectID
	}

	return nil, userMongoDBID
}
