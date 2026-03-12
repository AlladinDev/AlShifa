// Package repository provides repository functions for auth service
package repository

import (
	"context"

	"github.com/AlladinDev/AlShifa/auth/interfaces"
	"github.com/AlladinDev/AlShifa/auth/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

var _ interfaces.IRepository = (*Repository)(nil)

func (r *Repository) Register(ctx context.Context, credientials models.Credientials) error {
	_, err := r.db.Collection("Auth").InsertOne(ctx, credientials)
	return err
}

func (r *Repository) SearchCredientials(ctx context.Context, filter bson.M) (*models.Credientials, error) {
	res := r.db.Collection("Auth").FindOne(ctx, filter)
	var user models.Credientials
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *Repository) UpdateVerificationStatus(ctx context.Context, userID primitive.ObjectID, status bool) error {
	updateOperation := bson.D{{Key: "$set", Value: bson.M{"isVerified": status}}}
	result := r.db.Collection("Auth").FindOneAndUpdate(ctx, bson.M{"_id": userID}, updateOperation)
	return result.Err()
}
