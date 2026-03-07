package interfaces

import (
	"AlShifa/auth/models"
	"AlShifa/structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	Register(ctx context.Context, credientials models.Credientials) *structs.IAppError
	Login(ctx context.Context, email string, password string) (string, *structs.IAppError)
	SearchCredientials(ctx context.Context, filter bson.M) (*models.Credientials, *structs.IAppError)
	UpdateVerificationStatus(ctx context.Context, userID primitive.ObjectID, status bool) *structs.IAppError
}
