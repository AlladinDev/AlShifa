package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plan struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	AmountToDeduct float64            `json:"amountToDeduct" bson:"amountToDeduct"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}
