// Package dtos contains data transfer objects for the clinic application.
package dtos

import (
	"AlShifa/clinic/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clinicDetails struct {
	ID               primitive.ObjectID           `json:"id,omitempty"  bson:"_id"`
	RegistrationDate time.Time                    `json:"registrationDate,omitempty"  bson:"registrationDate"`
	Name             string                       `json:"name,omitempty"  bson:"name"`                  // 16 bytes
	Address          string                       `json:"address,omitempty" bson:"address"`             // 16 bytes
	SeasonTimings    []models.SeasonTimingDetails `json:"seasonTimings,omitempty" bson:"seasonTimings"` // 8 bytes (pointer)
	Mobile           int64                        `json:"mobile,omitempty" bson:"mobile"`               // 8 bytes (int64 for phone numbers)
	Pincode          int32                        `json:"pincode,omitempty" bson:"pincode"`             // 4 bytes
	MaxAppointments  int                          `json:"maxAppointments" bson:"maxAppointments"`
	Verified         bool                         `json:"verified" bson:"verified"`
	StartTime        time.Time                    `json:"StartTime" bson:"startTime"`
	EndTime          time.Time                    `json:"endTime" bson:"endTime"`
	WorkingDays      []string                     `json:"workingDays" bson:"workingDays"`
}

// DoctorAtclinicsDTO is for getting doctordetails and at which clinics that doctor is available
type DoctorAtclinicsDTO struct {
	DoctorDetails models.Doctor   `json:"doctorDetails" bson:"doctorDetails"`
	clinics       []clinicDetails `json:"clinics" bson:"clinics"`
}
