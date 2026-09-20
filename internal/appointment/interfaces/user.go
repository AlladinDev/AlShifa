package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUser interface {
	UserExists(ctx context.Context, userID primitive.ObjectID) (bool, error)
}
