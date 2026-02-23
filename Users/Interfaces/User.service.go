package interfaces

import (
	sharedModels "AlShifa/models"
	structs "AlShifa/structs"
	models "AlShifa/users/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	AddUser(ctx context.Context, user models.User) *structs.IAppError
	SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, *structs.IAppError)
	SearchUser(ctx context.Context, filter bson.M) (*models.User, *structs.IAppError)
	FetchAppointments(ctx context.Context, groupingID string, filter bson.M) ([]sharedModels.Appointment, *structs.IAppError)
	LoginUser(ctx context.Context, email string, password string) (string, *structs.IAppError)
}
