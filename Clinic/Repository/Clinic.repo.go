// Package repository provides the implementation of the repository layer for managing clinic data in MongoDB.
package repository

import (
	"AlShifa/Clinic/models"
	"context"
	"errors"
	"fmt"

	interfaces "AlShifa/Clinic/Interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo is the MongoDB implementation of the IRepository interface.
type Repo struct {
	DB *mongo.Database
}

///this ensures this service layer implements all methods of service layer interface
var _ interfaces.IRepository = (*Repo)(nil)

// NewRepository creates a new repository with the specified database and collection name.
func NewRepository(db *mongo.Database) interfaces.IRepository {

	//here call InitialiseIndexes for Slot table
	//because it we call this initialiseIndexes in module initialisation function developer may forget to call it and it will corrupt db in edge cases
	if err := InitialiseIndexes(db.Collection("Slot")); err != nil {
		panic(fmt.Sprintf("Failed to create index for slot collection error is %v :", err))
	}

	fmt.Println("Indexes created successfully for clinic repo")
	return &Repo{
		DB: db,
	}

}

func InitialiseIndexes(collection *mongo.Collection) error {
	//call create index  function for slot  document
	if err := CreateSlotIndex(collection, bson.D{
		{Key: "bookingDate", Value: 1},
		{Key: "doctorID", Value: 1},
		{Key: "clinicID", Value: 1},
	}); err != nil {
		return err
	}

	return nil
}

func (r *Repo) RegisterClinicOwner(ctx context.Context, owner models.Owner) error {
	_, err := r.DB.Collection("Owner").InsertOne(ctx, owner)
	return err
}

func (r *Repo) RegisterClinic(
	ctx context.Context,
	ownerID primitive.ObjectID,
	clinic models.Clinic,
) error {

	session, err := r.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (any, error) {
		// 1️⃣ Insert clinic
		res, err := r.DB.Collection("Clinic").InsertOne(sessCtx, clinic)
		if err != nil {
			return nil, err
		}

		clinicID := res.InsertedID.(primitive.ObjectID)

		// 2️⃣ Update owner with clinic ID
		_, err = r.DB.Collection("Owner").UpdateOne(
			sessCtx,
			bson.M{
				"_id": ownerID,
			},

			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "clinic", Value: clinicID},
				}},
			},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	fmt.Print(err)
	return err
}

func (r *Repo) GetOwnerDetails(ctx context.Context, filter bson.M) ([]models.Owner, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Clinic"},
				{Key: "localField", Value: "clinic"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "clinicDetails"},
			}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$clinicDetails",
			"preserveNullAndEmptyArrays": true,
		}},
		},
	}
	cursor, err := r.DB.Collection("Owner").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var owners []models.Owner

	if err := cursor.All(ctx, &owners); err != nil {
		return nil, err
	}

	return owners, nil
}

func (r *Repo) SearchClinic(ctx context.Context, filter bson.M) ([]models.Clinic, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Doctor"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "clinics.clinic"},
			{Key: "as", Value: "doctorDetails"},
		}}},

		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Owner"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "clinic"},
				{Key: "as", Value: "ownerDetails"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$ownerDetails"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},

		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "ownerDetails.password", Value: 0},
				{Key: "doctorDetails.password", Value: 0},
				{Key: "ownerDetails.gender", Value: 0},
				{Key: "ownerDetails.email", Value: 0},
				{Key: "ownerDetails.registrationDate", Value: 0},
				{Key: "ownerDetails._id", Value: 0},
				{Key: "ownerDetails.clinic", Value: 0},
				{Key: "ownerDetails.role", Value: 0},
				{Key: "doctorDetails.registrationDate", Value: 0},
				{Key: "doctorDetails._id", Value: 0},
				{Key: "doctorDetails.email", Value: 0},
				{Key: "doctorDetails.clinics", Value: 0},
				{Key: "doctorDetails.role", Value: 0},
			}},
		},
	}

	cursor, err := r.DB.Collection("Clinic").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Clinic
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repo) RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error {
	_, err := r.DB.Collection("Doctor").InsertOne(ctx, doctorDetails)
	return err
}

