package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Clinic struct {
	MappingID        primitive.ObjectID    `json:"mappingID" bson:"mappingID"`
	ID               primitive.ObjectID    `json:"_id" bson:"_id"`
	AvailableOn      []string              `json:"availableOn" bson:"availableOn"`
	ConsultationFees int                   `json:"consultationFees" bson:"consultationFees"`
	DoctorTimings    []ConsultationTimings `json:"doctorTimings" bson:"doctorTimings"`
	JoinedAt         time.Time             `json:"joinedAt" bson:"joinedAt"`
	Name             string                `json:"name" bson:"name"`
	Departments      []string              `json:"departments" bson:"departments"`
	WorkingDays      []string              `json:"workingDays" bson:"workingDays"`
	Address          string                `json:"address" bson:"address"`
	SeasonTimings    []SeasonTimingDetails `json:"seasonTimings" bson:"seasonTimings"`
}

type DoctorWithClinic struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	Experience     int                `json:"experience" bson:"experience"`
	Speciality     string             `json:"speciality" bson:"field"`
	PhotoURL       string             `json:"photoUrl" bson:"profilePhoto"`
	WorkingAt      string             `json:"workingAt" bson:"workingAt"`
	Qualifications string             `json:"qualifications" bson:"qualifications"`
	Post           string             `json:"post" bson:"post"`
	Clinics        []Clinic           `json:"clinics" bson:"clinics"`
}
