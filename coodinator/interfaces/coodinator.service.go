package interfaces

import (
	"context"

	sharedModels "github.com/AlladinDev/AlShifa/models"
	"github.com/AlladinDev/AlShifa/structs"
)

type ICoodinatorService interface {
	BookAppointment(ctx context.Context, appointmentDetails sharedModels.Appointment) (int, *structs.IAppError)
}
