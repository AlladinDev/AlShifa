package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/dtos"
	"github.com/AlladinDev/AlShifa/owner/models"
	"github.com/AlladinDev/AlShifa/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	RegisterOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, *structs.IAppError)
	GetOwnerDetails(ctx context.Context, filter bson.M) (*models.Owner, *structs.IAppError)
	LoginOwner(ctx context.Context, loginDetails dtos.LoginEmailPasswordDTO) (string, *structs.IAppError)
}
