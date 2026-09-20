package interfaces

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	IsAppointmentDatePossible(ctx context.Context, appointmentDate time.Time, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error)
}
