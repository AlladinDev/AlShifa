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

func (r *Repo) RegisterDoctor(ctx context.Context, doctorDetails models.Doctor) error {
	_, err := r.DB.Collection("Doctor").InsertOne(ctx, doctorDetails)
	return err
}

func (r *Repo) SearchDoctors(ctx context.Context, filter bson.M) ([]models.Doctor, error) {
	pipeline := mongo.Pipeline{
		// 1️⃣ Match doctors by filter
		bson.D{{Key: "$match", Value: filter}},

		// 2️⃣ Unwind clinics (Keep doctor even if clinics is empty/null)
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$clinics"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// 3️⃣ Lookup clinic info
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Clinic"},
			{Key: "localField", Value: "clinics.clinic"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "clinics.information"},
		}}},

		// 4️⃣ Flatten the information array (Lookup returns an array)
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "clinics.information", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$clinics.information", 0}},
			}},
		}}},

		// 5️⃣ Group back
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "password", Value: bson.D{{Key: "$first", Value: "$password"}}},
			{Key: "registrationDate", Value: bson.D{{Key: "$first", Value: "$registrationDate"}}},
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$name"}}},
			{Key: "qualifications", Value: bson.D{{Key: "$first", Value: "$qualifications"}}},
			{Key: "address", Value: bson.D{{Key: "$first", Value: "$address"}}},
			{Key: "email", Value: bson.D{{Key: "$first", Value: "$email"}}},

			{Key: "workingAt", Value: bson.D{{Key: "$first", Value: "$workingAt"}}},
			{Key: "mobile", Value: bson.D{{Key: "$first", Value: "$mobile"}}},
			{Key: "appointments", Value: bson.D{{Key: "$first", Value: "$appointments"}}},

			// ⭐ CRITICAL FIX: Only push to array if clinics.clinic exists
			{Key: "clinics", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gt", Value: bson.A{"$clinics.clinic", nil}}},
						"$clinics",
						"$$REMOVE",
					}},
				}},
			}},
		}}},

		// 6️⃣ Final Polish: If the array is empty, set the field to null
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "clinics", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{bson.D{{Key: "$size", Value: "$clinics"}}, 0}}},
					nil,
					"$clinics",
				}},
			}},
		}}},
	}

	cursor, err := r.DB.Collection("Doctor").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var doctors []models.Doctor
	if err := cursor.All(ctx, &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}
