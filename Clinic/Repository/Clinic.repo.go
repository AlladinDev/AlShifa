// Package repository provides the implementation of the repository layer for managing clinic data in MongoDB.
package repository

import (
	"AlShifa/Clinic/models"
	"context"
	"fmt"

	interfaces "AlShifa/Clinic/Interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repo is the MongoDB implementation of the IRepository interface.
type Repo struct {
	DB *mongo.Database
}

///this ensures this service layer implements all methods of service layer interface
var _ interfaces.IRepository = (*Repo)(nil)

// NewRepository creates a new repository with the specified database and collection name.
func NewRepository(db *mongo.Database) interfaces.IRepository {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) RegisterClinicOwner(ctx context.Context, owner models.Owner) error {
	fmt.Print("inside RegisterClinicOwner")
	_, err := r.DB.Collection("Owner").InsertOne(ctx, owner)
	return err
}

func (r *Repo) RegisterClinic(
	ctx context.Context,
	ownerID primitive.ObjectID,
	clinic models.Clinic,
) error {
	//add ownerId  to clinic field
	clinic.Owner = ownerID

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
			{Key: "localField", Value: "doctors"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "doctorDetails"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Owner"},
			{Key: "localField", Value: "owner"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "ownerDetails"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$ownerDetails"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
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
