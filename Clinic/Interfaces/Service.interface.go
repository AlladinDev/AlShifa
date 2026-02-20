package interfaces

import (
	DTO "AlShifa/clinic/dtos"
	"AlShifa/clinic/models"
	sharedModels "AlShifa/models"
	structs "AlShifa/structs"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IService interface contains functions that clinic service layer must implement( to beused by handlers)
type IService interface {
	Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) *structs.IAppError
	RegisterclinicOwner(ctx context.Context, ownerDetails models.Owner) *structs.IAppError
	Searchclinic(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, *structs.IAppError)
	LoginclinicOwner(ctx context.Context, email string, password string) (string, *structs.IAppError)
	LoginDoctor(ctx context.Context, email string, password string) (string, *structs.IAppError)
	SearchDoctor(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, *structs.IAppError)
	AddDoctorToclinic(ctx context.Context, clinicDetails models.AddDoctorToclinic) *structs.IAppError
	SearchOwner(ctx context.Context, filter bson.M) ([]models.Owner, *structs.IAppError)
	RegisterDoctor(ctx context.Context, doctor models.Doctor) *structs.IAppError
	VerifyAddDoctorToclinicOTP(ctx context.Context, otp string, doctorID primitive.ObjectID, clinicID primitive.ObjectID) *structs.IAppError
	AddAppointment(ctx context.Context, appointmentDetails models.Appointment) (int, *structs.IAppError)
	AppointmentSlotsBooked(ctx context.Context, slotDetails models.SlotDetails) ([]models.Slot, *structs.IAppError)
	FetchAppointments(ctx context.Context, groupBy string, filter bson.M) ([]sharedModels.Appointments, *structs.IAppError)
	DoctorWithItsclinics(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, *structs.IAppError)
}
