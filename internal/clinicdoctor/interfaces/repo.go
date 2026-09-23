package interfaces

import (
	"context"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/dtos"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepo interface {
	RegisterDoctorClinicMapping(ctx context.Context, details *models.ClinicDoctorMapping) error
	CheckClinicDoctorMappingExists(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error)
	CheckAppointmentDatePossible(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDate time.Time) (bool, error)
	FetchClinicWithDoctors(ctx context.Context, filter bson.M) ([]dtos.ClinicDoctorDTO, error)
	FetchDoctorWithClinics(ctx context.Context, filter bson.M) ([]dtos.DoctorWithClinic, error)
}
