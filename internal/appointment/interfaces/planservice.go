package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPlanService interface {
	FetchAmountToDeduct(ctx context.Context, filter bson.M) (planID primitive.ObjectID, amountToDeduct int, err error)
}
