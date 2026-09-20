// Package planinterface provides interfaces for plan module
package planinterface

import (
	"context"

	"github.com/AlladinDev/AlShifa/internal/clinicplan/models"
	"github.com/AlladinDev/AlShifa/internal/shared/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IService interface {
	RegisterPlan(ctx context.Context, details *models.Plan) *response.IAppError
	FetchAmountToDeduct(ctx context.Context, filter bson.M) (planID primitive.ObjectID, amountToDeduct int, err error)
	FetchPlanID(ctx context.Context, planName string) (primitive.ObjectID, error)
	FetchPlans(ctx context.Context, filters bson.M) ([]models.Plan, *response.IAppError)
	UpdatePlan(ctx context.Context, ID primitive.ObjectID, data *models.Plan) *response.IAppError
}