func (r *Repo) SearchDoctors(ctx context.Context, filter bson.M) ([]models.DoctorPublicDetails, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: filter},
		},

		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$clinics"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Clinic"},
			{Key: "localField", Value: "clinics.clinic"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "clinics.information"},
		}}},

		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$clinics.information"},
		}}},

		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$name"}}},
			{Key: "qualifications", Value: bson.D{{Key: "$first", Value: "$qualifications"}}},
			{Key: "workingAt", Value: bson.D{{Key: "$first", Value: "$workingAt"}}},
			{Key: "clinics", Value: bson.D{{Key: "$push", Value: "$clinics"}}},
		}}},
	}

	cursor, err := r.DB.Collection("Doctor").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var doctors []models.DoctorPublicDetails
	if err := cursor.All(ctx, &doctors); err != nil {
		return nil, err
	}
	fmt.Print("doctors are", doctors)

	return doctors, nil
}

//SearchDoctor searches for a single doctor based on the provided filter. it returns password and email fields so use it for internal apis like login only and not for public use
func (r *Repo) SearchDoctor(ctx context.Context, filter bson.M) (models.Doctor, error) {
	result := r.DB.Collection("Doctor").FindOne(ctx, filter)
	var doctor models.Doctor
	err := result.Decode(&doctor)
	return doctor, err
}

func (r *Repo) AddDoctorToClinic(ctx context.Context, clinicDetails models.AddDoctorToClinic) error {

	doctorUpdationRes := r.DB.Collection("Doctor").FindOneAndUpdate(ctx, bson.M{"_id": clinicDetails.DoctorID}, bson.M{
		"$push": bson.M{
			"clinics": models.ClinicDetails{
				StartTime:   clinicDetails.StartTime,
				EndTime:     clinicDetails.EndTime,
				WorkingDays: clinicDetails.WorkingDays,
				Clinic:      clinicDetails.ClinicID,
			},
		},
	})

	if doctorUpdationRes.Err() != nil {
		return doctorUpdationRes.Err()
	}

	return nil
}

func (r *Repo) AddAppointment(ctx context.Context, maxAppointments int, appointmentDetails models.Appointment) (*models.Appointment, error) {
	//here first update slot document using upsert to ensure if it is not present create it if present updates its slots booked by 1
	session, err := r.DB.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	//this is the transaction function which will contain db operations logic which must be executed in single operation
	transactionFn := func(sessCtx mongo.SessionContext) (any, error) {
		slotFilter := bson.D{
			{Key: "bookingDate", Value: appointmentDetails.AppointmentDate},
			{Key: "doctorID", Value: appointmentDetails.Doctor},
			{Key: "clinicID", Value: appointmentDetails.Clinic},
			{Key: "slotsBooked", Value: bson.M{
				"$lt": maxAppointments,
			}},
		}

		updateQuery := bson.D{
			{Key: "$inc", Value: bson.D{
				{Key: "slotsBooked", Value: 1},
			}},
			{Key: "$setOnInsert", Value: bson.D{
				{Key: "bookingDate", Value: appointmentDetails.AppointmentDate},
				{Key: "doctorID", Value: appointmentDetails.Doctor},
				{Key: "clinicID", Value: appointmentDetails.Clinic},
			}},
		}

		//create the update options to update slot first then get its slotsBooked and assign it to current appointment as slot
		updateOptions := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		slotUpdationRes := r.DB.Collection("Slot").FindOneAndUpdate(sessCtx, slotFilter, updateQuery, updateOptions)
		var slotUpdated models.Slot
		if err := slotUpdationRes.Decode(&slotUpdated); err != nil {
			return nil, err
		}

		//here get the updated slotsBooked and set it in appointment details it will represent the slot for this appointment
		appointmentDetails.Slot = slotUpdated.SlotsBooked

		//now save the appointment also
		_, err := r.DB.Collection("Appointment").InsertOne(sessCtx, appointmentDetails)
		if err != nil {
			return nil, err
		}

		return appointmentDetails, nil
	}

	res, transactionErr := session.WithTransaction(ctx, transactionFn)
	if transactionErr != nil {
		return nil, transactionErr
	}

	//now try to convert res into appointment model
	appointmentCreated, ok := res.(models.Appointment)
	if !ok {
		return nil, errors.New("failed to return appointment Created")
	}

	return &appointmentCreated, nil

}

func (r *Repo) AppointmentSlotsBooked(ctx context.Context, maxAppointments int, clinicID primitive.ObjectID, doctorID primitive.ObjectID) ([]models.Slot, error) {
	findQuery := bson.M{
		"clinicID":    clinicID,
		"doctorID":    doctorID,
		"slotsBooked": maxAppointments,
	}
	cur, err := r.DB.Collection("Slot").Find(ctx, findQuery)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var slots []models.Slot
	if err := cur.All(ctx, &slots); err != nil {
		return nil, err
	}

	return slots, nil
}
