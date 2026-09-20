package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsultationTimings struct {
	Day       string    `json:"day" bson:"day"`
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`
}

// SeasonTimingDetails represents the timing details for a season, grouped by size (time.Time is 24 bytes, string is 16).
type SeasonTimingDetails struct {
	Start time.Time `json:"start" bson:"start"` // 24 bytes
	End   time.Time `json:"end" bson:"end"`     // 24 bytes
	Name  string    `json:"name" bson:"name"`   // 16 bytes
}

type ClinicDoctorDTO struct {
	ID                  primitive.ObjectID    `json:"_id" bson:"_id"`
	DoctorID            primitive.ObjectID    `json:"doctorID" bson:"doctorID"`
	ClinicID            primitive.ObjectID    `json:"clinicID" bson:"clinicID"`
	AvailableOn         []string              `json:"availableOn" bson:"availableOn"`
	ConsultationFees    int                   `json:"consultationFees" bson:"consultationFees"`
	Timings             []ConsultationTimings `json:"timings" bson:"timings"`
	CreatedAt           time.Time             `json:"createdAt" bson:"createdAt"`
	DoctorName          string                `json:"doctorName" bson:"doctorName"`
	ClinicName          string                `json:"clinicName" bson:"clinicName"`
	Experience          int                   `json:"experience" bson:"experience"`
	Speciality          string                `json:"speciality" bson:"speciality"`
	ClinicWorkingDays   []string              `json:"clinicWorkingDays" bson:"clinicWorkingDays"`
	DoctorPhotoURL      string                `json:"doctorPhotoURL" bson:"doctorPhotoURL"`
	ClinicAddress       string                `json:"clinicAddress" bson:"clinicAddress"`
	ClinicSeasonTimings []SeasonTimingDetails `json:"clinicSeasonTimings,omitempty" bson:"clinicSeasonTimings"`
	DoctorJoinedOn      time.Time             `json:"doctorJoinedOn" bson:"doctorJoinedOn"`
}
