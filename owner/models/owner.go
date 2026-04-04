// Package models provides models for owner module
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Owner struct {
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"` // 24 bytes
	Name      string             `json:"name" bson:"name"`           // 16 bytes
	Address   string             `json:"address" bson:"address"`     // 16 bytes
	Gender    string             `json:"gender" bson:"gender"`       // 16 bytes
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Role      string             `json:"role" bson:"role"` // 12 bytes (placed near end)
	Mobile    string              `json:"mobile" bson:"mobile"`
	Password  string             `json:"password" bson:"password"`
	Email     string             `json:"email" bson:"email"`
}
