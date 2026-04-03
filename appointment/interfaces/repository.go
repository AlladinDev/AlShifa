// Package interfaces provides interfaces for appointment module
package interfaces

import (
	"github.com/AlladinDev/AlShifa/appointment/models"
	sharedModels "github.com/AlladinDev/AlShifa/models"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	AddAppointment(ctx context.Context, clinicMaxAppointment int, appointmentDetails sharedModels.Appointment) (int, error)
	FetchAppointments(ctx context.Context, filters bson.M) ([]sharedModels.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) error
	FetchAppointmentDaysBooked(ctx context.Context, maxAppointments int, doctorID primitive.ObjectID, clinicID primitive.ObjectID) ([]models.Slot, error)
}
