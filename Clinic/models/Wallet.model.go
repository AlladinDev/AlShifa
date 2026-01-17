package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletDetails struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
	AvailableBalance  int64              `json:"availableBalance" bson:"availableBalance"`
	Clinic            primitive.ObjectID `json:"clinic" bson:"clinic"`
	LatestTransaction primitive.ObjectID `json:"latestTransaction" bson:"latestTransaction"`
}
