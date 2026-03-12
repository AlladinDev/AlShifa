package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/structs"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthService interface {
	GetEmail(ctx context.Context, id primitive.ObjectID) (string, *structs.IAppError)
}
