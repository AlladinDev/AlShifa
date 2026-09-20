// Package doctorrepository provides repository functions for doctor module
package doctorrepository

import (
	"context"
	"errors"

	doctorinterfaces "github.com/AlladinDev/AlShifa/internal/doctor/interfaces"
	models "github.com/AlladinDev/AlShifa/internal/doctor/models"
	"github.com/AlladinDev/AlShifa/internal/shared/customerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo is the MongoDB implementation of the IRepository interface.
type Repo struct {
	db *mongo.Database
}

///this ensures this service layer implements all methods of service layer interface
var _ doctorinterfaces.IRepository = (*Repo)(nil)

// NewRepository creates a new repository with the specified database and collection name.
func NewRepository(DB *mongo.Database) doctorinterfaces.IRepository {
	return &Repo{
		db: DB,
	}
}

func (r *Repo) DoctorExists(ctx context.Context, filters bson.M) (bool, error) {
	if err := r.db.Collection("Doctor").FindOne(ctx, filters, options.FindOne().SetProjection(bson.M{"_id": 1})).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repo) RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error {
	_, err := r.db.Collection("Doctor").InsertOne(ctx, doctorDetails)
	return err
}

func (r *Repo) FetchDoctor(ctx context.Context, filter bson.M) (*models.Doctor, error) {
	res := r.db.Collection("Doctor").FindOne(ctx, filter)
	var doctor models.Doctor
	if err := res.Decode(&doctor); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *Repo) FetchDoctors(ctx context.Context, filters bson.M) ([]models.Doctor, error) {
	cur, err := r.db.Collection("Doctor").Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var doctors []models.Doctor
	if err := cur.All(ctx, &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *Repo) CheckClinicAllowedToOnboard(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) (bool, error) {
	if err := r.db.Collection("Doctor").FindOne(
		ctx,
		bson.M{
			"_id":                     doctorID,
			"clinicsAllowedToOnboard": clinicID,
		},
	).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repo) AllowClinicToOnboard(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) error {
	res, err := r.db.Collection("Doctor").UpdateOne(ctx, bson.M{"_id": doctorID}, bson.M{"$addToSet": bson.M{"clinicsAllowedToOnboard": clinicID}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	if res.ModifiedCount == 0 {
		return customerrors.ErrNoDocumentModified
	}

	return nil
}
