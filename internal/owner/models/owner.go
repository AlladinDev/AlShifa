// Package models provides models for owner module
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Owner struct {
	CreatedAt                time.Time          `json:"createdAt" bson:"createdAt"`            // 24 bytes
	Name                     string             `json:"name" Form:"name" bson:"name"`          // 16 bytes
	Address                  string             `json:"address" Form:"address" bson:"address"` // 16 bytes
	Gender                   string             `json:"gender" Form:"gender" bson:"gender"`    // 16 bytes
	ID                       primitive.ObjectID `json:"_id"  bson:"_id"`
	Role                     string             `json:"role"  bson:"role"` // 12 bytes (placed near end)
	Mobile                   string             `json:"mobile" Form:"mobile" bson:"mobile"`
	Password                 string             `json:"password" Form:"password" bson:"password"`
	Email                    string             `json:"email" Form:"email" bson:"email"`
	AadhaarNumber            string             `json:"aadhaarNumber" Form:"aadhaarNumber" bson:"aadhaarNumber"`
	Photo                    string             `json:"photo" Form:"photo" bson:"photo"`
	PhotoID                  string             `json:"photoID"  bson:"photoID"`
	MobileVerified           bool               `json:"mobileVerified"  bson:"mobileVerified"`
	EmailVerified            bool               `json:"emailVerified" bson:"emailVerified"`
	AadhaarVerified          bool               `json:"aadhaarVerified" bson:"aadhaarVerified"`
	AccountOpeningAmountPaid bool               `json:"accountOpeningAmountPaid" bson:"accountOpeningAmountPaid"`
}
