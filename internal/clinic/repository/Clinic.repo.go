// Package repository provides the implementation of the repository layer for managing clinic data in MongoDB.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlladinDev/AlShifa/internal/clinic/models"

	interfaces "github.com/AlladinDev/AlShifa/internal/clinic/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo is the MongoDB implementation of the IRepository interface.
type Repo struct {
	DB *mongo.Database
}

// /this ensures this service layer implements all methods of service layer interface
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

func (r *Repo) Searchclinic(ctx context.Context, filter bson.M) ([]models.Clinic, error) {

	fmt.Printf("%+v", filter)
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

func (r *Repo) ClinicExists(ctx context.Context, filter bson.M) (bool, error) {
	if err := r.DB.Collection("Clinic").FindOne(ctx, filter).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repo) FetchMaxAppointments(ctx context.Context, clinicID primitive.ObjectID) (int, error) {
	res := r.DB.Collection("Clinic").FindOne(ctx, bson.M{"_id": clinicID})
	var clinic models.Clinic
	if err := res.Decode(&clinic); err != nil {
		return 0, err
	}
	return clinic.MaxAppointments, nil
}

func (r *Repo) FetchSingleClinic(ctx context.Context, filter bson.M) (*models.Clinic, error) {
	clinic := models.Clinic{}
	res := r.DB.Collection("Clinic").FindOne(ctx, filter)

	if err := res.Decode(&clinic); err != nil {
		return nil, err
	}

	return &clinic, nil
}
