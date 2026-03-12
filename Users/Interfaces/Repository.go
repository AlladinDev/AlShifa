// Package interfaces provides interfaces for user module
package interfaces

import (
	"context"

	models "github.com/AlladinDev/AlShifa/users/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	RegisterUser(ctx context.Context, user models.User) error
	SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error)
}
