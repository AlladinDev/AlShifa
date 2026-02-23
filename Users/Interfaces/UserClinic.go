package interfaces

import (
	sharedModels "AlShifa/models"
	"AlShifa/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type IUserClinicContract interface {
	FetchAppointments(ctx context.Context, groupBy string, filter bson.M) ([]sharedModels.Appointment, *structs.IAppError)
}
