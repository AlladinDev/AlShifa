package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/users/models"
	"go.mongodb.org/mongo-driver/bson"
)

type IRepository interface {
	CreateUser(ctx context.Context, details *models.User) error
	GetUserDetails(ctx context.Context, filter bson.M) (*models.User, error)
}
