package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	DB *mongo.Database
}

func NewRepo(database *mongo.Database) interfaces.IRepo {
	return &Repo{
		DB: database,
	}
}

func (r *Repo) RegisterDoctorClinicMapping(ctx context.Context, details *models.ClinicDoctorMapping) error {
	_, err := r.DB.Collection("ClinicDoctorMapping").InsertOne(ctx, details)
	return err
}

//CheckClinicDoctorMappingExists function checks whether a clinic and doctor mapping exists if it doesnt return any error it means mapping exists if error is there it means otherwise like mapping dosnt exist or other db error
func (r *Repo) CheckClinicDoctorMappingExists(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error) {
	if err := r.DB.Collection("ClinicDoctor").FindOne(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID}).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repo) CheckAppointmentDatePossible(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDate time.Time) (bool, error) {
	var mapping models.ClinicDoctorMapping
	if err := r.DB.Collection("ClinicDoctor").FindOne(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID}).Decode(&mapping); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, err
		}
		return false, err
	}

	//now here it means mapping does exist but now we have to check whether this doctor mapping has availability for this appointment date passed or not
	for _, availableDay := range mapping.AvailableOn {
		if !strings.EqualFold(availableDay, appointmentDate.Weekday().String()) {
			return false, nil
		}
	}

	return true, nil
}
