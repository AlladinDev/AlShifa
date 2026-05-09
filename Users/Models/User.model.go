// Package models provides models for user module
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Role      string             `json:"role" bson:"role"`
	Password  string             `json:"password" bson:"password"`
	Address   string             `json:"address" bson:"address"`
	Email     string             `json:"email" bson:"email"`
	Mobile    string             `json:"mobile" bson:"mobile"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
