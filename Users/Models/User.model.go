// Package models provides models for user module
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Role    string             `json:"role" bson:"role"`
	Address string             `json:"address" bson:"address"`
}
