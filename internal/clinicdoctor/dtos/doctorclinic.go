package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Clinic struct {
	ID               primitive.ObjectID    `json:"_id" bson:"_id"`
	AvailableOn      []string              `json:"availableOn" bson:"availableOn"`
	ConsultationFees int                   `json:"consultationFees" bson:"consultationFees"`
	Timings          []ConsultationTimings `json:"timings" bson:"timings"`
	CreatedAt        time.Time             `json:"createdAt" bson:"createdAt"`
	Name             string                `json:"name" bson:"name"`
	Departments      []string              `json:"departments" bson:"departments"`
	WorkingDays      []string              `json:"workingDays" bson:"workingDays"`
	Address          string                `json:"address" bson:"address"`
	SeasonTimings    []SeasonTimingDetails `json:"seasonTimings,omitempty" bson:"seasonTimings"`
}

type DoctorWithClinic struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	DoctorName     string             `json:"doctorName" bson:"doctorName"`
	Experience     int                `json:"experience" bson:"experience"`
	Speciality     string             `json:"speciality" bson:"speciality"`
	DoctorPhotoURL string             `json:"doctorPhotoURL" bson:"doctorPhotoURL"`
	WorkingAt      string             `json:"workingAt" bson:"workingAt"`
	Qualifications string             `json:"qualifications" bson:"qualifications"`
	Post           string             `json:"post" bson:"post"`
}
