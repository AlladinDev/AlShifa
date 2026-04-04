// Package repository provides Repository functions for owner Repository
package repository

import (
	"context"
	"fmt"

	"github.com/AlladinDev/AlShifa/owner/interfaces"
	"github.com/AlladinDev/AlShifa/owner/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(DB *mongo.Database) *Repository {
	return &Repository{
		db: DB,
	}
}

var _ interfaces.IRepository = (*Repository)(nil)

func (r *Repository) RegisterOwner(ctx context.Context, owner models.Owner) error {
	_, err := r.db.Collection("Owner").InsertOne(ctx, owner)
	return err
}
func (r *Repository) GetOwnerDetails(ctx context.Context, filter bson.M) (*models.Owner, error) {
	res := r.db.Collection("Owner").FindOne(ctx, filter)
	var owner models.Owner
	if err := res.Decode(&owner); err != nil {
		return nil, err
	}

	return &owner, nil
}

func (r *Repository) GetOwnerByID(ctx context.Context, ownerID primitive.ObjectID) (*models.Owner, error) {
	fmt.Println(ownerID)
	res := r.db.Collection("Owner").FindOne(ctx, bson.M{"_id": ownerID})
	var owner models.Owner
	if err := res.Decode(&owner); err != nil {
		return nil, err
	}

	return &owner, nil
}
