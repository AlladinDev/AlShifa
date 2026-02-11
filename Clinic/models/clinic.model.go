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
	ID               primitive.ObjectID    `json:"id,omitempty"  bson:"_id"`
	RegistrationDate time.Time             `json:"registrationDate,omitempty"  bson:"registrationDate"`
	Name             string                `json:"name,omitempty"  bson:"name"`                  // 16 bytes
	Address          string                `json:"address,omitempty" bson:"address"`             // 16 bytes
	SeasonTimings    []SeasonTimingDetails `json:"seasonTimings,omitempty" bson:"seasonTimings"` // 8 bytes (pointer)
	Mobile           int64                 `json:"mobile,omitempty" bson:"mobile"`               // 8 bytes (int64 for phone numbers)
	Pincode          int32                 `json:"pincode,omitempty" bson:"pincode"`             // 4 bytes
	Wallet           *WalletDetails        `json:"wallet,omitempty" bson:"wallet"`
	OwnerDetails     *Owner                `bson:"ownerDetails,omitempty"`
	DoctorDetails    []Doctor              `bson:"doctorDetails,omitempty"`
	PlanType         string                `json:"planType,omitempty" bson:"planType"`
	MaxAppointments  int                   `json:"maxAppointments" bson:"maxAppointments"`
}
