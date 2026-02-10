package interfaces

import (
	"AlShifa/Clinic/models"
	structs "AlShifa/Structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IService interface contains functions that clinic service layer must implement( to beused by handlers)
type IService interface {
	RegisterClinic(ctx context.Context, ownerID string, clinic models.Clinic) *structs.IAppError
	RegisterClinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, *structs.IAppError)
	LoginClinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError)
	LoginDoctor(ctx context.Context, email string, password string) (string, *structs.IAppError)
	SearchDoctor(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError)
	AddDoctorToClinic(ctx context.Context, clinicDetails models.AddDoctorToClinic) *structs.IAppError
	SearchOwner(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError)
	RegisterDoctor(ctx context.Context, doctor models.Doctor) *structs.IAppError
	VerifyAddDoctorToClinicOTP(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError
	AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (*models.Appointment, *structs.IAppError)
	AppointmentSlotsBooked(ctx context.Context, slotDetails models.SlotDetails) ([]models.Slot, *structs.IAppError)
}
