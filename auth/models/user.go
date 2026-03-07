// Package models provides models for auth module
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credientials struct {
	ID               primitive.ObjectID `json:"userID" bson:"userID"`
	Email            string             `json:"email" bson:"email"`
	Password         string             `json:"password" bson:"password"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	Role             string             `json:"role" bson:"role"`
	Mobile           string             `json:"mobile" bson:"mobile"`
	ProfileCompleted bool               `json:"profileCompleted" bson:"profileCompleted"`
	IsVerified       bool               `json:"isVerified" bson:"isverified"`
}
