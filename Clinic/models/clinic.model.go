package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SeasonTimingDetails represents the timing details for a season, grouped by size (time.Time is 24 bytes, string is 16).
type SeasonTimingDetails struct {
	Start time.Time `json:"start" bson:"start"` // 24 bytes
	End   time.Time `json:"end" bson:"end"`     // 24 bytes
	Name  string    `json:"name" bson:"name"`   // 16 bytes
}

// Clinic represents the details of a clinic, reordered for alignment.
type Clinic struct {
	ID               primitive.ObjectID    `json:"id" bson:"_id"`
	Owner            primitive.ObjectID    `json:"owner" bson:"owner"`
	RegistrationDate time.Time             `json:"registrationDate" bson:"registrationDate"`
	Name             string                `json:"name" bson:"name"`                   // 16 bytes
	Address          string                `json:"address" bson:"address"`             // 16 bytes
	SeasonTimings    []SeasonTimingDetails `json:"seasonTimings" bson:"seasonTimings"` // 8 bytes (pointer)
	Mobile           int64                 `json:"mobile" bson:"mobile"`               // 8 bytes (int64 for phone numbers)
	Pincode          int32                 `json:"pincode" bson:"pincode"`             // 4 bytes
	Doctors          []primitive.ObjectID  `json:"doctors" bson:"doctors"`
	Wallet           primitive.ObjectID    `json:"wallet" bson:"wallet"`
	OwnerDetails     *Owner                `bson:"ownerDetails,omitempty"`
	DoctorDetails    []Doctor              `bson:"doctorDetails,omitempty"`
}
