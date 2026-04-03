// Package repository provides the implementation of the repository layer for managing clinic data in MongoDB.
package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AlladinDev/AlShifa/clinic/models"
	"github.com/AlladinDev/AlShifa/constants"
	"github.com/AlladinDev/AlShifa/structs"

	interfaces "github.com/AlladinDev/AlShifa/clinic/interfaces"

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

func (r *Repo) FetchDoctorProfile(ctx context.Context, filter bson.M) (*models.Doctor, error) {
	res := r.DB.Collection("Doctor").FindOne(ctx, filter)

	var doctor models.Doctor
	if err := res.Decode(&doctor); err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *Repo) GetClinicIDIfExists(ctx context.Context, filters bson.M) (ID primitive.ObjectID, err error) {
	options := options.FindOne().SetProjection(bson.M{"_id": 1})
	res := r.DB.Collection("Clinic").FindOne(ctx, filters, options)
	type Result struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	var dbRes Result
	if err := res.Decode(&dbRes); err != nil {
		return primitive.NilObjectID, nil
	}

	return dbRes.ID, nil
}

func (r *Repo) SearchclinicByID(ctx context.Context, clinicID primitive.ObjectID) (*models.Clinic, error) {
	res := r.DB.Collection("Clinic").FindOne(ctx, bson.M{"_id": clinicID})
	var clinic models.Clinic
	if err := res.Decode(&clinic); err != nil {
		return nil, err
	}
	return &clinic, nil
}

func (r *Repo) GetClinicIDByReceptionist(ctx context.Context, receptionistID primitive.ObjectID) (clinicID primitive.ObjectID, err error) {
	type Result struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	var dbRes Result
	res := r.DB.Collection("Clinic").FindOne(ctx, bson.M{"receptionistID": receptionistID})
	if err := res.Decode(&dbRes); err != nil {
		return primitive.NilObjectID, err
	}

	return dbRes.ID, nil
}

func (r *Repo) Registerclinic(
	ctx context.Context,
	ownerID primitive.ObjectID,
	clinic models.Clinic,
) error {

	//find the basic plan for clinics
	basicPlanRes := r.DB.Collection("ClinicPlan").FindOne(ctx, bson.M{"type": constants.ClinicPlanBasic})
	var plan models.ClinicPlan
	if err := basicPlanRes.Decode(&plan); err != nil {
		return err
	}

	//add the plan id to clinic payload to attach plan to it
	clinic.PlanID = plan.ID

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

func (r *Repo) GetPlanDetails(ctx context.Context, planID primitive.ObjectID) (*models.ClinicPlan, error) {
	res := r.DB.Collection("ClinicPlan").FindOne(ctx, bson.M{"_id": planID})
	var planDetails models.ClinicPlan
	if err := res.Decode(&planDetails); err != nil {
		return nil, err
	}

	return &planDetails, nil
}

func (r *Repo) Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, error) {

	cursor, err := r.DB.Collection("Clinic").Find(ctx, filter)
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

func (r *Repo) FetchDoctors(ctx context.Context, filter bson.M) ([]models.Doctor, error) {

	res, err := r.DB.Collection("Doctor").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var doctors []models.Doctor
	if err := res.All(ctx, &doctors); err != nil {
		return nil, err
	}
	return doctors, nil
}

func (r *Repo) AddDoctorToclinic(ctx context.Context, clinicDetails models.ClinicDoctor) error {
	_, err := r.DB.Collection("ClinicDoctor").InsertOne(ctx, clinicDetails)
	return err
}

func (r *Repo) DeductClinicWallet(ctx context.Context, amountToDeduct int, clinicID primitive.ObjectID) error {
	filter := bson.M{"_id": clinicID, "wallet.availableBalance": bson.M{"$gte": amountToDeduct}}
	res, err := r.DB.Collection("Clinic").UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"wallet.availableBalance": -amountToDeduct}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed to deduct clinic wallet")
	}

	return nil
}

// FetchDoctorClinicMappings function clinic with its associated doctors using filters
func (r *Repo) FetchDoctorClinicMappings(ctx context.Context, filter bson.M) ([]models.ClinicDoctor, error) {

	cursor, err := r.DB.Collection("ClinicDoctor").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var doctorClinicMapping []models.ClinicDoctor
	if err := cursor.All(ctx, &doctorClinicMapping); err != nil {
		return nil, err
	}

	return doctorClinicMapping, nil

}

