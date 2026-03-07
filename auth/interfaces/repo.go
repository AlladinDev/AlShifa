// Package interfaces provides interfaces for auth module
package interfaces

import (
	"AlShifa/auth/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	Register(ctx context.Context, credientials models.Credientials) error
	UpdateVerificationStatus(ctx context.Context, userID primitive.ObjectID, status bool) error
	SearchCredientials(ctx context.Context, filter bson.M) (*models.Credientials, error)
}
