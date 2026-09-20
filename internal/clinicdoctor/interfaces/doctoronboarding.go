package interfaces

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDoctorOnboarding interface {
	InitiateDoctorClinicOnboarding(ctx context.Context, details *models.ClinicDoctorMapping, userID primitive.ObjectID) (string, *response.IAppError)
}