//ClinicDoctorDetails function is for getting some details from clinic doctor mapping corresponding to clinicid and doctorid it returns some details but in primitive individual format
func (r *Repo) ClinicDoctorDetails(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID, appointmentDate time.Time) (error *structs.IAppError, doctorName string, clinicName string, clinicAddress string) {
	//first check whether this clinic exists or not
	clinicCheckingErr := r.DB.Collection("Clinic").FindOne(ctx, bson.M{"_id": clinicID}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if clinicCheckingErr != nil {
		if clinicCheckingErr == mongo.ErrNoDocuments {
			return &structs.IAppError{
				Message:    "This Clinic Doesnt Exist",
				Reason:     clinicCheckingErr.Error(),
				StatusCode: http.StatusNotFound,
			}, "", "", ""
		}

		return &structs.IAppError{
			Message:    "Failed to check clinic existence",
			Reason:     clinicCheckingErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}, "", "", ""
	}

	//now check if this doctor exists or not
	doctorCheckingErr := r.DB.Collection("Doctor").FindOne(ctx, bson.M{"_id": clinicID}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if doctorCheckingErr != nil {
		if doctorCheckingErr == mongo.ErrNoDocuments {
			return &structs.IAppError{
				Message:    "This Doctor Doesnt Exist",
				Reason:     doctorCheckingErr.Error(),
				StatusCode: http.StatusNotFound,
			}, "", "", ""
		}

		return &structs.IAppError{
			Message:    "Failed to check clinic existence",
			Reason:     doctorCheckingErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}, "", "", ""
	}

	result := r.DB.Collection("ClinicDoctor").FindOne(ctx, bson.M{"clinicID": clinicID, "doctorID": doctorID})
	var doctorClinicDetails models.ClinicDoctor
	if err := result.Decode(&doctorClinicDetails); err != nil {
		return &structs.IAppError{
			Message:    "Failed to get doctor clinic details",
			Reason:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, "", "", ""
	}

	//now here check whether this appointmentDate actually is valid or not whether doctor is available on this date or not
	requestedDateDay := appointmentDate.Weekday()
	appointmentDatePossible := false
	for _, day := range doctorClinicDetails.WorkingDays {
		if strings.EqualFold(day, requestedDateDay.String()) {
			appointmentDatePossible = true
			break
		}
	}

	if !appointmentDatePossible {
		return &structs.IAppError{
			Message:    "Doctor is not available on this date",
			Reason:     "Doctor is not available on this date",
			StatusCode: http.StatusBadRequest,
		}, "", "", ""
	}

	return nil, doctorClinicDetails.DoctorName, doctorClinicDetails.ClinicName, doctorClinicDetails.ClinicAddress
}

func (r *Repo) RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error {
	_, err := r.DB.Collection("Doctor").InsertOne(ctx, doctorDetails)
	return err
}

func (r *Repo) ClinicExists(ctx context.Context, clinicID primitive.ObjectID) error {
	return r.DB.Collection("Clinic").FindOne(ctx, bson.M{"_id": clinicID}).Err()
}

func (r *Repo) DoctorExists(ctx context.Context, doctorID primitive.ObjectID) error {
	return r.DB.Collection("Doctor").FindOne(ctx, bson.M{"_id": doctorID}).Err()
}

func (r *Repo) DoctorClinicMappingExists(ctx context.Context, clinicID primitive.ObjectID, doctorID primitive.ObjectID) error {
	return r.DB.Collection("ClinicDoctor").FindOne(ctx, bson.M{"doctorID": doctorID, "clinicID": clinicID}).Err()
}

func (r *Repo) FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error) {
	res := r.DB.Collection("Clinic").FindOne(ctx, bson.M{"_id": clinicID})
	var clinic models.Clinic
	if err := res.Decode(&clinic); err != nil {
		return 0, err
	}

	return clinic.MaxAppointments, nil
}

func (r *Repo) FetchPlanDetails(ctx context.Context, filter bson.M) (*models.ClinicPlan, error) {
	res := r.DB.Collection("ClinicPlan").FindOne(ctx, filter)
	var planDetails models.ClinicPlan
	if err := res.Decode(&planDetails); err != nil {
		return nil, err
	}

	return &planDetails, nil
}
