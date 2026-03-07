// Package interfaces provides interfaces for owner module
package interfaces

import (
	"AlShifa/owner/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	RegisterOwner(ctx context.Context, owner models.Owner) error
	GetOwnerDetails(ctx context.Context, filter bson.M) (*models.Owner, error)
	GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, error)
}
