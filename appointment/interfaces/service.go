package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/appointment/models"
	"github.com/AlladinDev/AlShifa/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (int, *structs.IAppError)
	FetchAppointments(ctx context.Context, filters bson.M) ([]models.Appointment, *structs.IAppError)
	UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) *structs.IAppError
	FetchAppointmentDaysBooked(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) ([]models.Slot, *structs.IAppError)
}
