// Package interfaces contains interfaces for Clinic module
package interfaces

import (
	"AlShifa/Clinic/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IRepository defines the methods for loose coupling between the repository and its implementation.
type IRepository interface {
	RegisterClinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) error
	RegisterClinicOwner(ctx context.Context, owner models.Owner) error
	GetOwnerDetails(ctx context.Context, filter bson.M) ([]models.Owner, error)
	SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, error)
	RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error
	SearchDoctors(ctx context.Context, filter bson.M) ([]models.Doctor, error)
}
