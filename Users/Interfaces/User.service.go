package interfaces

import (
	"context"

	structs "github.com/AlladinDev/AlShifa/structs"
	models "github.com/AlladinDev/AlShifa/users/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	AddUser(ctx context.Context, user models.User) *structs.IAppError
	SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError)
	
}
