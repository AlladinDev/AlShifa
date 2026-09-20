package repository

import (
	"context"
	"errors"

	planinterface "github.com/AlladinDev/AlShifa/internal/clinicplan/interfaces"
	"github.com/AlladinDev/AlShifa/internal/clinicplan/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(DB *mongo.Database) planinterface.IRepository {
	return &Repository{
		db: DB,
	}
}

func (r *Repository) AddPlan(ctx context.Context, details *models.Plan) error {
	_, err := r.db.Collection("Plan").InsertOne(ctx, details)
	return err
}

func (r *Repository) FetchAmountToDeduct(ctx context.Context, filter bson.M) (primitive.ObjectID, int, error) {
	var result bson.M
	if err := r.db.Collection("Plan").FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"amountToDeduct": 1, "_id": 1})).Decode(&result); err != nil {
		return primitive.NilObjectID, 0, err
	}

	amountToDeduct, ok := result["amountToDeduct"].(int)
	if !ok {
		return primitive.NilObjectID, 0, errors.New("failed to get amount to deduct from plan details")
	}

	planID, ok := result["_id"].(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, 0, errors.New("failed to get planid from plan details")
	}

	return planID, amountToDeduct, nil
}

func (r *Repository) FetchPlanID(ctx context.Context, planName string) (primitive.ObjectID, error) {
	var res bson.M
	if err := r.db.Collection("Plan").FindOne(ctx, bson.M{"name": planName}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&res); err != nil {
		return primitive.NilObjectID, err
	}

	planID, ok := res["_id"].(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("plan id retrieved is not of primitive.object type")
	}

	return planID, nil
}

func (r *Repository) FetchPlans(ctx context.Context, filter bson.M) ([]models.Plan, error) {
	cur, err := r.db.Collection("Plan").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var plans []models.Plan
	if err := cur.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

func (r *Repository) UpdatePlan(ctx context.Context, ID primitive.ObjectID, data *models.Plan) error {
	res, err := r.db.Collection("Plan").UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": data})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("no document found with match")
	}

	if res.ModifiedCount == 0 {
		return errors.New("no document updated ")
	}

	return nil
}
