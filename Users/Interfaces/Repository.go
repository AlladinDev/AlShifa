// Package interfaces provides interfaces for user module
package interfaces

import (
	models "AlShifa/Users/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	RegisterUser(ctx context.Context, user models.User) error
	SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error)
	SearchUser(ctx context.Context, filter bson.M) (*models.User, error)
}
