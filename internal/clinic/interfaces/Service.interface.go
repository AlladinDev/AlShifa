package clinicinterfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/clinic/models"
	response "github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IService interface contains functions that clinic service layer must implement( to beused by handlers)
type IService interface {
	Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) *response.IAppError
	Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, *response.IAppError)
	ClinicExists(ctx context.Context, filter bson.M) (bool, error)
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error)
	DeductClinicMoneyForAppointment(ctx context.Context, clinicID primitive.ObjectID) error
	GetClinicName(ctx context.Context, clinicID primitive.ObjectID) (string, error)
	GetClinicIDByOwnerID(ctx context.Context, ownerID primitive.ObjectID) (clinicID primitive.ObjectID, err error)
}
