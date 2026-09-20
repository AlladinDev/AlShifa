// Package clinicinterfaces contains interfaces for clinic module
package clinicinterfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/clinic/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IRepository defines the methods for loose coupling between the repository and its implementation.
type IRepository interface {
	ClinicExists(ctx context.Context, filter bson.M) (bool, error)
	Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) error
	Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, error)
	FetchSingleClinic(ctx context.Context, filter bson.M) (*models.Clinic, error)
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error)
	DeductClinicWallet(ctx context.Context, amountToDeduct int, clinicID primitive.ObjectID) error
}
