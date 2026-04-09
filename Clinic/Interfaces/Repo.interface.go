// Package interfaces contains interfaces for clinic module
package interfaces

import (
	"context"
	"time"

	"github.com/AlladinDev/AlShifa/clinic/dtos"
	"github.com/AlladinDev/AlShifa/clinic/models"
	structs "github.com/AlladinDev/AlShifa/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IRepository defines the methods for loose coupling between the repository and its implementation.
type IRepository interface {
	Registerclinic(ctx context.Context, ownerID primitive.ObjectID, clinic models.Clinic) error
	Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, error)
	RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error
	AddDoctorToclinic(ctx context.Context, clinicDetails models.ClinicDoctor) error
	SearchclinicByID(ctx context.Context, clinicID primitive.ObjectID) (*models.Clinic, error)
	ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDate time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string)
	FetchDoctors(ctx context.Context, filter bson.M) ([]models.Doctor, error)
	ClinicWithDoctors(ctx context.Context, filter bson.M) ([]dtos.ClinicWithDoctors, error)
	DoctorWithClinics(ctx context.Context, filter bson.M) ([]dtos.DoctorWithClinics, error)
	FetchDoctorClinicMappings(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error)
	ClinicExists(ctx context.Context, clinicID primitive.ObjectID) error
	DoctorExists(ctx context.Context, doctorID primitive.ObjectID) error
	FetchDoctorProfile(ctx context.Context, filter bson.M) (*models.Doctor, error)
	GetPlanDetails(ctx context.Context, planID primitive.ObjectID) (*models.ClinicPlan, error)
	FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error)
	DeductClinicWallet(ctx context.Context, amountToDeduct int, clinicID primitive.ObjectID) error
	GetClinicIDIfExists(ctx context.Context, filters bson.M) (ID primitive.ObjectID, err error)
	GetClinicIDByReceptionist(ctx context.Context, receptionistID primitive.ObjectID) (clinicID primitive.ObjectID, err error)
	DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) error
	FetchPlanDetails(ctx context.Context, filter bson.M) (*models.ClinicPlan, error)
}
