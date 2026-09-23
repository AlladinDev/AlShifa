package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AlladinDev/AlShifa/internal/clinicdoctor/dtos"
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

// CheckClinicDoctorMappingExists function checks whether a clinic and doctor mapping exists if it doesnt return any error it means mapping exists if error is there it means otherwise like mapping dosnt exist or other db error
func (r *Repo) CheckClinicDoctorMappingExists(ctx context.Context, doctorID primitive.ObjectID, clinicID primitive.ObjectID) (bool, error) {
	if err := r.DB.Collection("ClinicDoctorMapping").FindOne(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID}).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repo) CheckAppointmentDatePossible(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDate time.Time) (bool, error) {
	var mapping models.ClinicDoctorMapping
	if err := r.DB.Collection("ClinicDoctorMapping").FindOne(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID}).Decode(&mapping); err != nil {
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

func (r *Repo) FetchClinicWithDoctors(ctx context.Context, filter bson.M) ([]dtos.ClinicDoctorDTO, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Doctor"},
			{Key: "localField", Value: "doctorID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "doctor"},
		}}},

		bson.D{{Key: "$unwind", Value: "$doctor"}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Clinic"},
			{Key: "localField", Value: "clinicID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "clinic"},
		}}},

		bson.D{{Key: "$unwind", Value: "$clinic"}},

		bson.D{{Key: "$match", Value: filter}},

		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$clinic._id"},

			{Key: "mappingId", Value: bson.D{
				{Key: "$first", Value: "$_id"},
			}},

			{Key: "clinic", Value: bson.D{
				{Key: "$first", Value: "$clinic"},
			}},

			{Key: "doctors", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "mappingID", Value: "$_id"},
					{Key: "_id", Value: "$doctor._id"},
					{Key: "name", Value: "$doctor.name"},
					{Key: "experience", Value: "$doctor.experience"},
					{Key: "joinedOn", Value: "$createdAt"},
					{Key: "photoUrl", Value: "$doctor.profilePhoto"},
					{Key: "availableOn", Value: "$availableOn"},
					{Key: "timings", Value: "$timings"},
					{Key: "consultationFees", Value: "$consultationFees"},
					{Key: "speciality", Value: "$doctor.field"},
					{Key: "qualifications", Value: "$doctor.qualifications"},
					{Key: "post", Value: "$doctor.post"},
					{Key: "workingAt", Value: "$doctor.workingAt"},
				}},
			}},
		}}},

		bson.D{{Key: "$replaceWith", Value: bson.D{
			{Key: "$mergeObjects", Value: bson.A{
				"$clinic",
				bson.D{
					{Key: "mappingId", Value: "$mappingId"},
					{Key: "doctors", Value: "$doctors"},
				},
			}},
		}}},
	}

	cur, err := r.DB.Collection("ClinicDoctorMapping").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var clinicWithDoctors []dtos.ClinicDoctorDTO
	fmt.Printf("%v", clinicWithDoctors)

	if err := cur.All(ctx, &clinicWithDoctors); err != nil {
		return nil, err
	}

	return clinicWithDoctors, nil
}

func (r *Repo) FetchDoctorWithClinics(ctx context.Context, filter bson.M) ([]dtos.DoctorWithClinic, error) {
	if filter == nil {
		filter = bson.M{}
	}
	pipeline := mongo.Pipeline{
		// 1. Lookup Doctor
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Doctor"},
				{Key: "localField", Value: "doctorID"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "doctor"},
			}},
		},

		// 2. Unwind Doctor
		bson.D{
			{Key: "$unwind", Value: "$doctor"},
		},

		// 3. Lookup Clinic
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Clinic"},
				{Key: "localField", Value: "clinicID"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "clinic"},
			}},
		},

		// 4. Unwind Clinic
		bson.D{
			{Key: "$unwind", Value: "$clinic"},
		},

		// 5. Match Filter
		bson.D{
			{Key: "$match", Value: filter},
		},

		// 6. Group Stage
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$doctor._id"},
				{Key: "doctor", Value: bson.D{
					{Key: "$first", Value: "$doctor"},
				}},
				{Key: "clinics", Value: bson.D{
					{Key: "$push", Value: bson.D{
						{Key: "_id", Value: "$clinic._id"},
						{Key: "mappingID", Value: "$_id"},
						{Key: "name", Value: "$clinic.name"},
						{Key: "address", Value: "$clinic.address"},
						{Key: "consultationFees", Value: "$consultationFees"},
						{Key: "availableOn", Value: "$availableOn"},
						{Key: "doctorTimings", Value: "$timings"},
						{Key: "workingDays", Value: "$clinic.workingDays"},
						{Key: "seasonTimings", Value: "$clinic.seasonTimings"},
						{Key: "departments", Value: "$clinic.departments"},
						{Key: "joinedAt", Value: "$createdAt"},
					}},
				}},
			}},
		},

		// 7. Replace Root Stage
		bson.D{
			{Key: "$replaceRoot", Value: bson.D{
				{Key: "newRoot", Value: bson.D{
					{Key: "$mergeObjects", Value: bson.A{
						"$doctor",
						bson.D{
							{Key: "_id", Value: "$_id"},
							{Key: "clinics", Value: "$clinics"},
						},
					}},
				}},
			}},
		},
	}

	cur, err := r.DB.Collection("ClinicDoctorMapping").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var doctorWithClinics []dtos.DoctorWithClinic
	if err := cur.All(ctx, &doctorWithClinics); err != nil {
		return nil, err
	}

	return doctorWithClinics, nil
}
