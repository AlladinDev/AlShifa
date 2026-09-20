package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicModule interface {
	GetClinicName(ctx context.Context, clinicID primitive.ObjectID) (string, error)
	ClinicExists(ctx context.Context, filter bson.M) (bool, error)
	GetClinicIDByOwnerID(ctx context.Context, ownerID primitive.ObjectID) (clinicID primitive.ObjectID, err error)
}
