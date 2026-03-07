package interfaces

import (
	structs "AlShifa/structs"
	models "AlShifa/users/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	AddUser(ctx context.Context, user models.User) *structs.IAppError
	SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError)
	
}
