package interfaces

import (
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// IService interface contains functions that clinic service layer must implement( to beused by handlers)
type IService interface {
	RegisterClinic(ctx context.Context, ownerID string, clinic models.Clinic) *structs.IAppError
	RegisterClinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError)
	LoginClinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError)
	LoginDoctor(ctx context.Context, email string, password string) (string, *structs.IAppError)
}
