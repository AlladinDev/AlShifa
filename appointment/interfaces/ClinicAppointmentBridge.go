package interfaces

import (
	"AlShifa/structs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicModule interface {
	ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDateRequested time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string, clinicMaxAppointments int)
	ClinicExists(ctx context.Context, clinicID primitive.ObjectID) *structs.IAppError
	DoctorExists(ctx context.Context, doctorID primitive.ObjectID) *structs.IAppError
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, *structs.IAppError)
	DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) *structs.IAppError
}
