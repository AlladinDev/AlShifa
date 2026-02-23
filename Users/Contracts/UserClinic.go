//Package contracts provides layers to communicate between modules with least dependency either through dependency injection or http calls
package contracts

import (
	"AlShifa/models"
	"AlShifa/structs"
	"AlShifa/users/interfaces"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserClinicContract struct {
	clinicService interfaces.IUserClinicContract
}

func NewUserClinicContract(service interfaces.IUserClinicContract) *UserClinicContract {
	return &UserClinicContract{
		clinicService: service,
	}
}

var _ = interfaces.IUserClinicContract(nil)

func (c *UserClinicContract) FetchAppointments(ctx context.Context, groupingID string, filter bson.M) ([]models.Appointment, *structs.IAppError) {
	return c.clinicService.FetchAppointments(ctx, groupingID, filter)
}
