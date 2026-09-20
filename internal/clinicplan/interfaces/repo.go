package planinterface

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/clinicplan/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	FetchAmountToDeduct(ctx context.Context, filter bson.M) (primitive.ObjectID, int, error)
	AddPlan(ctx context.Context, details *models.Plan) error
	FetchPlans(ctx context.Context, filter bson.M) ([]models.Plan, error)
	UpdatePlan(ctx context.Context, ID primitive.ObjectID, data *models.Plan) error
	FetchPlanID(ctx context.Context, planName string) (primitive.ObjectID, error)
}
