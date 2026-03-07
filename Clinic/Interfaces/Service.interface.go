package interfaces

import (
	"AlShifa/clinic/models"
	structs "AlShifa/structs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IService interface contains functions that clinic service layer must implement( to beused by handlers)
type IService interface {
	Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) *structs.IAppError
	Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError)
	RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) *structs.IAppError
	ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDateRequested time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string, maxAppointmentsAllowed int)
	AddDoctorToclinic(ctx context.Context, userID primitive.ObjectID, clinicDetails models.ClinicDoctor) *structs.IAppError
	FetchDoctors(ctx context.Context, filter bson.M) ([]models.Doctor, *structs.IAppError)
	VerifyAddDoctorToclinicOTP(ctx context.Context, otp string, userID primitive.ObjectID) *structs.IAppError
	FetchDoctorClinicMappings(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError)
	ClinicExists(ctx context.Context, clinicID primitive.ObjectID) *structs.IAppError
	DoctorExists(ctx context.Context, doctorID primitive.ObjectID) *structs.IAppError
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, *structs.IAppError)
	DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) *structs.IAppError
}
