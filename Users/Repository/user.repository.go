// Package repository provides repository functions for user module
package repository

import (
	interfaces "AlShifa/Users/Interfaces"
	models "AlShifa/Users/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB *mongo.Database
}

func ReturnNewRepository(db *mongo.Database) *Repository {
	return &Repository{
		DB: db,
	}
}

var _ interfaces.IRepository = (*Repository)(nil)

func (repo *Repository) RegisterUser(ctx context.Context, user models.User) error {
	_, err := repo.DB.Collection("User").InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) SearchUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	result := repo.DB.Collection("User").FindOne(ctx, bson.M{"_id": userID})
	var user models.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Repository) SearchUser(ctx context.Context, filter bson.M) (*models.User, error) {
	result := repo.DB.Collection("User").FindOne(ctx, filter)
	var user models.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
