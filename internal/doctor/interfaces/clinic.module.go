package doctorinterfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type IClinic interface {
	ClinicExists(ctx context.Context, filter bson.M) (bool, error)
}
