package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type IDoctorModule interface {
	DoctorExists(ctx context.Context, filter bson.M) (bool, error)
}
