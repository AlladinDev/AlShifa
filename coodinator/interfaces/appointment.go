package interfaces

import (
	"context"

	sharedModels "github.com/AlladinDev/AlShifa/models"
	"github.com/AlladinDev/AlShifa/structs"
)

type IAppointmentService interface {
	AddAppointment(ctx context.Context, maxAppointments int, appointmentDetails sharedModels.Appointment) (int, *structs.IAppError)
}
