package interfaces

import (
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	"context"
)

// IService interface contains functions that clinic service layer must implement( to beused by handlers)
type IService interface {
	RegisterClinic(ctx context.Context, ownerID string, clinic models.Clinic) *structs.IAppError
	RegisterClinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	SearchClinic(ctx context.Context, clinicName string) ([]models.Clinic, *structs.IAppError)
}
