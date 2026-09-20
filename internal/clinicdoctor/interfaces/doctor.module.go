package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDoctorModule interface {
	IsClinicAllowedToOnboardDoctor(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) (bool, error)
	GetDoctorEmail(ctx context.Context, doctorID primitive.ObjectID) (string, error)
	DoctorExists(ctx context.Context, filter bson.M) (bool, error)
}
