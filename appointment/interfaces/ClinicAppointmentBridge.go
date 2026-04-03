package interfaces

import (
	"context"
	"time"

	"github.com/AlladinDev/AlShifa/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicModule interface {
	ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDateRequested time.Time) (doctorName string, clinicName string, clinicAddress string, clinicMaxAppointments int, error *structs.IAppError)
	ClinicExists(ctx context.Context, clinicID primitive.ObjectID) *structs.IAppError
	DoctorExists(ctx context.Context, doctorID primitive.ObjectID) *structs.IAppError
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, *structs.IAppError)
	GetClinicIDIfExists(ctx context.Context, filters bson.M) (ID primitive.ObjectID, error *structs.IAppError)
	GetClinicIDByReceptionist(ctx context.Context, receptionistID primitive.ObjectID) (clinicID primitive.ObjectID, error *structs.IAppError)
	DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) *structs.IAppError
}
