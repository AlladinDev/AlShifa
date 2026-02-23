// Package repository provides the implementation of the repository layer for managing clinic data in MongoDB.
package repository

import (
	DTO "AlShifa/clinic/dtos"
	"AlShifa/clinic/models"
	sharedModels "AlShifa/models"
	"context"
	"errors"
	"fmt"

	interfaces "AlShifa/clinic/interfaces"

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

func (r *Repo) RegisterclinicOwner(ctx context.Context, owner models.Owner) error {
	_, err := r.DB.Collection("Owner").InsertOne(ctx, owner)
	return err
}

func (r *Repo) SearchclinicByID(ctx context.Context, clinicID primitive.ObjectID) (*models.Clinic, error) {
	res := r.DB.Collection("clinic").FindOne(ctx, bson.M{"_id": clinicID})
	var clinic models.Clinic
	if err := res.Decode(&clinic); err != nil {
		return nil, err
	}
	return &clinic, nil
}

func (r *Repo) Registerclinic(
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
		res, err := r.DB.Collection("clinic").InsertOne(sessCtx, clinic)
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
				{Key: "from", Value: "clinic"},
				{Key: "localField", Value: "clinic"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "clinicDetails"},
			}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$clinicDetails",
			"preserveNullAndEmptyArrays": true,
		}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "registrationDate", Value: 0},
				{Key: "clinicDetails.registrationDate", Value: 0},
				{Key: "clinicDetails.ownerDetails", Value: 0},
				{Key: "clinicDetails.doctorDetails", Value: 0},
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

func (r *Repo) Searchclinic(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Doctor"},
			{Key: "localField", Value: "doctorID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "doctorDetails"},
		}}},

		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$doctorDetails"},
			}}},

		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Clinic"},
				{Key: "localField", Value: "clinicID"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "clinicDetails"},
			}},
		},

		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$clinicDetails"},
			}},
		},

		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "doctorDetails.password", Value: 0},
			}},
		},

		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$clinicID"},
				{Key: "clinicDetails", Value: bson.D{
					{Key: "$first", Value: "$clinicDetails"},
				}},
				{Key: "doctors", Value: bson.D{
					{Key: "$push", Value: bson.D{
						{Key: "$mergeObjects", Value: bson.A{
							"$doctorDetails",
							bson.D{{Key: "startTime", Value: "$startTime"}},
							bson.D{{Key: "endTime", Value: "$endTime"}},
							bson.D{{Key: "workingDays", Value: "$workingDays"}},
						}},
					}},
				}},
			}},
		},
	}

	cursor, err := r.DB.Collection("ClinicDoctor").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ClinicDoctor
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

	res, err := r.DB.Collection("Doctor").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var doctors []models.DoctorPublicDetails
	if err := res.All(ctx, &doctors); err != nil {
		return nil, err
	}
	return doctors, nil
}

//SearchDoctor searches for a single doctor based on the provided filter. it returns password and email fields so use it for internal apis like login only and not for public use
func (r *Repo) SearchDoctor(ctx context.Context, filter bson.M) (models.Doctor, error) {
	result := r.DB.Collection("Doctor").FindOne(ctx, filter)
	var doctor models.Doctor
	err := result.Decode(&doctor)
	return doctor, err
}

func (r *Repo) AddDoctorToclinic(ctx context.Context, clinicDetails models.AddDoctorToclinic) error {
	_, err := r.DB.Collection("ClinicDoctor").InsertOne(ctx, clinicDetails)
	return err
}

func (r *Repo) AddAppointment(ctx context.Context, maxAppointments int, appointmentDetails models.Appointment) (int, error) {
	//here first update slot document using upsert to ensure if it is not present create it if present updates its slots booked by 1
	session, err := r.DB.Client().StartSession()
	if err != nil {
		return 0, err
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
		return 0, transactionErr
	}

	//now try to convert res into appointment model
	appointmentCreated, ok := res.(models.Appointment)
	if !ok {
		return 0, errors.New("failed to return appointment Created")
	}

	return appointmentCreated.Slot, nil

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

// FetchDoctorWithclinics function fetches a doctor using doctor name or doctorid along with its clinics where he/she works
func (r *Repo) FetchDoctorAtclinics(ctx context.Context, filter bson.M) ([]DTO.DoctorAtclinicsDTO, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Doctor"},
			{Key: "foreignField", Value: "_id"},
			{Key: "localField", Value: "doctorID"},
			{Key: "as", Value: "doctorDetails"},
		},
		}},

		bson.D{
			{Key: "$unwind", Value: "$doctorDetails"},
		},

		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "clinic"},
				{Key: "foreignField", Value: "_id"},
				{Key: "localField", Value: "clinicID"},
				{Key: "as", Value: "clinicDetails"},
			}},
		},

		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$clinicDetails"},
			}},
		},

		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$doctorID"},
				{Key: "doctorDetails", Value: bson.D{{Key: "$first", Value: "$doctorDetails"}}},
				{Key: "clinics", Value: bson.D{
					{Key: "$push", Value: bson.D{
						{Key: "$mergeObjects", Value: bson.A{
							"$clinicDetails",
							bson.D{{Key: "workingDays", Value: "$workingDays"}},
							bson.D{{Key: "startTime", Value: "$startTiming"}},
							bson.D{{Key: "endTime", Value: "$endTime"}},
						}},
					}},
				}},
			}},
		},

		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "doctorDetails.password", Value: 0},
				{Key: "doctorDetails.email", Value: 0},
			}},
		},
	}

	cursor, err := r.DB.Collection("clinicDoctor").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var doctorAtclinics []DTO.DoctorAtclinicsDTO
	if err := cursor.All(ctx, &doctorAtclinics); err != nil {
		return nil, err
	}

	return doctorAtclinics, nil

}

func (r *Repo) FetchAppointments(ctx context.Context, groupingID string, filter bson.M) ([]sharedModels.Appointments, error) {

	//this is by which documents will be grouped in mongo pipeline
	groupBy := "$" + groupingID

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: groupBy},
			{Key: "totalAppointments", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "appointments", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "doctorName", Value: "$doctorName"},
					{Key: "patientName", Value: "$patientName"},
					{Key: "patientMobile", Value: "$patientMobile"},
					{Key: "userName", Value: "$userName"},
					{Key: "patientAddress", Value: "$patientAddress"},
					{Key: "status", Value: "$status"},
					{Key: "registrationDate", Value: "$registrationDate"},
					{Key: "appointmentDate", Value: "$appointmentDate"},
					{Key: "clinicName", Value: "$clinicName"},
					{Key: "slot", Value: "$slot"},
				}}}},
		}}},
	}

	cursor, err := r.DB.Collection("Appointment").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var appointments []sharedModels.Appointments
	if err := cursor.All(ctx, &appointments); err != nil {
		return nil, err
	}

	return appointments, nil

}
