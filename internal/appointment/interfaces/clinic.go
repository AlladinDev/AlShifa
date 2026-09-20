package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicModule interface {
	ClinicExists(ctx context.Context, filter bson.M) (bool, error)
	DeductClinicMoneyForAppointment(ctx context.Context, clinicID primitive.ObjectID) error
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error)
}
