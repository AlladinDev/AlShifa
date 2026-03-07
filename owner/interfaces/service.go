package interfaces

import (
	"AlShifa/owner/models"
	"AlShifa/structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	RegisterOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, *structs.IAppError)
	GetOwnerDetails(ctx context.Context, filter bson.M) (*models.Owner, *structs.IAppError)
}
