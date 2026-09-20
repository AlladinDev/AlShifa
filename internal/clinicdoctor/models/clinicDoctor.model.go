package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsultationTimings struct {
	Day       string    `json:"day" bson:"day"`
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`
}

type ClinicDoctorMapping struct {
	ID               primitive.ObjectID    `json:"_id" bson:"_id"`
	DoctorID         primitive.ObjectID    `json:"doctorID" bson:"doctorID"`
	ClinicID         primitive.ObjectID    `json:"clinicID" bson:"clinicID"`
	AvailableOn      []string              `json:"availableOn" bson:"availableOn"`
	ConsultationFees int                   `json:"consultationFees" bson:"consultationFees"`
	Timings          []ConsultationTimings `json:"timings" bson:"timings"`
	CreatedAt        time.Time             `json:"createdAt" bson:"createdAt"`
}
