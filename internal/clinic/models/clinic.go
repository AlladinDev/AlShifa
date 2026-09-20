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

type WalletDetails struct {
	AvailableBalance int64 `json:"availableBalance" bson:"availableBalance"`
}

// Clinic represents the details of a clinic, reordered for alignment.
type Clinic struct {
	ID              primitive.ObjectID    `json:"_id,omitempty"  bson:"_id"`
	CreatedAt       time.Time             `json:"registrationDate"  bson:"registrationDate"`
	Name            string                `json:"name"  bson:"name"`                  // 16 bytes
	Address         string                `json:"address" bson:"address"`             // 16 bytes
	SeasonTimings   []SeasonTimingDetails `json:"seasonTimings" bson:"seasonTimings"` // 8 bytes (pointer)
	Mobile          string                `json:"mobile" bson:"mobile"`               // 8 bytes (int64 for phone numbers)
	Pincode         int32                 `json:"pincode" bson:"pincode"`             // 4 bytes
	Wallet          *WalletDetails        `json:"wallet" bson:"wallet"`
	PlanID          primitive.ObjectID    `json:"planID" bson:"planID"`
	MaxAppointments int                   `json:"maxAppointments" bson:"maxAppointments"`
	Verified        bool                  `json:"verified" bson:"verified"`
	OwnerID         primitive.ObjectID    `json:"ownerID" bson:"ownerID"`
	WorkingDays     []string              `json:"workingDays" bson:"workingDays"`
}
