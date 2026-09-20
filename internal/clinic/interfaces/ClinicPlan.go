package clinicinterfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IClinicPlan interface {
	FetchAmountToDeduct(ctx context.Context, filter bson.M) (primitive.ObjectID, int, error)
	FetchPlanID(ctx context.Context, planName string) (primitive.ObjectID, error)
}
