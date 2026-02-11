// Package models stores database models for clinic
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model OwnerListResponse

// Owner represents the owner of a clinic, including personal and clinic details.
type Owner struct {
	RegistrationDate time.Time          `json:"registrationDate,omitzero" bson:"registrationDate"` // 24 bytes
	Name             string             `json:"name,omitempty" bson:"name"`                        // 16 bytes
	Address          string             `json:"address,omitempty" bson:"address"`                  // 16 bytes
	Password         string             `json:"password,omitempty" bson:"password"`                // 16 bytes
	Email            string             `json:"email,omitempty" bson:"email"`                      // 16 bytes
	Gender           string             `json:"gender,omitempty" bson:"gender"`                    // 16 bytes
	Clinic           primitive.ObjectID `json:"clinic,omitempty" bson:"clinic"`                    // 8 bytes (pointer)
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Role             string             `json:"role,omitempty" bson:"role"` // 12 bytes (placed near end)
	Mobile           int64              `json:"mobile,omitempty" bson:"mobile"`
	ClinicDetails    *Clinic            `json:"clinicDetails,omitempty" bson:"clinicDetails"` // 8 bytes
}
