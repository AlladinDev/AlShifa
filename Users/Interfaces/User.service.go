package interfaces

import (
	structs "AlShifa/Structs"
	models "AlShifa/Users/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	AddUser(ctx context.Context, user models.User) *structs.IAppError
	SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError)
	SearchUser(ctx context.Context, filter bson.M) (*models.User, *structs.IAppError)
	LoginUser(ctx context.Context, email string, password string) (string, *structs.IAppError)
}
