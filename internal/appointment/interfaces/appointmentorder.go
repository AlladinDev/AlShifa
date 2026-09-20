package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/appointment/models"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAppointmentOrder interface {
	CreateOrder(ctx context.Context, appointmentID primitive.ObjectID, userID primitive.ObjectID, planID primitive.ObjectID) (appointmentOrder *models.AppointmentBookingPaymentOrder, err *response.IAppError)
}
