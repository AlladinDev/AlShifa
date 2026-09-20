package interfaces

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicDoctor interface {
	IsAppointmentDatePossible(ctx context.Context, appoinmentDate time.Time, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error)
}
