package repository

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/users/interfaces"
	"github.com/AlladinDev/AlShifa/internal/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(DB *mongo.Database) interfaces.IRepository {
	return &Repository{
		db: DB,
	}
}

func (r *Repository) CreateUser(ctx context.Context, details *models.User) error {
	_, err := r.db.Collection("User").InsertOne(ctx, details)
	return err
}

func (r *Repository) GetUserDetails(ctx context.Context, filter bson.M) (*models.User, error) {
	user := models.User{}
	if err := r.db.Collection("User").FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
