package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/appointment/models"

	"github.com/AlladinDev/AlShifa/internal/shared/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	AddAppointment(ctx context.Context, userID primitive.ObjectID, appointmentDetails models.Appointment) (slot int, appointmentToken string, err *response.IAppError)
	//FetchAppointments(ctx context.Context, filters bson.M) ([]models.Appointment, *response.IAppError)
	UpdateAppointmentStatus(ctx context.Context, appointmentID primitive.ObjectID, status bool) *response.IAppError
}
