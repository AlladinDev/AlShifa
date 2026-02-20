// Package interfaces contains interfaces for clinic module
package interfaces

import (
	DTO "AlShifa/clinic/dtos"
	"AlShifa/clinic/models"
	sharedModels "AlShifa/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IRepository defines the methods for loose coupling between the repository and its implementation.
type IRepository interface {
	Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) error
	RegisterclinicOwner(ctx context.Context, owner models.Owner) error
	GetOwnerDetails(ctx context.Context, filter bson.M) ([]models.Owner, error)
	Searchclinic(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error)
	RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error
	SearchDoctors(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, error)
	SearchDoctor(ctx context.Context, filter bson.M) (models.Doctor, error)
	AddDoctorToclinic(ctx context.Context, clinicDetails models.AddDoctorToclinic) error
	AddAppointment(ctx context.Context, maxAppointments int, appointmentDetails models.Appointment) (int, error)
	AppointmentSlotsBooked(ctx context.Context, maxAppointments int, clinicID primitive.ObjectID, doctorID primitive.ObjectID) ([]models.Slot, error)
	SearchclinicByID(ctx context.Context, clinicID primitive.ObjectID) (*models.Clinic, error)
	FetchDoctorAtclinics(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, error)
	FetchAppointments(ctx context.Context, groupBy string, filter bson.M) ([]sharedModels.Appointments, error)
}
