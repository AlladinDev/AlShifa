package doctorinterfaces

import (
	"context"

	models "github.com/AlladinDev/AlShifa/internal/doctor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error
	FetchDoctor(ctx context.Context, filter bson.M) (*models.Doctor, error)
	DoctorExists(ctx context.Context, filters bson.M) (bool, error)
	FetchDoctors(ctx context.Context, filters bson.M) ([]models.Doctor, error)
	CheckClinicAllowedToOnboard(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) (bool, error)
	AllowClinicToOnboard(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) error
}
