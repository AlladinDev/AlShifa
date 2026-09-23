package interfaces

import (
	"context"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/dtos"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	IsAppointmentDatePossible(ctx context.Context, appointmentDate time.Time, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error)
	FetchDoctorWithClinics(ctx context.Context, filter bson.M) ([]dtos.DoctorWithClinic, *response.IAppError)
	FetchClinicWithDoctors(ctx context.Context, filter bson.M) ([]dtos.ClinicDoctorDTO, *response.IAppError)
}
