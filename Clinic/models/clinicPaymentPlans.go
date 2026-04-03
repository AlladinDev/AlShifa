package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClinicPlan struct {
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Type           string             `json:"type"  bson:"type"`
	AmountToDeduct int                `json:"amountToDeduct" bson:"amountToDeduct"`
}
