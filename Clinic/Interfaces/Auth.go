package interfaces

import (
	"AlShifa/structs"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthService interface {
	GetEmail(ctx context.Context, id primitive.ObjectID) (string, *structs.IAppError)
}
