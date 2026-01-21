// Package models stores database models for clinic
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model OwnerListResponse

// Owner represents the owner of a clinic, including personal and clinic details.
type Owner struct {
	RegistrationDate time.Time          `json:"registrationDate" bson:"registrationDate"` // 24 bytes
	Name             string             `json:"name" bson:"name"`                         // 16 bytes
	Address          string             `json:"address" bson:"address"`                   // 16 bytes
	Password         string             `json:"password" bson:"password"`                 // 16 bytes
	Email            string             `json:"email" bson:"email"`                       // 16 bytes
	Gender           string             `json:"gender" bson:"gender"`                     // 16 bytes
	Clinic           primitive.ObjectID `json:"clinic" bson:"clinic"`                     // 8 bytes (pointer)
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Role             string             `json:"role" bson:"role"` // 12 bytes (placed near end)
	Mobile           int64              `json:"mobile" bson:"mobile"`
	ClinicDetails    *Clinic            `json:"clinicDetails" bson:"clinicDetails,omitempty"` // 8 bytes

}
